package filters

import (
	"bufio"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

const maxLoggedPathLen = 256

// AccessLog is a Beego filter chain (web.InsertFilterChain) wrapping the whole
// router. It emits exactly one access-log line per request, including for
// static files, 404/405, CORS preflight, early filter responses and panics.
//
// Logged: method, path (no query string; path parameter values replaced by
// their names), route pattern, status, duration_ms, bytes, plus the IDs.
// Never logged: query values, headers, cookies, bodies.
func AccessLog(next beego.FilterFunc) beego.FilterFunc {
	return func(ctx *context.Context) {
		start := time.Now()
		EnsureTraceIDs(ctx)

		cw := &countingWriter{ResponseWriter: ctx.ResponseWriter.ResponseWriter}
		ctx.ResponseWriter.ResponseWriter = cw

		defer func() {
			// Safety net: Beego's RecoverFunc (RecoverPanic) normally handles
			// handler panics; this catches anything outside it.
			if rec := recover(); rec != nil {
				if rec != beego.ErrAbort {
					handlePanic(ctx, rec, debug.Stack())
				}
			}
			logAccess(ctx, start, cw)
		}()

		next(ctx)
	}
}

func logAccess(ctx *context.Context, start time.Time, cw *countingWriter) {
	status := cw.status
	if status == 0 {
		// Nothing written: net/http will send 200.
		status = http.StatusOK
	}
	slog.InfoContext(ctx.Request.Context(), "request completed",
		"method", ctx.Request.Method,
		"path", loggablePath(ctx),
		"route", routePattern(ctx),
		"status", status,
		"duration_ms", float64(time.Since(start).Microseconds())/1000.0,
		"bytes", cw.bytes,
	)
}

func routePattern(ctx *context.Context) string {
	if p, ok := ctx.Input.GetData("RouterPattern").(string); ok {
		return p
	}
	return ""
}

// loggablePath returns the URL path without the query string. Path parameter
// values (e.g. a username in /change-customer-password/:username) are
// replaced by the parameter name so no PII lands in the access log.
func loggablePath(ctx *context.Context) string {
	p := ctx.Request.URL.EscapedPath()
	// Only mask parameters declared by the matched route (filters with "*"
	// patterns also add a ":splat" param that is not a route parameter).
	route := routePattern(ctx)
	params := map[string]string{}
	for name, val := range ctx.Input.Params() {
		if strings.HasPrefix(name, ":") && name != ":splat" && route != "" && strings.Contains(route, name) {
			params[name] = val
		}
	}
	if len(params) > 0 {
		segs := strings.Split(p, "/")
		for i, seg := range segs {
			if seg == "" {
				continue
			}
			unesc, err := url.PathUnescape(seg)
			if err != nil {
				unesc = seg
			}
			for name, val := range params {
				if val != "" && (val == unesc || val == seg) {
					segs[i] = name
					break
				}
			}
		}
		p = strings.Join(segs, "/")
	}
	if len(p) > maxLoggedPathLen {
		p = p[:maxLoggedPathLen] + "..."
	}
	return p
}

// countingWriter records the status code and body bytes written.
type countingWriter struct {
	http.ResponseWriter
	status int
	bytes  int64
}

func (w *countingWriter) WriteHeader(code int) {
	if w.status == 0 {
		w.status = code
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *countingWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += int64(n)
	return n, err
}

func (w *countingWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *countingWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, errors.New("hijack not supported")
}

func (w *countingWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }
