package filters_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"authentication_service/api"
	"authentication_service/filters"
	"authentication_service/internal/logger"
	"authentication_service/internal/trace"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

// syncBuffer is a goroutine-safe log sink.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func (b *syncBuffer) lines(t *testing.T) []map[string]interface{} {
	t.Helper()
	var out []map[string]interface{}
	for _, l := range strings.Split(strings.TrimSpace(b.String()), "\n") {
		if l == "" {
			continue
		}
		m := map[string]interface{}{}
		if err := json.Unmarshal([]byte(l), &m); err != nil {
			t.Fatalf("non-JSON log line %q: %v", l, err)
		}
		out = append(out, m)
	}
	return out
}

func (b *syncBuffer) accessLines(t *testing.T) []map[string]interface{} {
	var out []map[string]interface{}
	for _, l := range b.lines(t) {
		if l["msg"] == "request completed" {
			out = append(out, l)
		}
	}
	return out
}

var downstream struct {
	mu      sync.Mutex
	headers http.Header
}

// newApp builds a fresh Beego server wired exactly like main.go.
func newApp(t *testing.T) (http.Handler, *syncBuffer) {
	t.Helper()
	logs := &syncBuffer{}
	logger.Init(logs, slog.LevelDebug)

	cfg := *beego.BConfig
	cfg.RunMode = beego.DEV // dev mode: Beego would print its own router line
	cfg.RecoverPanic = true
	cfg.CopyRequestBody = true
	app := beego.NewHttpServerWithCfg(&cfg)

	filters.Install(app)

	// An early-responding BeforeRouter filter (like CORS preflight).
	app.InsertFilter("/early", beego.BeforeRouter, func(ctx *context.Context) {
		ctx.Output.SetStatus(204)
		_ = ctx.Output.Body([]byte(""))
	})

	app.Get("/ok", func(ctx *context.Context) {
		slog.InfoContext(ctx.Request.Context(), "handler ran")
		_ = ctx.Output.Body([]byte("hello"))
	})
	app.Get("/panic", func(ctx *context.Context) {
		panic("boom: something broke")
	})
	app.Get("/fail", func(ctx *context.Context) {
		filters.RespondError(ctx, http.StatusBadRequest, "bad_request", "Invalid input")
	})
	app.Post("/users/:username", func(ctx *context.Context) {
		_ = ctx.Output.Body([]byte("ok"))
	})
	app.Get("/outbound", func(ctx *context.Context) {
		srv := ctx.Input.Query("target")
		third := ctx.Input.Query("third") == "1"
		req := api.NewInternalRequest(srv, "/downstream", api.GET)
		if third {
			req = api.NewRequest(srv, "/downstream", api.GET)
		}
		client := api.Client{Request: req, Ctx: ctx.Request.Context()}
		res, err := client.SendRequest()
		if err != nil {
			t.Errorf("outbound call failed: %v", err)
			return
		}
		res.Body.Close()
		_ = ctx.Output.Body([]byte("called"))
	})

	app.Handlers.Init()
	return app.Handlers, logs
}

func do(h http.Handler, method, target string, hdr map[string]string, body string) *httptest.ResponseRecorder {
	var r io.Reader
	if body != "" {
		r = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, target, r)
	for k, v := range hdr {
		req.Header[http.CanonicalHeaderKey(k)] = []string{v}
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

const (
	validReq   = "req_0123456789abcdef"
	validTrace = "trace-0123456789abcdef"
	validSess  = "sess_0123456789"
)

func TestValidHeadersArePropagated(t *testing.T) {
	h, logs := newApp(t)
	rec := do(h, "GET", "/ok", map[string]string{
		trace.HeaderRequestID: validReq, trace.HeaderTraceID: validTrace, trace.HeaderSessionID: validSess,
	}, "")

	if rec.Code != 200 || rec.Body.String() != "hello" {
		t.Fatalf("unexpected response %d %q", rec.Code, rec.Body.String())
	}
	if rec.Header().Get(trace.HeaderRequestID) != validReq || rec.Header().Get(trace.HeaderTraceID) != validTrace {
		t.Fatalf("response headers = %v", rec.Header())
	}
	if rec.Header().Get(trace.HeaderSessionID) != "" {
		t.Fatalf("session id must not be echoed on the response")
	}
	for _, l := range logs.lines(t) {
		if l["service"] != logger.ServiceName {
			t.Fatalf("line without service: %v", l)
		}
		if l["msg"] == "handler ran" || l["msg"] == "request completed" {
			if l["request_id"] != validReq || l["trace_id"] != validTrace || l["session_id"] != validSess {
				t.Fatalf("line missing ids: %v", l)
			}
		}
	}
}

func TestMissingHeadersStillWork(t *testing.T) {
	h, logs := newApp(t)
	rec := do(h, "GET", "/ok", nil, "")
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	reqID := rec.Header().Get(trace.HeaderRequestID)
	if !trace.Valid(reqID) {
		t.Fatalf("request id not generated: %q", reqID)
	}
	if rec.Header().Get(trace.HeaderTraceID) != reqID {
		t.Fatalf("trace id should fall back to request id")
	}
	access := logs.accessLines(t)
	if len(access) != 1 || access[0]["request_id"] != reqID || access[0]["session_id"] != "" {
		t.Fatalf("access log = %v", access)
	}
}

func TestInvalidAndInjectionHeadersAreDiscardedAndNeverLogged(t *testing.T) {
	cases := map[string]string{
		"too short":      "abc",
		"too long":       strings.Repeat("a", 65),
		"crlf injection": "abcdefgh\r\nX-Injected: yes",
		"json injection": `abcdefgh","level":"ERROR`,
		"html":           "<script>alert('xss')</script>",
		"sql":            "1' OR '1'='1' --",
		"unicode":        "ááááááááááá",
		"spaces":         "abcd efgh ijkl",
	}
	for name, bad := range cases {
		t.Run(name, func(t *testing.T) {
			h, logs := newApp(t)
			rec := do(h, "GET", "/ok", map[string]string{
				trace.HeaderRequestID: bad, trace.HeaderTraceID: bad, trace.HeaderSessionID: bad,
			}, "")
			if rec.Code != 200 {
				t.Fatalf("request with invalid headers must still succeed, got %d", rec.Code)
			}
			reqID := rec.Header().Get(trace.HeaderRequestID)
			if reqID == bad || !trace.Valid(reqID) || rec.Header().Get(trace.HeaderTraceID) != reqID {
				t.Fatalf("invalid id not replaced: req=%q trace=%q", reqID, rec.Header().Get(trace.HeaderTraceID))
			}
			out := logs.String()
			if strings.Contains(out, bad) || strings.Contains(out, "X-Injected") || strings.Contains(out, "script") {
				t.Fatalf("raw invalid value was logged:\n%s", out)
			}
			warned := map[string]bool{}
			for _, l := range logs.lines(t) {
				if l["msg"] == "invalid tracing header discarded" {
					warned[l["header"].(string)] = true
					if l["request_id"] != reqID {
						t.Fatalf("warning must carry the generated request id: %v", l)
					}
				}
			}
			for _, hname := range []string{trace.HeaderRequestID, trace.HeaderTraceID, trace.HeaderSessionID} {
				if !warned[hname] {
					t.Fatalf("no warning for %s; got %v", hname, warned)
				}
			}
			if n := len(logs.accessLines(t)); n != 1 {
				t.Fatalf("want 1 access line, got %d", n)
			}
		})
	}
}

func TestValidRequestIDWithInvalidTraceID(t *testing.T) {
	h, _ := newApp(t)
	rec := do(h, "GET", "/ok", map[string]string{trace.HeaderRequestID: validReq, trace.HeaderTraceID: "bad"}, "")
	if rec.Header().Get(trace.HeaderRequestID) != validReq || rec.Header().Get(trace.HeaderTraceID) != validReq {
		t.Fatalf("headers %v", rec.Header())
	}
}

func decodeErr(t *testing.T, body []byte) filters.ErrorBody {
	t.Helper()
	var eb filters.ErrorBody
	if err := json.Unmarshal(body, &eb); err != nil {
		t.Fatalf("error body not JSON: %q", body)
	}
	return eb
}

func TestErrorResponseHasRequestID(t *testing.T) {
	h, logs := newApp(t)
	rec := do(h, "GET", "/fail", map[string]string{trace.HeaderRequestID: validReq}, "")
	if rec.Code != 400 {
		t.Fatalf("status %d", rec.Code)
	}
	if rec.Header().Get(trace.HeaderRequestID) != validReq || rec.Header().Get(trace.HeaderTraceID) != validReq {
		t.Fatalf("response headers missing on error: %v", rec.Header())
	}
	if !strings.HasPrefix(rec.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("content type %q", rec.Header().Get("Content-Type"))
	}
	eb := decodeErr(t, rec.Body.Bytes())
	if eb.Error.Code != "bad_request" || eb.Error.Message != "Invalid input" || eb.Error.RequestID != validReq {
		t.Fatalf("error body %+v", eb)
	}
	access := logs.accessLines(t)
	if len(access) != 1 || access[0]["status"] != float64(400) {
		t.Fatalf("access %v", access)
	}
}

func TestPanicIsRecovered(t *testing.T) {
	h, logs := newApp(t)
	rec := do(h, "GET", "/panic", map[string]string{trace.HeaderRequestID: validReq, trace.HeaderTraceID: validTrace}, "")

	if rec.Code != 500 {
		t.Fatalf("status %d", rec.Code)
	}
	if rec.Header().Get(trace.HeaderRequestID) != validReq || rec.Header().Get(trace.HeaderTraceID) != validTrace {
		t.Fatalf("response headers missing on panic: %v", rec.Header())
	}
	body := rec.Body.String()
	eb := decodeErr(t, rec.Body.Bytes())
	if eb.Error.Code != "internal_error" || eb.Error.RequestID != validReq || eb.Error.Message == "" {
		t.Fatalf("body %+v", eb)
	}
	if strings.Contains(body, "boom") || strings.Contains(body, ".go:") || strings.Contains(body, "goroutine") {
		t.Fatalf("panic details leaked to client: %s", body)
	}

	var panicLine map[string]interface{}
	for _, l := range logs.lines(t) {
		if l["msg"] == "panic recovered" {
			panicLine = l
		}
	}
	if panicLine == nil {
		t.Fatalf("panic not logged:\n%s", logs.String())
	}
	if panicLine["request_id"] != validReq || panicLine["trace_id"] != validTrace || panicLine["level"] != "ERROR" {
		t.Fatalf("panic line ids %v", panicLine)
	}
	if !strings.Contains(panicLine["stack"].(string), "filters_test.go") || panicLine["panic"] != "boom: something broke" {
		t.Fatalf("panic line missing stack/panic: %v", panicLine)
	}
	access := logs.accessLines(t)
	if len(access) != 1 || access[0]["status"] != float64(500) || access[0]["request_id"] != validReq {
		t.Fatalf("want exactly one access line with status 500, got %v", access)
	}
}

func TestExactlyOneAccessLine(t *testing.T) {
	cases := []struct {
		name, method, path string
		status             float64
		route              string
	}{
		{"success", "GET", "/ok", 200, "/ok"},
		{"not found", "GET", "/does-not-exist", 404, ""},
		{"early filter response", "GET", "/early", 204, ""},
		{"panic", "GET", "/panic", 500, "/panic"},
		{"error helper", "GET", "/fail", 400, "/fail"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h, logs := newApp(t)
			rec := do(h, c.method, c.path, nil, "")
			access := logs.accessLines(t)
			if len(access) != 1 {
				t.Fatalf("want 1 access line, got %d:\n%s", len(access), logs.String())
			}
			a := access[0]
			if a["status"] != c.status || float64(rec.Code) != c.status {
				t.Fatalf("status logged=%v recorded=%d want=%v", a["status"], rec.Code, c.status)
			}
			if a["method"] != c.method || a["path"] != c.path || a["route"] != c.route {
				t.Fatalf("access line %v", a)
			}
			if a["request_id"] != rec.Header().Get(trace.HeaderRequestID) {
				t.Fatalf("access line request id %v vs header %q", a["request_id"], rec.Header().Get(trace.HeaderRequestID))
			}
			if _, ok := a["duration_ms"].(float64); !ok {
				t.Fatalf("duration missing: %v", a)
			}
			if a["bytes"] != float64(rec.Body.Len()) {
				t.Fatalf("bytes logged %v, body %d", a["bytes"], rec.Body.Len())
			}
			// Beego's built-in dev-mode router line must not appear.
			if strings.Contains(logs.String(), "match|") {
				t.Fatalf("duplicate built-in access log present:\n%s", logs.String())
			}
		})
	}
}

func TestOutboundCallsForwardSameIDs(t *testing.T) {
	ds := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		downstream.mu.Lock()
		downstream.headers = r.Header.Clone()
		downstream.mu.Unlock()
		w.WriteHeader(200)
	}))
	defer ds.Close()
	got := func() http.Header {
		downstream.mu.Lock()
		defer downstream.mu.Unlock()
		return downstream.headers
	}

	t.Run("internal call forwards inbound ids unchanged", func(t *testing.T) {
		h, _ := newApp(t)
		rec := do(h, "GET", "/outbound?target="+ds.URL, map[string]string{
			trace.HeaderRequestID: validReq, trace.HeaderTraceID: validTrace, trace.HeaderSessionID: validSess,
		}, "")
		if rec.Code != 200 {
			t.Fatalf("status %d %s", rec.Code, rec.Body.String())
		}
		dh := got()
		if dh.Get(trace.HeaderRequestID) != validReq || dh.Get(trace.HeaderTraceID) != validTrace || dh.Get(trace.HeaderSessionID) != validSess {
			t.Fatalf("downstream headers %v", dh)
		}
	})

	t.Run("generated ids are forwarded, not regenerated", func(t *testing.T) {
		h, _ := newApp(t)
		rec := do(h, "GET", "/outbound?target="+ds.URL, nil, "")
		dh := got()
		if dh.Get(trace.HeaderRequestID) == "" || dh.Get(trace.HeaderRequestID) != rec.Header().Get(trace.HeaderRequestID) {
			t.Fatalf("downstream request id %q, response %q", dh.Get(trace.HeaderRequestID), rec.Header().Get(trace.HeaderRequestID))
		}
		if dh.Get(trace.HeaderTraceID) != rec.Header().Get(trace.HeaderTraceID) {
			t.Fatalf("trace id mismatch")
		}
		if _, ok := dh[trace.HeaderSessionID]; ok {
			t.Fatalf("absent session id must not be invented")
		}
	})

	t.Run("third-party call does not receive internal ids", func(t *testing.T) {
		h, _ := newApp(t)
		do(h, "GET", "/outbound?third=1&target="+ds.URL, map[string]string{trace.HeaderRequestID: validReq}, "")
		dh := got()
		for _, k := range []string{trace.HeaderRequestID, trace.HeaderTraceID, trace.HeaderSessionID} {
			if dh.Get(k) != "" {
				t.Fatalf("%s forwarded to third party", k)
			}
		}
	})
}

func TestNoSensitiveDataInLogs(t *testing.T) {
	h, logs := newApp(t)
	secrets := []string{
		"Bearer-SECRET-ACCESS-TOKEN", "SESSIONCOOKIEVALUE", "hunter2", "QUERYTOKENVALUE",
		"john.doe@example.com", "BODYPASSWORD",
	}
	rec := do(h, "POST", "/users/john.doe@example.com?password=hunter2&token=QUERYTOKENVALUE", map[string]string{
		"Authorization":       "Bearer Bearer-SECRET-ACCESS-TOKEN",
		"Cookie":              "session=SESSIONCOOKIEVALUE",
		"Content-Type":        "application/json",
		trace.HeaderRequestID: validReq,
	}, `{"Username":"john.doe@example.com","Password":"BODYPASSWORD"}`)
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	out := logs.String()
	for _, s := range secrets {
		if strings.Contains(out, s) {
			t.Fatalf("sensitive value %q found in logs:\n%s", s, out)
		}
	}
	access := logs.accessLines(t)
	if len(access) != 1 || access[0]["path"] != "/users/:username" || access[0]["route"] != "/users/:username" {
		t.Fatalf("path params must be masked, query dropped: %v", access)
	}
}
