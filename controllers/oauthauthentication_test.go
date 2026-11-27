package controllers

import (
	"authentication_service/internal/observability"
	"authentication_service/internal/trace"
	"authentication_service/structs/responsesDTOs"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	beecontext "github.com/beego/beego/v2/server/web/context"
)

func TestGoogleOAuthEnvironmentConfiguration(t *testing.T) {
	for _, test := range []struct {
		name, clientID, clientSecret string
		callback                     bool
		status                       int
	}{
		{"authorization missing ID", "", "", false, 500},
		{"callback missing credentials", "", "", true, 500},
		{"callback missing secret", "test-client", "", true, 500},
		{"callback missing ID", "", "test-secret", true, 500},
		{"authorization configured", "test-client&literal", "test-secret", false, 200},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("GOOGLE_OAUTH_CLIENT_ID", test.clientID)
			t.Setenv("GOOGLE_OAUTH_CLIENT_SECRET", test.clientSecret)
			var output bytes.Buffer
			if err := observability.ConfigureLogging(&output, "DEBUG"); err != nil {
				t.Fatal(err)
			}
			request := httptest.NewRequest("GET", "/?third-party=google", nil)
			request.Header.Set(trace.RequestHeader, "request_123")
			response := httptest.NewRecorder()
			observability.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := beecontext.NewContext()
				ctx.Reset(w, r)
				c := &OAuthAuthenticationController{}
				c.Init(ctx, "OAuthAuthenticationController", "test", nil)
				if test.callback {
					c.OAuthGoogle()
				} else {
					c.OAuthThirdPartyLogin()
				}
			})).ServeHTTP(response, request)
			if response.Code != test.status {
				t.Fatalf("status = %d, want %d", response.Code, test.status)
			}
			if test.status == 500 {
				if !strings.Contains(response.Body.String(), `"request_id":"request_123"`) || !strings.Contains(response.Body.String(), `"code":"internal_error"`) {
					t.Fatal("configuration error is not safely correlated")
				}
			} else {
				var body responsesDTOs.StringResponseDTO
				if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}
				location, err := url.Parse(body.Value)
				if err != nil || location.Query().Get("client_id") != test.clientID {
					t.Fatal("authorization URL did not use and escape the configured ID")
				}
				if strings.Contains(body.Value, test.clientSecret) {
					t.Fatal("authorization URL leaked the client secret")
				}
			}
			if strings.Contains(output.String(), "test-client") || strings.Contains(output.String(), "test-secret") {
				t.Fatal("OAuth configuration leaked into logs")
			}
		})
	}
}
