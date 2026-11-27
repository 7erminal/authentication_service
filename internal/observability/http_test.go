package observability

import (
	"authentication_service/internal/trace"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

func testRouter() http.Handler {
	cfg := *beego.BConfig
	cfg.RecoverFunc = Recover
	cfg.RunMode = "dev" // exercise suppression of the duplicate dev access line
	cfg.Log.AccessLogs = false
	cfg.WebConfig.AutoRender = false
	router := beego.NewControllerRegisterWithCfg(&cfg)
	router.InsertFilter("*", beego.BeforeRouter, BeforeRouter)
	router.InsertFilter("*", beego.BeforeExec, BeforeExec)
	router.InsertFilter("/early-panic", beego.BeforeRouter, func(c *beecontext.Context) { panic("PANIC_SECRET") })
	router.InsertFilter("/early", beego.BeforeRouter, func(c *beecontext.Context) {
		c.Output.SetStatus(204)
		_ = c.Output.Body(nil)
	})
	router.Get("/success/:username", func(c *beecontext.Context) { _ = c.Output.Body([]byte("OK")) })
	router.Get("/error", func(c *beecontext.Context) {
		RespondError(c.ResponseWriter, c.Request, 401, "unauthorized", "Not authorized")
	})
	router.Get("/panic", func(c *beecontext.Context) { panic("password=PANIC_SECRET token=PANIC_SECRET") })
	router.Get("/panic-slice", func(c *beecontext.Context) { panic([]string{"PANIC_SECRET"}) })
	router.Get("/abort", func(c *beecontext.Context) { _ = c.Output.Body([]byte("done")); panic(beego.ErrAbort) })
	router.Get("/partial", func(c *beecontext.Context) {
		c.Output.Header("Set-Cookie", "token=SHOULD_NOT_SURVIVE")
		_ = c.Output.Body([]byte("partial"))
		panic("PANIC_SECRET")
	})
	router.Get("/committed", func(c *beecontext.Context) {
		_ = c.Output.Body([]byte("partial"))
		_ = http.NewResponseController(c.ResponseWriter.ResponseWriter).Flush()
		panic("PANIC_SECRET")
	})
	router.Init()
	return Middleware(router)
}

func records(t *testing.T, output string) []map[string]interface{} {
	t.Helper()
	var result []map[string]interface{}
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if line == "" {
			continue
		}
		var record map[string]interface{}
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("non-JSON log: %s", line)
		}
		for _, key := range []string{"service", "session_id", "trace_id", "request_id"} {
			if _, ok := record[key]; !ok {
				t.Fatalf("missing %s: %s", key, line)
			}
		}
		result = append(result, record)
	}
	return result
}

func TestBeegoLifecycle(t *testing.T) {
	var output bytes.Buffer
	if err := ConfigureLogging(&output, "DEBUG"); err != nil {
		t.Fatal(err)
	}
	handler := testRouter()
	for _, test := range []struct {
		path   string
		status int
		body   string
		panic  bool
	}{
		{"/success/private-email@example.com?password=QUERY_SECRET", 200, "OK", false},
		{"/error", 401, "", false},
		{"/early", 204, "", false},
		{"/early-panic", 500, "", true},
		{"/panic", 500, "", true},
		{"/panic-slice", 500, "", true},
		{"/abort", 200, "done", false},
		{"/partial", 500, "", true},
		{"/committed", 200, "partial", true},
	} {
		t.Run(test.path, func(t *testing.T) {
			output.Reset()
			req := httptest.NewRequest("GET", test.path, nil)
			req.Header.Set(trace.RequestHeader, "request_123")
			req.Header.Set(trace.TraceHeader, "trace_123")
			req.Header.Set(trace.SessionHeader, "session_123")
			req.Header.Set("Authorization", "Bearer AUTH_SECRET")
			req.Header.Set("Cookie", "session=COOKIE_SECRET")
			res := httptest.NewRecorder()
			handler.ServeHTTP(res, req)
			if res.Code != test.status {
				t.Fatalf("status = %d, body = %s", res.Code, res.Body.String())
			}
			if res.Header().Get(trace.RequestHeader) != "request_123" || res.Header().Get(trace.TraceHeader) != "trace_123" {
				t.Fatal("missing response correlation headers")
			}
			if test.body != "" && res.Body.String() != test.body {
				t.Fatalf("body: %s", res.Body.String())
			}
			if test.path == "/partial" && res.Result().Header.Get("Set-Cookie") != "" {
				t.Fatal("partial success headers survived panic")
			}
			if test.status >= 400 {
				var body errorBody
				if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}
				if body.Error.RequestID != "request_123" {
					t.Fatal("missing error request ID")
				}
				if test.panic && (body.Error.Code != "internal_error" || body.Error.Message != "Internal server error") {
					t.Fatalf("unsafe panic body: %s", res.Body.String())
				}
			}
			access, panics := 0, 0
			for _, record := range records(t, output.String()) {
				if record["request_id"] != "request_123" || record["trace_id"] != "trace_123" || record["session_id"] != "session_123" {
					t.Fatalf("uncorrelated log: %+v", record)
				}
				if record["msg"] == "http access" {
					access++
					if record["status"] != float64(test.status) || record["bytes"] != float64(res.Body.Len()) {
						t.Fatalf("incorrect access accounting: %+v", record)
					}
					if _, ok := record["duration_ms"]; !ok {
						t.Fatal("no duration")
					}
					if strings.HasPrefix(test.path, "/success/") && record["route"] != "/success/:username" {
						t.Fatalf("route missing: %+v", record)
					}
				}
				if record["msg"] == "request panic" {
					panics++
					if record["stack"] == "" {
						t.Fatal("panic stack missing")
					}
				}
			}
			if access != 1 {
				t.Fatalf("access lines = %d; %s", access, output.String())
			}
			if test.panic != (panics == 1) {
				t.Fatalf("panic lines = %d", panics)
			}
			for _, secret := range []string{"private-email@example.com", "QUERY_SECRET", "PANIC_SECRET", "AUTH_SECRET", "COOKIE_SECRET"} {
				if strings.Contains(output.String(), secret) || strings.Contains(res.Body.String(), "PANIC_SECRET") {
					t.Fatalf("secret leaked: %s", secret)
				}
			}
		})
	}
}

func TestMissingInvalidAndPreRouterHeaders(t *testing.T) {
	var output bytes.Buffer
	_ = ConfigureLogging(&output, "INFO")
	for _, value := range []string{"", "evil\r\nINJECTED_SECRET", "has spaces", strings.Repeat("a", 65)} {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/static/private-path?token=QUERY_SECRET", nil)
		req.Header.Set(trace.RequestHeader, value)
		req.Header.Set(trace.TraceHeader, value)
		req.Header.Set(trace.SessionHeader, value)
		Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ids := trace.FromContext(r.Context())
			if !trace.Valid(ids.RequestID) || ids.RequestID != ids.TraceID || ids.SessionID != "" {
				t.Fatalf("bad IDs: %+v", ids)
			}
			w.WriteHeader(204)
		})).ServeHTTP(res, req)
		if !trace.Valid(res.Header().Get(trace.RequestHeader)) || res.Header().Get(trace.TraceHeader) != res.Header().Get(trace.RequestHeader) {
			t.Fatal("missing fallback headers")
		}
	}
	if strings.Contains(output.String(), "SECRET") || strings.Contains(output.String(), "private-path") {
		t.Fatal("sensitive input leaked")
	}
	if len(records(t, output.String())) != 4 {
		t.Fatal("pre-router access log missing")
	}
}

func TestLoggingLevelAndUnsafeFrameworkDiagnostics(t *testing.T) {
	var output bytes.Buffer
	if err := ConfigureLogging(&output, "warn"); err != nil {
		t.Fatal(err)
	}
	trace.Logger(context.Background()).Info("filtered")
	logs.Debug("duplicate request line")
	logs.Error("token=FRAMEWORK_SECRET")
	orm.DebugLog.Println("password=ORM_SECRET")
	if strings.Contains(output.String(), "SECRET") || strings.Contains(output.String(), "filtered") || strings.Contains(output.String(), "duplicate") {
		t.Fatalf("unsafe logs: %s", output.String())
	}
	records(t, output.String())
	if err := ConfigureLogging(&output, "invalid"); err == nil {
		t.Fatal("invalid level accepted")
	}
}

// lockingBuffer permits concurrent requests under the race detector.
type lockingBuffer struct {
	sync.Mutex
	bytes.Buffer
}

func (b *lockingBuffer) Write(p []byte) (int, error) {
	b.Lock()
	defer b.Unlock()
	return b.Buffer.Write(p)
}

func TestConcurrentRequestsDoNotMixIDs(t *testing.T) {
	var output lockingBuffer
	_ = ConfigureLogging(&output, "INFO")
	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { trace.Logger(r.Context()).Info("working") }))
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req := httptest.NewRequest("GET", "/", nil)
			id := fmt.Sprintf("request_%03d", i)
			req.Header.Set(trace.RequestHeader, id)
			req.Header.Set(trace.TraceHeader, id)
			req.Header.Set(trace.SessionHeader, id)
			handler.ServeHTTP(httptest.NewRecorder(), req)
		}(i)
	}
	wg.Wait()
	for _, record := range records(t, output.String()) {
		if record["request_id"] != record["trace_id"] || record["request_id"] != record["session_id"] {
			t.Fatalf("mixed IDs: %+v", record)
		}
	}
}

func TestInstalledFrameworkErrorsAndAccessAtErrorLevel(t *testing.T) {
	var output bytes.Buffer
	_ = ConfigureLogging(&output, "ERROR")
	originalApp, originalConfig := beego.BeeApp, beego.BConfig
	originalErrors := maps.Clone(beego.ErrorMaps)
	defer func() { beego.BeeApp, beego.BConfig, beego.ErrorMaps = originalApp, originalConfig, originalErrors }()
	cfg := *originalConfig
	cfg.WebConfig.AutoRender = false
	beego.BConfig = &cfg
	beego.BeeApp = beego.NewHttpServerWithCfg(&cfg)
	Install()
	beego.BeeApp.Handlers.Init()
	for _, test := range []struct {
		method string
		status int
	}{{"GET", 404}, {"UNKNOWN", 405}} {
		output.Reset()
		req := httptest.NewRequest(test.method, "/private-email@example.com?token=QUERY_SECRET", nil)
		res := httptest.NewRecorder()
		Middleware(beego.BeeApp.Handlers).ServeHTTP(res, req)
		if res.Code != test.status {
			t.Fatalf("framework status: %d, %s", res.Code, res.Body.String())
		}
		if !strings.HasPrefix(res.Result().Header.Get("Content-Type"), "application/json") {
			t.Fatal("framework error headers committed before JSON content type")
		}
		var body errorBody
		if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if !trace.Valid(body.Error.RequestID) || body.Error.RequestID != res.Header().Get(trace.RequestHeader) {
			t.Fatal("framework error not correlated")
		}
		if len(records(t, output.String())) != 1 || !strings.Contains(output.String(), `"msg":"http access"`) {
			t.Fatalf("access must remain enabled at ERROR: %s", output.String())
		}
		if strings.Contains(output.String(), "private-email") || strings.Contains(output.String(), "QUERY_SECRET") {
			t.Fatal("framework error leaked URL")
		}
	}
}

func TestBufferedWriterBoundaries(t *testing.T) {
	var output bytes.Buffer
	_ = ConfigureLogging(&output, "INFO")
	for _, size := range []int{0, 32, maxBufferedResponse, maxBufferedResponse + 1} {
		payload := strings.Repeat("x", size)
		res := httptest.NewRecorder()
		Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "preserved")
			w.WriteHeader(202)
			if n, err := w.Write([]byte(payload)); err != nil || n != size {
				t.Errorf("write: %d %v", n, err)
			}
		})).ServeHTTP(res, httptest.NewRequest("GET", "/", nil))
		if res.Code != 202 || res.Body.String() != payload || res.Result().Header.Get("X-Test") != "preserved" {
			t.Fatalf("buffer changed response at size %d", size)
		}
	}
}
