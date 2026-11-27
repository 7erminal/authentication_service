// Package trace carries the gateway tracing IDs (session, trace, request)
// through this service.
//
// LOCAL COPY: the gateways do not yet publish a shared trace module. This
// package mirrors the API we expect that module to expose so it can be
// swapped for the shared one later with an import-path change only.
//
// The IDs are for observability (logging/correlation) ONLY. They must never
// be used for authentication, authorization, idempotency or as storage keys.
package trace

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"
)

// Header names of the gateway contract.
const (
	HeaderSessionID = "X-Session-Id"
	HeaderTraceID   = "X-Trace-Id"
	HeaderRequestID = "X-Request-Id"
)

// Metadata keys used when IDs travel in message/job metadata instead of
// HTTP headers. They match the structured-log field names.
const (
	MetaSessionID = "session_id"
	MetaTraceID   = "trace_id"
	MetaRequestID = "request_id"
)

const (
	minLen = 8
	maxLen = 64
)

// IDs is the set of correlation IDs attached to a unit of work.
type IDs struct {
	SessionID string
	TraceID   string
	RequestID string
}

type ctxKey struct{}

// Valid reports whether v satisfies the contract: 8-64 characters drawn from
// [A-Za-z0-9_-]. Anything else must be discarded, never logged.
func Valid(v string) bool {
	if len(v) < minLen || len(v) > maxLen {
		return false
	}
	for i := 0; i < len(v); i++ {
		c := v[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9', c == '_', c == '-':
		default:
			return false
		}
	}
	return true
}

// Sanitize returns v if it is valid and "" otherwise.
func Sanitize(v string) string {
	if Valid(v) {
		return v
	}
	return ""
}

// NewRequestID returns a fresh 32-character hex ID.
func NewRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand failing is practically impossible; fall back to a
		// time-based value that still satisfies the contract.
		return "t" + strconv.FormatInt(time.Now().UnixNano(), 36) + "0000000"
	}
	return hex.EncodeToString(b)
}

// Resolve applies the contract to raw values: invalid values are dropped, a
// missing request ID is generated and a missing trace ID falls back to the
// request ID. It returns the names (header names) of values that were present
// but rejected so callers can log THAT they were rejected without logging the
// raw value.
func Resolve(rawSessionID, rawTraceID, rawRequestID string) (IDs, []string) {
	var invalid []string
	check := func(name, raw string) string {
		if raw == "" {
			return ""
		}
		if v := Sanitize(raw); v != "" {
			return v
		}
		invalid = append(invalid, name)
		return ""
	}
	ids := IDs{
		SessionID: check(HeaderSessionID, rawSessionID),
		TraceID:   check(HeaderTraceID, rawTraceID),
		RequestID: check(HeaderRequestID, rawRequestID),
	}
	if ids.RequestID == "" {
		ids.RequestID = NewRequestID()
	}
	if ids.TraceID == "" {
		ids.TraceID = ids.RequestID
	}
	return ids, invalid
}

// FromHeaders reads and validates the tracing headers. A header sent more than
// once is ambiguous and treated as invalid.
func FromHeaders(h http.Header) (IDs, []string) {
	var dup []string
	get := func(name string) string {
		vals := h.Values(name)
		switch len(vals) {
		case 0:
			return ""
		case 1:
			return vals[0]
		default:
			dup = append(dup, name)
			return ""
		}
	}
	ids, invalid := Resolve(get(HeaderSessionID), get(HeaderTraceID), get(HeaderRequestID))
	return ids, append(dup, invalid...)
}

// WithIDs stores ids in ctx.
func WithIDs(ctx context.Context, ids IDs) context.Context {
	return context.WithValue(ctx, ctxKey{}, ids)
}

// FromContext returns the IDs stored in ctx, if any.
func FromContext(ctx context.Context) (IDs, bool) {
	if ctx == nil {
		return IDs{}, false
	}
	ids, ok := ctx.Value(ctxKey{}).(IDs)
	return ids, ok
}

// RequestID returns the request ID stored in ctx or "".
func RequestID(ctx context.Context) string {
	ids, _ := FromContext(ctx)
	return ids.RequestID
}

// InjectHeaders sets the SAME IDs from ctx on an outbound request's headers.
// It never generates new IDs.
func InjectHeaders(ctx context.Context, h http.Header) {
	ids, ok := FromContext(ctx)
	if !ok {
		return
	}
	set := func(name, v string) {
		if v != "" {
			h.Set(name, v)
		}
	}
	set(HeaderRequestID, ids.RequestID)
	set(HeaderTraceID, ids.TraceID)
	set(HeaderSessionID, ids.SessionID)
}

// ToMetadata returns the IDs from ctx as message/job metadata.
func ToMetadata(ctx context.Context) map[string]string {
	md := map[string]string{}
	ids, ok := FromContext(ctx)
	if !ok {
		return md
	}
	if ids.RequestID != "" {
		md[MetaRequestID] = ids.RequestID
	}
	if ids.TraceID != "" {
		md[MetaTraceID] = ids.TraceID
	}
	if ids.SessionID != "" {
		md[MetaSessionID] = ids.SessionID
	}
	return md
}

// FromMetadata restores IDs from message/job metadata into ctx (consumer
// side), applying the same validation and fallbacks as inbound headers.
func FromMetadata(ctx context.Context, md map[string]string) (context.Context, []string) {
	ids, invalid := Resolve(md[MetaSessionID], md[MetaTraceID], md[MetaRequestID])
	return WithIDs(ctx, ids), invalid
}

// Detach returns a context that keeps the IDs (and other values) of ctx but is
// not cancelled when the request finishes. Use it for in-process background
// work started from a request.
func Detach(ctx context.Context) context.Context {
	return context.WithoutCancel(ctx)
}
