package filters

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

// RecoverPanic replaces Beego's default RecoverFunc (BConfig.RecoverFunc).
// It logs the panic with the request IDs and stack trace and returns a JSON
// 500 without any stack trace, regardless of run mode.
func RecoverPanic(ctx *context.Context, cfg *beego.Config) {
	rec := recover()
	if rec == nil {
		return
	}
	if rec == beego.ErrAbort {
		// Controller.StopRun / CustomAbort: response already written.
		return
	}
	if cfg != nil && !cfg.RecoverPanic {
		panic(rec)
	}
	// c.Abort("404") panics with the status code string: answer with the
	// matching error envelope, this is not a crash.
	if s, ok := rec.(string); ok {
		if code, err := strconv.Atoi(s); err == nil && code >= 400 && code <= 599 {
			if _, registered := beego.ErrorMaps[s]; registered {
				RespondError(ctx, code, errorCodeFor(code), http.StatusText(code))
				return
			}
		}
	}
	handlePanic(ctx, rec, debug.Stack())
}

// handlePanic logs a recovered panic and writes the 500 envelope.
func handlePanic(ctx *context.Context, rec interface{}, stack []byte) {
	EnsureTraceIDs(ctx)
	slog.ErrorContext(ctx.Request.Context(), "panic recovered",
		"panic", fmt.Sprint(rec),
		"method", ctx.Request.Method,
		"route", routePattern(ctx),
		"stack", string(stack),
	)
	RespondError(ctx, http.StatusInternalServerError, "internal_error", "An internal error occurred")
}

// errorCodeFor maps an HTTP status to the machine-readable error code used in
// the JSON envelope.
func errorCodeFor(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "bad_request"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusForbidden:
		return "forbidden"
	case http.StatusNotFound:
		return "not_found"
	case http.StatusMethodNotAllowed:
		return "method_not_allowed"
	case http.StatusRequestEntityTooLarge:
		return "payload_too_large"
	case http.StatusNotImplemented:
		return "not_implemented"
	case http.StatusBadGateway:
		return "bad_gateway"
	case http.StatusServiceUnavailable:
		return "service_unavailable"
	case http.StatusGatewayTimeout:
		return "gateway_timeout"
	case http.StatusInternalServerError:
		return "internal_error"
	default:
		if status >= 500 {
			return "internal_error"
		}
		return "error"
	}
}
