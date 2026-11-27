package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"authentication_service/internal/trace"

	"github.com/beego/beego/v2/core/logs"
	"github.com/go-sql-driver/mysql"
)

func decodeLines(t *testing.T, buf *bytes.Buffer) []map[string]interface{} {
	t.Helper()
	var out []map[string]interface{}
	for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		if line == "" {
			continue
		}
		m := map[string]interface{}{}
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("log line is not JSON: %q: %v", line, err)
		}
		out = append(out, m)
	}
	return out
}

func TestEveryLineHasServiceAndIDs(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, slog.LevelDebug)
	ctx := trace.WithIDs(context.Background(), trace.IDs{SessionID: "session_12345", TraceID: "trace_12345", RequestID: "request_12345"})

	l.InfoContext(ctx, "with ids", "user_id", 42)
	l.Info("without ctx")

	lines := decodeLines(t, &buf)
	if len(lines) != 2 {
		t.Fatalf("want 2 lines, got %d", len(lines))
	}
	first := lines[0]
	for k, want := range map[string]interface{}{
		"service": ServiceName, "session_id": "session_12345", "trace_id": "trace_12345",
		"request_id": "request_12345", "msg": "with ids", "user_id": float64(42),
	} {
		if first[k] != want {
			t.Fatalf("field %s = %v, want %v (line %v)", k, first[k], want, first)
		}
	}
	second := lines[1]
	for _, k := range []string{"service", "session_id", "trace_id", "request_id"} {
		if _, ok := second[k]; !ok {
			t.Fatalf("field %s missing on line without context: %v", k, second)
		}
	}
}

func TestLevel(t *testing.T) {
	for in, want := range map[string]slog.Level{
		"debug": slog.LevelDebug, "INFO": slog.LevelInfo, "warn": slog.LevelWarn,
		"warning": slog.LevelWarn, "error": slog.LevelError, "": slog.LevelInfo, "nonsense": slog.LevelInfo,
	} {
		if got := ParseLevel(in); got != want {
			t.Fatalf("ParseLevel(%q) = %v, want %v", in, got, want)
		}
	}

	t.Setenv("LOG_LEVEL", "error")
	if got := ResolveLevel("debug"); got != slog.LevelError {
		t.Fatalf("env must override config, got %v", got)
	}
	t.Setenv("LOG_LEVEL", "")
	if got := ResolveLevel("warn"); got != slog.LevelWarn {
		t.Fatalf("config fallback, got %v", got)
	}

	var buf bytes.Buffer
	l := New(&buf, slog.LevelWarn)
	l.Info("dropped")
	l.Warn("kept")
	if lines := decodeLines(t, &buf); len(lines) != 1 || lines[0]["msg"] != "kept" {
		t.Fatalf("level filtering failed: %v", lines)
	}
}

func TestBeegoBridge(t *testing.T) {
	var buf bytes.Buffer
	Init(&buf, slog.LevelDebug)

	logs.Info("legacy %s", "message")
	logs.Error("legacy error")
	// Beego's dev-mode per-request line (duplicate of our access log).
	logs.Debug("|      127.0.0.1|\x1b[97;42m 200 \x1b[0m|     1.234ms|   match|\x1b[97;44m POST    \x1b[0m /v1/auth/login/token   r:/v1/auth/login/token")
	logs.Debug("|      127.0.0.1|\x1b[97;43m 404 \x1b[0m|      51.2µs| nomatch|\x1b[97;44m GET     \x1b[0m /nope")

	lines := decodeLines(t, &buf)
	if len(lines) != 2 {
		t.Fatalf("want 2 bridged lines (dev router lines dropped), got %d: %v", len(lines), lines)
	}
	if lines[0]["msg"] != "legacy message" || lines[0]["logger"] != "beego" || lines[0]["level"] != "INFO" {
		t.Fatalf("unexpected bridged line %v", lines[0])
	}
	if lines[1]["level"] != "ERROR" || lines[1]["service"] != ServiceName {
		t.Fatalf("unexpected bridged line %v", lines[1])
	}
}

func TestErrHidesMySQLValues(t *testing.T) {
	me := &mysql.MySQLError{Number: 1062, Message: "Duplicate entry 'SECRET_TOKEN_VALUE' for key 'token'"}
	attr := Err(fmt.Errorf("insert: %w", me))
	if strings.Contains(attr.Value.String(), "SECRET_TOKEN_VALUE") || attr.Value.String() != "mysql error 1062" {
		t.Fatalf("mysql error value leaked: %q", attr.Value.String())
	}
	if got := Err(errors.New("<QuerySeter> no row found")).Value.String(); got != "<QuerySeter> no row found" {
		t.Fatalf("got %q", got)
	}
}
