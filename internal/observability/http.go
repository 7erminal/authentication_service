package observability

import (
	"authentication_service/internal/trace"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

type requestState struct{ route string }
type stateKey struct{}

func initialize(r *http.Request) {
	if trace.FromContext(r.Context()).RequestID == "" {
		*r = *r.WithContext(trace.WithIDs(r.Context(), trace.ReadHeaders(r.Header)))
	}
}

// BeforeRouter runs before CORS, so early preflight responses retain the IDs.
// Middleware also initializes them for static files and pre-router failures.
func BeforeRouter(c *beecontext.Context) {
	initialize(c.Request)
	ids := trace.FromContext(c.Request.Context())
	c.Output.Header(trace.RequestHeader, ids.RequestID)
	c.Output.Header(trace.TraceHeader, ids.TraceID)
}

func BeforeExec(c *beecontext.Context) {
	if state, ok := c.Request.Context().Value(stateKey{}).(*requestState); ok {
		state.route, _ = c.Input.GetData("RouterPattern").(string)
	}
}

// Recover hands real panics to the outer middleware after Beego releases its
// pooled context. ErrAbort is Beego's intentional early-response sentinel.
func Recover(_ *beecontext.Context, _ *beego.Config) {
	if value := recover(); value != nil {
		if err, ok := value.(error); ok && errors.Is(err, beego.ErrAbort) {
			return
		}
		panic(value)
	}
}

type errorBody struct {
	Error struct {
		Code      string `json:"code"`
		Message   string `json:"message"`
		RequestID string `json:"request_id"`
	} `json:"error"`
}

func RespondError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	ids := trace.FromContext(r.Context())
	w.Header().Set(trace.RequestHeader, ids.RequestID)
	w.Header().Set(trace.TraceHeader, ids.TraceID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var body errorBody
	body.Error.Code, body.Error.Message, body.Error.RequestID = code, message, ids.RequestID
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// stackLocations omits panic values and argument dumps, which can contain secrets.
func stackLocations() string {
	pcs := make([]uintptr, 64)
	n := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	var stack strings.Builder
	for {
		frame, more := frames.Next()
		stack.WriteString(frame.Function)
		stack.WriteString(" at " + frame.File + ":" + strconv.Itoa(frame.Line) + "\n")
		// No runtime.Stack argument dump or recovered value is emitted.
		if !more {
			break
		}
	}
	return stack.String()
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		initialize(r)
		state := &requestState{route: "<unmatched>"}
		r = r.WithContext(context.WithValue(r.Context(), stateKey{}, state))
		rw := &responseWriter{ResponseWriter: w, header: w.Header().Clone()}
		ids := trace.FromContext(r.Context())
		rw.Header().Set(trace.RequestHeader, ids.RequestID)
		rw.Header().Set(trace.TraceHeader, ids.TraceID)
		defer func() {
			if value := recover(); value != nil {
				trace.Logger(r.Context()).Error("request panic", "stack", stackLocations(), "response_committed", rw.committed)
				if !rw.committed {
					rw.resetForError()
					RespondError(rw, r, http.StatusInternalServerError, "internal_error", "Internal server error")
				}
			}
			rw.commit()
			status := rw.status
			if status == 0 {
				status = http.StatusOK
			}
			// Use the route template as the redacted path. Unmatched/static paths
			// may contain arbitrary PII and are deliberately never logged raw.
			root := accessLogger
			if root == nil {
				root = slog.Default()
			}
			trace.LoggerFrom(r.Context(), root).Info("http access", "method", r.Method, "path", state.route, "route", state.route,
				"status", status, "duration_ms", float64(time.Since(start).Microseconds())/1000, "bytes", rw.bytes)
		}()
		next.ServeHTTP(rw, r)
	})
}

// Buffer ordinary API responses so a panic after a partial JSON write can still
// be replaced by a clean 500. Memory is bounded per request. Explicit streaming
// (Flush/Hijack) and responses over 1 MiB commit to the network and cannot later
// change HTTP status; those exceptional panics are still safely logged.
const maxBufferedResponse = 1 << 20

type responseWriter struct {
	http.ResponseWriter
	header    http.Header
	body      bytes.Buffer
	status    int
	bytes     int64
	committed bool
}

func (w *responseWriter) Header() http.Header         { return w.header }
func (w *responseWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }
func (w *responseWriter) copyHeaders() {
	target := w.ResponseWriter.Header()
	clear(target)
	for key, values := range w.header {
		target[key] = append([]string(nil), values...)
	}
}
func (w *responseWriter) commit() {
	if w.committed {
		return
	}
	w.committed = true
	if w.status == 0 {
		w.status = http.StatusOK
	}
	w.copyHeaders()
	w.ResponseWriter.WriteHeader(w.status)
	if w.body.Len() > 0 {
		n, _ := w.ResponseWriter.Write(w.body.Bytes())
		w.bytes = int64(n)
		w.body.Reset()
	}
}
func (w *responseWriter) resetForError() {
	w.status, w.bytes = 0, 0
	w.body.Reset()
	for key := range w.header {
		// Retain only correlation and CORS headers. Cookies, redirects, compression
		// metadata and other partial success headers must not survive a panic.
		if key != trace.RequestHeader && key != trace.TraceHeader && !strings.HasPrefix(key, "Access-Control-") {
			w.header.Del(key)
		}
	}
}
func (w *responseWriter) WriteHeader(status int) {
	if w.status != 0 || w.committed {
		return
	}
	if status < 100 || status > 999 {
		panic("invalid HTTP status")
	}
	if status >= 100 && status < 200 && status != http.StatusSwitchingProtocols {
		w.copyHeaders()
		w.ResponseWriter.WriteHeader(status)
		return
	}
	w.status = status
	if status == http.StatusSwitchingProtocols {
		w.commit()
	}
}
func (w *responseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	if w.status == http.StatusNoContent || w.status == http.StatusNotModified {
		return 0, http.ErrBodyNotAllowed
	}
	if !w.committed && w.body.Len()+len(b) <= maxBufferedResponse {
		n, err := w.body.Write(b)
		w.bytes += int64(n)
		return n, err
	}
	w.commit()
	n, err := w.ResponseWriter.Write(b)
	w.bytes += int64(n)
	return n, err
}
func (w *responseWriter) Flush() {
	w.commit()
	_ = http.NewResponseController(w.ResponseWriter).Flush()
}
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	conn, reader, err := http.NewResponseController(w.ResponseWriter).Hijack()
	if err == nil {
		w.committed = true
	}
	return conn, reader, err
}

func Install() {
	beego.BConfig.Log.AccessLogs = false
	beego.BConfig.Log.EnableStaticLogs = false
	beego.BConfig.EnableErrorsRender = false
	beego.BConfig.RecoverPanic = true
	beego.BConfig.RecoverFunc = Recover
	beego.InsertFilter("*", beego.BeforeRouter, BeforeRouter)
	beego.InsertFilter("*", beego.BeforeExec, BeforeExec)
	// ErrorController (not ErrorHandler) is deliberate: Beego's handler-style
	// errors commit headers before invoking the callback, losing Content-Type.
	beego.ErrorController(&frameworkErrors{})
}

type frameworkErrors struct{ beego.Controller }

func (c *frameworkErrors) respond(status int) {
	RespondError(c.Ctx.ResponseWriter, c.Ctx.Request, status, "http_"+strconv.Itoa(status), http.StatusText(status))
}
func (c *frameworkErrors) Error400() { c.respond(400) }
func (c *frameworkErrors) Error401() { c.respond(401) }
func (c *frameworkErrors) Error402() { c.respond(402) }
func (c *frameworkErrors) Error403() { c.respond(403) }
func (c *frameworkErrors) Error404() { c.respond(404) }
func (c *frameworkErrors) Error405() { c.respond(405) }
func (c *frameworkErrors) Error413() { c.respond(413) }
func (c *frameworkErrors) Error417() { c.respond(417) }
func (c *frameworkErrors) Error422() { c.respond(422) }
func (c *frameworkErrors) Error500() { c.respond(500) }
func (c *frameworkErrors) Error501() { c.respond(501) }
func (c *frameworkErrors) Error502() { c.respond(502) }
func (c *frameworkErrors) Error503() { c.respond(503) }
func (c *frameworkErrors) Error504() { c.respond(504) }
