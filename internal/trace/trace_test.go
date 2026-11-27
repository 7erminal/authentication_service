package trace

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestValid(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"min length", "abcdefgh", true},
		{"max length", strings.Repeat("a", 64), true},
		{"uuid", "3f2b8c1e-9d4a-4b7e-8f21-0c6d5e4a3b2c", true},
		{"underscore", "req_1234_ABCD", true},
		{"empty", "", false},
		{"too short", "abcdefg", false},
		{"too long", strings.Repeat("a", 65), false},
		{"space", "abcd efgh", false},
		{"CRLF injection", "abcdefgh\r\nX-Evil: 1", false},
		{"newline log injection", "abcdefgh\n{\"level\":\"ERROR\"}", false},
		{"quote", `abcdefgh"`, false},
		{"html", "<script>alert(1)</script>", false},
		{"sql", "' OR 1=1 --abcd", false},
		{"unicode", "ááááááááá", false},
		{"null byte", "abcdefgh\x00", false},
		{"dot", "abcd.efgh", false},
		{"slash", "abcd/efgh", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Valid(c.in); got != c.want {
				t.Fatalf("Valid(%q) = %v, want %v", c.in, got, c.want)
			}
			if got := Sanitize(c.in); (got != "") != c.want {
				t.Fatalf("Sanitize(%q) = %q", c.in, got)
			}
		})
	}
}

func TestNewRequestID(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 1000; i++ {
		id := NewRequestID()
		if !Valid(id) {
			t.Fatalf("generated invalid id %q", id)
		}
		if seen[id] {
			t.Fatalf("duplicate id %q", id)
		}
		seen[id] = true
	}
}

func TestResolve(t *testing.T) {
	t.Run("all valid", func(t *testing.T) {
		ids, invalid := Resolve("session_12345", "trace_12345", "request_12345")
		if ids != (IDs{SessionID: "session_12345", TraceID: "trace_12345", RequestID: "request_12345"}) {
			t.Fatalf("unexpected ids %+v", ids)
		}
		if len(invalid) != 0 {
			t.Fatalf("unexpected invalid %v", invalid)
		}
	})
	t.Run("all missing", func(t *testing.T) {
		ids, invalid := Resolve("", "", "")
		if !Valid(ids.RequestID) {
			t.Fatalf("request id not generated: %+v", ids)
		}
		if ids.TraceID != ids.RequestID {
			t.Fatalf("trace id should fall back to request id: %+v", ids)
		}
		if ids.SessionID != "" {
			t.Fatalf("session id must not be invented: %+v", ids)
		}
		if len(invalid) != 0 {
			t.Fatalf("missing headers are not invalid: %v", invalid)
		}
	})
	t.Run("trace missing uses request id", func(t *testing.T) {
		ids, _ := Resolve("", "", "request_12345")
		if ids.TraceID != "request_12345" {
			t.Fatalf("got %+v", ids)
		}
	})
	t.Run("invalid values are discarded and reported by name only", func(t *testing.T) {
		ids, invalid := Resolve("bad session!", "x", "evil\r\nheader")
		if ids.SessionID != "" || !Valid(ids.RequestID) || ids.TraceID != ids.RequestID {
			t.Fatalf("got %+v", ids)
		}
		want := []string{HeaderSessionID, HeaderTraceID, HeaderRequestID}
		if strings.Join(invalid, ",") != strings.Join(want, ",") {
			t.Fatalf("invalid = %v, want %v", invalid, want)
		}
		for _, v := range invalid {
			if strings.Contains(v, "evil") || strings.Contains(v, "bad") {
				t.Fatalf("raw value leaked in %q", v)
			}
		}
	})
}

func TestFromHeadersDuplicateIsInvalid(t *testing.T) {
	h := http.Header{}
	h.Add(HeaderRequestID, "request_aaaaaaaa")
	h.Add(HeaderRequestID, "request_bbbbbbbb")
	ids, invalid := FromHeaders(h)
	if ids.RequestID == "request_aaaaaaaa" || ids.RequestID == "request_bbbbbbbb" {
		t.Fatalf("ambiguous duplicate header must not be trusted: %+v", ids)
	}
	if len(invalid) != 1 || invalid[0] != HeaderRequestID {
		t.Fatalf("invalid = %v", invalid)
	}
}

func TestInjectHeadersForwardsSameIDs(t *testing.T) {
	ctx := WithIDs(context.Background(), IDs{SessionID: "session_12345", TraceID: "trace_12345", RequestID: "request_12345"})
	h := http.Header{}
	InjectHeaders(ctx, h)
	if h.Get(HeaderRequestID) != "request_12345" || h.Get(HeaderTraceID) != "trace_12345" || h.Get(HeaderSessionID) != "session_12345" {
		t.Fatalf("headers = %v", h)
	}

	// No IDs in context: nothing is invented.
	h2 := http.Header{}
	InjectHeaders(context.Background(), h2)
	if len(h2) != 0 {
		t.Fatalf("expected no headers, got %v", h2)
	}

	// Empty session ID is not sent.
	ctx3 := WithIDs(context.Background(), IDs{TraceID: "trace_12345", RequestID: "request_12345"})
	h3 := http.Header{}
	InjectHeaders(ctx3, h3)
	if _, ok := h3[HeaderSessionID]; ok {
		t.Fatalf("empty session id must not be forwarded: %v", h3)
	}
}

func TestMetadataRoundTrip(t *testing.T) {
	ids := IDs{SessionID: "session_12345", TraceID: "trace_12345", RequestID: "request_12345"}
	md := ToMetadata(WithIDs(context.Background(), ids))
	ctx, invalid := FromMetadata(context.Background(), md)
	got, ok := FromContext(ctx)
	if !ok || got != ids || len(invalid) != 0 {
		t.Fatalf("round trip got %+v ok=%v invalid=%v", got, ok, invalid)
	}

	// Tampered metadata is sanitized on the consumer side.
	ctx, invalid = FromMetadata(context.Background(), map[string]string{MetaRequestID: "bad\nvalue", MetaTraceID: "trace_12345"})
	got, _ = FromContext(ctx)
	if !Valid(got.RequestID) || got.RequestID == "bad\nvalue" || got.TraceID != "trace_12345" || len(invalid) != 1 {
		t.Fatalf("got %+v invalid=%v", got, invalid)
	}
}

func TestDetachKeepsIDsButNotCancellation(t *testing.T) {
	parent, cancel := context.WithCancel(WithIDs(context.Background(), IDs{RequestID: "request_12345", TraceID: "request_12345"}))
	detached := Detach(parent)
	cancel()
	select {
	case <-detached.Done():
		t.Fatal("detached context must not be cancelled with the request")
	case <-time.After(10 * time.Millisecond):
	}
	if RequestID(detached) != "request_12345" {
		t.Fatalf("detached context lost ids")
	}
}
