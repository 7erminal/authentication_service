// Package logger configures structured JSON logging (log/slog) for the
// service. Every record carries service, session_id, trace_id and request_id;
// the IDs are read from the context passed to the *Context logging functions
// (slog.InfoContext, slog.ErrorContext, ...).
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	"authentication_service/internal/trace"
)

// ServiceName is the value of the "service" field on every log line.
const ServiceName = "authentication_service"

// ParseLevel maps debug|info|warn|warning|error (case-insensitive) to a slog
// level. Unknown or empty values return info.
func ParseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// ResolveLevel returns the level from the LOG_LEVEL env var, falling back to
// confValue (app.conf "loglevel"), then info.
func ResolveLevel(confValue string) slog.Level {
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		return ParseLevel(v)
	}
	return ParseLevel(confValue)
}

// New builds a JSON logger writing to w at the given level.
func New(w io.Writer, level slog.Leveler) *slog.Logger {
	jh := slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level})
	return slog.New(&ContextHandler{inner: jh.WithAttrs([]slog.Attr{slog.String("service", ServiceName)})})
}

// Init installs the JSON logger as the slog default and routes Beego's
// logs.* calls through it. It returns the logger.
func Init(w io.Writer, level slog.Leveler) *slog.Logger {
	l := New(w, level)
	slog.SetDefault(l)
	installBeegoBridge(l)
	return l
}

// ContextHandler adds the tracing IDs from the context to every record.
type ContextHandler struct {
	inner slog.Handler
}

// Enabled implements slog.Handler.
func (h *ContextHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.inner.Enabled(ctx, l)
}

// Handle implements slog.Handler. The ID fields are always present (empty
// when unknown, e.g. at startup) so log queries can rely on them.
func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	ids, _ := trace.FromContext(ctx)
	out := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	out.AddAttrs(
		slog.String("session_id", ids.SessionID),
		slog.String("trace_id", ids.TraceID),
		slog.String("request_id", ids.RequestID),
	)
	r.Attrs(func(a slog.Attr) bool {
		out.AddAttrs(a)
		return true
	})
	return h.inner.Handle(ctx, out)
}

// WithAttrs implements slog.Handler.
func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{inner: h.inner.WithAttrs(attrs)}
}

// WithGroup implements slog.Handler.
func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{inner: h.inner.WithGroup(name)}
}
