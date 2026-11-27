package trace

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestHeaders(t *testing.T) {
	valid := []string{"abcdefgh", "a_B-012345", strings.Repeat("Z", 64)}
	invalid := []string{"", "short", strings.Repeat("a", 65), " username1", "username1 ", "abc.defgh", "abcdefgh\r\nX-Fake: yes", "abcdefgh\x00", "éabcdefgh", "abcdefgh,ijklmnop", "abcdefgh\t"}
	for _, value := range append(valid, invalid...) {
		t.Run(value, func(t *testing.T) {
			h := http.Header{}
			for _, key := range []string{RequestHeader, TraceHeader, SessionHeader} {
				h.Set(key, value)
			}
			ids := ReadHeaders(h)
			if Valid(value) {
				if ids != (IDs{value, value, value}) {
					t.Fatalf("valid IDs changed: %+v", ids)
				}
			} else {
				if !Valid(ids.RequestID) || ids.TraceID != ids.RequestID || ids.SessionID != "" {
					t.Fatalf("invalid fallback: %+v", ids)
				}
				if value != "" && ids.RequestID == value {
					t.Fatal("invalid ID retained")
				}
			}
		})
	}
	for _, key := range []string{RequestHeader, TraceHeader, SessionHeader} {
		h := http.Header{}
		h.Add(key, "abcdefgh")
		h.Add(key, "ijklmnop")
		ids := ReadHeaders(h)
		if ids.SessionID != "" || ids.TraceID != ids.RequestID || ids.RequestID == "abcdefgh" {
			t.Fatal("ambiguous header accepted")
		}
	}
}

func TestIndependentFallbacks(t *testing.T) {
	h := make(http.Header)
	h.Set(RequestHeader, "request_123")
	h.Set(TraceHeader, "bad.value")
	h.Set(SessionHeader, "session_123")
	if got := ReadHeaders(h); got != (IDs{"session_123", "request_123", "request_123"}) {
		t.Fatalf("%+v", got)
	}
	h.Del(RequestHeader)
	h.Set(TraceHeader, "trace_123")
	if got := ReadHeaders(h); !Valid(got.RequestID) || got.TraceID != "trace_123" {
		t.Fatalf("%+v", got)
	}
	if ReadHeaders(nil).RequestID == ReadHeaders(nil).RequestID {
		t.Fatal("generated IDs repeated")
	}
}

func TestForwardAndBackgroundMetadata(t *testing.T) {
	ids := IDs{"session_123", "trace_123", "request_123"}
	parent, cancel := context.WithCancel(WithIDs(context.Background(), ids))
	ctx := Background(parent)
	cancel()
	if ctx.Err() != nil || FromContext(ctx) != ids {
		t.Fatal("background metadata was lost or remained canceled")
	}
	// A future queue can serialize IDs as metadata and restore through WithIDs.
	data, _ := json.Marshal(FromContext(ctx))
	var restored IDs
	_ = json.Unmarshal(data, &restored)
	h := make(http.Header)
	Forward(WithIDs(context.Background(), restored), h)
	if h.Get(RequestHeader) != ids.RequestID || h.Get(TraceHeader) != ids.TraceID || h.Get(SessionHeader) != ids.SessionID {
		t.Fatal("forwarding changed IDs")
	}
	Forward(context.Background(), h)
	if len(h) != 0 {
		t.Fatal("stale caller headers retained")
	}
	if got := FromContext(WithIDs(context.Background(), IDs{"invalid.value", "bad", "invalid.value"})); got != (IDs{}) {
		t.Fatal("untrusted metadata not sanitized")
	}
}
