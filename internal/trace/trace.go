// Package trace is a local compatibility boundary until the gateway publishes
// its shared trace module. These IDs are observability metadata, never auth keys.
package trace

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
)

const (
	SessionHeader = "X-Session-Id"
	TraceHeader   = "X-Trace-Id"
	RequestHeader = "X-Request-Id"
	Service       = "authentication_service"
)

type IDs struct {
	SessionID string `json:"session_id"`
	TraceID   string `json:"trace_id"`
	RequestID string `json:"request_id"`
}

type contextKey struct{}

// Valid rejects the entire value; it never trims or repairs invalid input.
func Valid(value string) bool {
	if len(value) < 8 || len(value) > 64 {
		return false
	}
	for _, c := range []byte(value) {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c == '-') {
			return false
		}
	}
	return true
}

func header(h http.Header, key string) string {
	values := h.Values(key)
	if len(values) != 1 || !Valid(values[0]) {
		return ""
	}
	return values[0]
}

func ReadHeaders(h http.Header) IDs {
	ids := IDs{header(h, SessionHeader), header(h, TraceHeader), header(h, RequestHeader)}
	if ids.RequestID == "" {
		var b [16]byte
		// crypto/rand failure is fatal in modern Go; never use predictable fallbacks.
		if _, err := rand.Read(b[:]); err != nil {
			panic("request ID entropy unavailable")
		}
		ids.RequestID = hex.EncodeToString(b[:])
	}
	if ids.TraceID == "" {
		ids.TraceID = ids.RequestID
	}
	return ids
}

// WithIDs also validates restored background-job metadata, without generating IDs.
func WithIDs(ctx context.Context, ids IDs) context.Context {
	if !Valid(ids.SessionID) {
		ids.SessionID = ""
	}
	if !Valid(ids.TraceID) {
		ids.TraceID = ""
	}
	if !Valid(ids.RequestID) {
		ids.RequestID = ""
	}
	return context.WithValue(ctx, contextKey{}, ids)
}

func FromContext(ctx context.Context) IDs {
	if ctx == nil {
		return IDs{}
	}
	ids, _ := ctx.Value(contextKey{}).(IDs)
	return ids
}

func Logger(ctx context.Context) *slog.Logger {
	return LoggerFrom(ctx, slog.Default())
}

// LoggerFrom attaches the same correlation fields to a dedicated sink.
func LoggerFrom(ctx context.Context, root *slog.Logger) *slog.Logger {
	ids := FromContext(ctx)
	return root.With("service", Service, "session_id", ids.SessionID, "trace_id", ids.TraceID, "request_id", ids.RequestID)
}

// Forward replaces caller-supplied values; it never creates a new request ID.
func Forward(ctx context.Context, h http.Header) {
	ids := FromContext(ctx)
	for key, value := range map[string]string{SessionHeader: ids.SessionID, TraceHeader: ids.TraceID, RequestHeader: ids.RequestID} {
		h.Del(key)
		if Valid(value) {
			h.Set(key, value)
		}
	}
}

// Background detaches cancellation and arbitrary request values, retaining only
// validated trace metadata. A worker should impose its own timeout.
func Background(ctx context.Context) context.Context {
	return WithIDs(context.Background(), FromContext(ctx))
}
