package controllers

import (
	"authentication_service/internal/observability"
	"authentication_service/internal/trace"
	"authentication_service/structs/responsesDTOs"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

func TestLegacyErrorCompatibility(t *testing.T) {
	for _, test := range []struct {
		name       string
		value      interface{}
		status     int
		correlated bool
	}{
		{"business error with HTTP 200", responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Incorrect password"}, 200, true},
		{"HTTP error", responsesDTOs.UserResponseDTO{StatusCode: 401, StatusDesc: "Invalid access token"}, 401, true},
		{"invalid token object", responsesDTOs.UserTokenResponseDTO{IsValid: false}, 200, true},
		{"success", responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "OK", StatusDesc: "Successful"}, 200, false},
		{"legacy string", "original error", 200, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			_ = observability.ConfigureLogging(&output, "DEBUG")
			res := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set(trace.RequestHeader, "request_123")
			observability.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := beecontext.NewContext()
				ctx.Reset(w, r)
				observability.BeforeRouter(ctx)
				controller := &tracedController{}
				controller.Init(ctx, "test", "test", nil)
				controller.Data["json"] = test.value
				controller.Ctx.Output.SetStatus(test.status)
				if err := controller.ServeJSON(); err != nil {
					t.Fatal(err)
				}
			})).ServeHTTP(res, req)
			if res.Code != test.status {
				t.Fatalf("status changed: %d", res.Code)
			}
			original, _ := json.Marshal(test.value)
			var want, got interface{}
			_ = json.Unmarshal(original, &want)
			if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if test.correlated {
				object := got.(map[string]interface{})
				if object["request_id"] != "request_123" {
					t.Fatal("error ID missing")
				}
				delete(object, "request_id")
			}
			gotJSON, _ := json.Marshal(got)
			wantJSON, _ := json.Marshal(want)
			if string(gotJSON) != string(wantJSON) {
				t.Fatalf("legacy contract changed: %s != %s", gotJSON, wantJSON)
			}
		})
	}
}

func TestRespondError(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(trace.RequestHeader, "request_123")
	observability.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := beecontext.NewContext()
		ctx.Reset(w, r)
		c := &tracedController{}
		c.Init(ctx, "test", "test", nil)
		c.respondError(403, "forbidden", "Not permitted")
	})).ServeHTTP(res, req)
	if res.Code != 403 || !strings.Contains(res.Body.String(), `"request_id":"request_123"`) {
		t.Fatalf("bad error: %s", res.Body.String())
	}
}

// Use a real application endpoint and Beego controller initialization to ensure
// embedding tracedController does not break route dispatch or ServeJSON.
func TestTokenVerificationEndpoint(t *testing.T) {
	var output bytes.Buffer
	_ = observability.ConfigureLogging(&output, "DEBUG")
	cfg := *beego.BConfig
	cfg.CopyRequestBody = true
	cfg.WebConfig.AutoRender = false
	cfg.RecoverFunc = observability.Recover
	router := beego.NewControllerRegisterWithCfg(&cfg)
	router.InsertFilter("*", beego.BeforeRouter, observability.BeforeRouter)
	router.InsertFilter("*", beego.BeforeExec, observability.BeforeExec)
	router.Add("/verify-token", &AuthenticationController{}, beego.WithRouterMethods(&AuthenticationController{}, "post:VerifyToken"))
	router.Init()
	req := httptest.NewRequest("POST", "/verify-token?token=QUERY_SECRET", strings.NewReader(`{"Value":"BODY_SECRET"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(trace.RequestHeader, "request_123")
	res := httptest.NewRecorder()
	observability.Middleware(router).ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fatalf("status changed: %d %s", res.Code, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), `"request_id"`) {
		t.Fatalf("error correlation missing: %s", res.Body.String())
	}
	if strings.Contains(output.String(), "SECRET") {
		t.Fatal("endpoint leaked input")
	}
	if strings.Count(output.String(), `"msg":"http access"`) != 1 {
		t.Fatal("incorrect number of access logs")
	}
}
