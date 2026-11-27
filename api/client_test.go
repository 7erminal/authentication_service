package api

import (
	"authentication_service/internal/observability"
	"authentication_service/internal/trace"
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func tracedContext() context.Context {
	return trace.WithIDs(context.Background(), trace.IDs{SessionID: "session_123", TraceID: "trace_123", RequestID: "request_123"})
}

func TestOutboundForwardingAndSensitiveData(t *testing.T) {
	var output bytes.Buffer
	_ = observability.ConfigureLogging(&output, "DEBUG")
	for _, internal := range []bool{true, false} {
		for _, withContext := range []bool{true, false} {
			for _, session := range []bool{true, false} {
				ctx := tracedContext()
				if !session {
					ids := trace.FromContext(ctx)
					ids.SessionID = ""
					ctx = trace.WithIDs(ctx, ids)
				}
				if !withContext {
					ctx = context.Background()
				}
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					for key, want := range map[string]string{trace.RequestHeader: "request_123", trace.TraceHeader: "trace_123", trace.SessionHeader: "session_123"} {
						if !internal || !withContext || key == trace.SessionHeader && !session {
							want = ""
						}
						if got := r.Header.Get(key); got != want {
							t.Errorf("%s = %q, want %q", key, got, want)
						}
					}
					if r.Header.Get("Authorization") != "Bearer AUTH_SECRET" {
						t.Error("auth behavior changed")
					}
					body, _ := io.ReadAll(r.Body)
					if !strings.Contains(string(body), "BODY_SECRET") {
						t.Error("JSON body missing")
					}
					w.WriteHeader(201)
					_, _ = w.Write([]byte("RESPONSE_SECRET"))
				}))
				request := NewRequest(server.URL, "/private-email@example.com?token=QUERY_SECRET", POST)
				request.HeaderField["Authorization"] = "Bearer AUTH_SECRET"
				request.HeaderField["Cookie"] = "session=COOKIE_SECRET"
				request.HeaderField[trace.RequestHeader] = "untrusted_123"
				request.InterfaceParams["password"] = "BODY_SECRET"
				client := Client{Request: request, Type_: "body", Internal: internal}
				res, err := client.SendRequestContext(ctx)
				if err != nil {
					t.Fatal(err)
				}
				_ = res.Body.Close()
				server.Close()
			}
		}
	}
	for _, secret := range []string{"SECRET", "private-email@example.com", "untrusted_123", "Authorization", "Cookie"} {
		if strings.Contains(output.String(), secret) {
			t.Fatalf("sensitive outbound log contains %s: %s", secret, output.String())
		}
	}
	if !strings.Contains(output.String(), `"request_id":"request_123"`) {
		t.Fatal("outbound logs are not correlated")
	}
}

func TestRedirectTraceBoundaries(t *testing.T) {
	external := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, key := range []string{trace.RequestHeader, trace.TraceHeader, trace.SessionHeader} {
			if r.Header.Get(key) != "" {
				t.Errorf("trace leaked to other origin: %s", key)
			}
		}
		w.WriteHeader(204)
	}))
	defer external.Close()
	internal := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(trace.RequestHeader) != "request_123" {
			t.Error("same-origin redirect lost IDs")
		}
		if r.URL.Path == "/start" {
			http.Redirect(w, r, "/next", 302)
		} else {
			http.Redirect(w, r, external.URL, 302)
		}
	}))
	defer internal.Close()
	client := Client{Request: NewRequest(internal.URL, "/start", GET), Internal: true}
	response, err := client.SendRequestContext(tracedContext())
	if err != nil {
		t.Fatal(err)
	}
	_ = response.Body.Close()
}

func TestOutboundCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { t.Error("canceled request reached upstream") }))
	defer server.Close()
	ctx, cancel := context.WithCancel(tracedContext())
	cancel()
	client := Client{Request: NewRequest(server.URL, "/", GET), Internal: true}
	start := time.Now()
	_, err := client.SendRequestContext(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation lost: %v", err)
	}
	if time.Since(start) > time.Second {
		t.Fatal("cancellation did not stop the call promptly")
	}
}
