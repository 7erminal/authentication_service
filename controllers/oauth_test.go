package controllers

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func thirdPartyLoginURL(t *testing.T) (string, string) {
	t.Helper()
	h, logs := newTestApp(t)
	req := httptest.NewRequest("GET", "/v1/oauth/third-party/login?third-party=google", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	var body struct {
		StatusCode int
		Value      string
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.StatusCode != 200 {
		t.Fatalf("unexpected response %d %q", rec.Code, rec.Body.String())
	}
	return body.Value, logs.String()
}

func TestGoogleClientIDComesFromEnv(t *testing.T) {
	t.Setenv(EnvGoogleClientID, "test-client.apps.example")
	t.Setenv(EnvGoogleClientSecret, "test-secret-value")
	if !GoogleOAuthConfigured() {
		t.Fatal("expected configured")
	}

	u, logs := thirdPartyLoginURL(t)
	if !strings.Contains(u, "client_id=test-client.apps.example&") {
		t.Fatalf("client id not taken from env: %s", u)
	}
	if strings.Contains(logs, "test-secret-value") || strings.Contains(logs, "test-client.apps.example") {
		t.Fatalf("credentials logged:\n%s", logs)
	}
}

func TestGoogleOAuthMissingEnvLogsError(t *testing.T) {
	t.Setenv(EnvGoogleClientID, "")
	t.Setenv(EnvGoogleClientSecret, "")
	if GoogleOAuthConfigured() {
		t.Fatal("expected not configured")
	}
	_, logs := thirdPartyLoginURL(t)
	if !strings.Contains(logs, "google oauth client id not configured") {
		t.Fatalf("missing config error not logged:\n%s", logs)
	}
}
