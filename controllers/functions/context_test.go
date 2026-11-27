package functions

import (
	"authentication_service/internal/observability"
	"authentication_service/internal/trace"
	"authentication_service/structs/requestsDTOs"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

func TestServiceHelpersKeepRequestContext(t *testing.T) {
	var output bytes.Buffer
	_ = observability.ConfigureLogging(&output, "DEBUG")
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Header.Get(trace.RequestHeader) != "request_123" || r.Header.Get(trace.TraceHeader) != "trace_123" || r.Header.Get(trace.SessionHeader) != "session_123" {
			t.Error("helper lost trace context")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"StatusCode":200,"StatusDesc":"OK","Result":{"UserId":"42","CustomerId":"42"},"Role":{"RoleId":"1"}}`))
	}))
	defer server.Close()
	original, _ := beego.AppConfig.String("customerBaseUrl")
	_ = beego.AppConfig.Set("customerBaseUrl", server.URL)
	defer beego.AppConfig.Set("customerBaseUrl", original)
	ctx := trace.WithIDs(context.Background(), trace.IDs{SessionID: "session_123", TraceID: "trace_123", RequestID: "request_123"})
	c := &beego.Controller{Ctx: beecontext.NewContext()}
	c.Ctx.Request = httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	if _, err := GetUser(c, requestsDTOs.GetUserRequest{UserId: "42"}); err != nil {
		t.Fatal(err)
	}
	if _, err := GetUserWithUsername(c, requestsDTOs.GetUserWithUsernameRequest{Username: "PRIVATE_EMAIL_SECRET"}); err != nil {
		t.Fatal(err)
	}
	if _, err := GetRole(c, "1"); err != nil {
		t.Fatal(err)
	}
	if _, err := GetCustomer(c, requestsDTOs.GetCustomerRequest{CustomerId: "42"}); err != nil {
		t.Fatal(err)
	}
	if _, err := GetCustomerWithUsername(c, requestsDTOs.GetCustomerWithUsernameRequest{Username: "PRIVATE_EMAIL_SECRET"}); err != nil {
		t.Fatal(err)
	}
	if _, err := UpdateUserPassword(c, requestsDTOs.UpdateUserPasswordRequest{UserId: "42", NewPassword: "PASSWORD_SECRET"}); err != nil {
		t.Fatal(err)
	}
	token, _, _, err := CreateAccessToken("42", "PRIVATE_EMAIL_SECRET", "1", "role", nil, "1", "1")
	if err != nil {
		t.Fatal(err)
	}
	if user, err := GetUserFromBearerTokenContext(ctx, "Bearer "+token); err != nil || user.UserId != "42" {
		t.Fatalf("bearer lookup: %v", err)
	}
	if calls != 7 {
		t.Fatalf("calls: %d", calls)
	}
	if strings.Contains(output.String(), "SECRET") || strings.Contains(output.String(), token) {
		t.Fatal("sensitive helper logs")
	}
}
