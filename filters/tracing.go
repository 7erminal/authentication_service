// Package filters holds the Beego filters that implement request tracing,
// access logging, panic recovery and the JSON error envelope.
package filters

import (
	"log/slog"

	"authentication_service/internal/trace"

	"github.com/beego/beego/v2/server/web/context"
)

// Tracing is the BeforeRouter filter that reads and validates the gateway
// tracing headers, stores the IDs in the request context and echoes
// X-Request-Id / X-Trace-Id on the response.
//
// It is idempotent: the access-log filter chain (which wraps the whole
// router, including static files, 404/405 and CORS preflight) calls it first
// so early responses carry IDs too; the BeforeRouter invocation is then a
// cheap no-op. Requests without the headers are never rejected.
func Tracing(ctx *context.Context) {
	EnsureTraceIDs(ctx)
}

// EnsureTraceIDs resolves the request's IDs once and returns them.
func EnsureTraceIDs(ctx *context.Context) trace.IDs {
	if ids, ok := trace.FromContext(ctx.Request.Context()); ok {
		return ids
	}

	ids, invalid := trace.FromHeaders(ctx.Request.Header)
	rctx := trace.WithIDs(ctx.Request.Context(), ids)
	ctx.Request = ctx.Request.WithContext(rctx)

	// Normalise inbound headers so any code reading them directly only ever
	// sees sanitized values.
	h := ctx.Request.Header
	h.Set(trace.HeaderRequestID, ids.RequestID)
	h.Set(trace.HeaderTraceID, ids.TraceID)
	if ids.SessionID != "" {
		h.Set(trace.HeaderSessionID, ids.SessionID)
	} else {
		h.Del(trace.HeaderSessionID)
	}

	ctx.Output.Header(trace.HeaderRequestID, ids.RequestID)
	ctx.Output.Header(trace.HeaderTraceID, ids.TraceID)

	// Never log the rejected raw value, only which header was rejected.
	for _, name := range invalid {
		slog.WarnContext(rctx, "invalid tracing header discarded", "header", name)
	}
	return ids
}
