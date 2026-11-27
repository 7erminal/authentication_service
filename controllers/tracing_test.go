package controllers

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"authentication_service/controllers/functions"
	"authentication_service/filters"
	"authentication_service/internal/logger"
	"authentication_service/internal/trace"
	"authentication_service/models"

	beego "github.com/beego/beego/v2/server/web"
)

const testReqID = "req_controller_test_01"

func newTestApp(t *testing.T) (http.Handler, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer
	logger.Init(&buf, slog.LevelDebug)

	cfg := *beego.BConfig
	cfg.RunMode = beego.DEV
	cfg.CopyRequestBody = true
	// AutoRender is left at its default (true) on purpose: error pages
	// must not depend on conf/app.conf's "autorender = false".
	app := beego.NewHttpServerWithCfg(&cfg)
	filters.Install(app)
	beego.ErrorController(&ErrorController{})

	app.Router("/v1/auth/customer-token/check", &AuthenticationController{}, "post:CheckCustomerTokenExpiry")
	app.Router("/v1/auth/token/verify", &AuthenticationController{}, "post:VerifyToken")
	app.Router("/v1/oauth/third-party/login", &OAuthAuthenticationController{}, "get:OAuthThirdPartyLogin")
	app.Handlers.Init()
	return app.Handlers, &buf
}

func post(h http.Handler, path, body string, hdr map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range hdr {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestLegacyErrorBodyGetsRequestIDOnly(t *testing.T) {
	h, logs := newTestApp(t)
	rec := post(h, "/v1/auth/customer-token/check", `{"Value":"not-a-jwt"}`, map[string]string{trace.HeaderRequestID: testReqID})

	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("body %q", rec.Body.String())
	}
	// Existing contract unchanged...
	if rec.Code != 200 || body["StatusCode"] != float64(703) || body["StatusDesc"] != "Error validating token" {
		t.Fatalf("legacy contract changed: %d %v", rec.Code, body)
	}
	// ...plus the additive request_id.
	if body["request_id"] != testReqID || rec.Header().Get(trace.HeaderRequestID) != testReqID {
		t.Fatalf("request_id missing: %v / %v", body, rec.Header())
	}
	if strings.Contains(logs.String(), "not-a-jwt") {
		t.Fatalf("token value logged:\n%s", logs.String())
	}
}

func TestLegacySuccessBodyUnchanged(t *testing.T) {
	h, logs := newTestApp(t)
	cust := models.Customers{CustomerId: 77, Email: "jane@example.com", PhoneNumber: "+233200000000", FullName: "Jane Secret"}
	token, _, err := functions.CreateCustomerAccessToken("jane_user", cust)
	if err != nil {
		t.Fatal(err)
	}
	rec := post(h, "/v1/auth/customer-token/check", `{"Value":"`+token+`"}`, map[string]string{trace.HeaderRequestID: testReqID})

	var body map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if body["StatusCode"] != float64(200) {
		t.Fatalf("unexpected body %v", body)
	}
	if _, ok := body["request_id"]; ok {
		t.Fatalf("success bodies must stay unchanged: %v", body)
	}
	out := logs.String()
	for _, s := range []string{token, "jane@example.com", "+233200000000", "Jane Secret", "jane_user"} {
		if strings.Contains(out, s) {
			t.Fatalf("sensitive value %q logged:\n%s", s, out)
		}
	}
	// Business ID is logged with the tracing IDs on the same line.
	found := false
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		var m map[string]interface{}
		_ = json.Unmarshal([]byte(line), &m)
		if m["msg"] == "customer token is valid" {
			found = m["customer_id"] == "77" && m["request_id"] == testReqID
		}
	}
	if !found {
		t.Fatalf("expected business id + request id on one line:\n%s", out)
	}
}

func TestVerifyTokenRejectsMissingAuthWithoutLoggingIt(t *testing.T) {
	h, logs := newTestApp(t)
	rec := post(h, "/v1/auth/token/verify", `{"Token":"REFRESH_SECRET"}`, map[string]string{
		"Authorization": "Bearer",
	})
	var body map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if rec.Code != 401 || body["StatusCode"] != float64(605) || body["request_id"] != rec.Header().Get(trace.HeaderRequestID) {
		t.Fatalf("got %d %v", rec.Code, body)
	}
	if strings.Contains(logs.String(), "REFRESH_SECRET") {
		t.Fatalf("token logged")
	}
}

func TestUnknownRouteReturnsJSONError(t *testing.T) {
	beego.BConfig.WebConfig.AutoRender = true // beego default; executeError reads the global
	h, logs := newTestApp(t)
	rec := post(h, "/v1/nope", "", map[string]string{trace.HeaderRequestID: testReqID})
	if rec.Code != 404 {
		t.Fatalf("status %d", rec.Code)
	}
	var eb filters.ErrorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &eb); err != nil {
		t.Fatalf("404 body not JSON: %q", rec.Body.String())
	}
	if eb.Error.Code != "not_found" || eb.Error.RequestID != testReqID {
		t.Fatalf("body %+v", eb)
	}
	if !strings.HasPrefix(rec.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("content-type %q", rec.Header().Get("Content-Type"))
	}
	if strings.Contains(logs.String(), "panic recovered") {
		t.Fatalf("error page rendering panicked:\n%s", logs.String())
	}
	if n := strings.Count(logs.String(), `"msg":"request completed"`); n != 1 {
		t.Fatalf("want 1 access line, got %d", n)
	}
}
