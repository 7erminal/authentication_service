package filters

import (
	"authentication_service/internal/trace"

	"github.com/beego/beego/v2/server/web/context"
)

// ErrorBody is the JSON error envelope:
//
//	{"error": {"code": "...", "message": "...", "request_id": "..."}}
type ErrorBody struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail is the inner object of ErrorBody.
type ErrorDetail struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

// RespondError writes the JSON error envelope with the given HTTP status. The
// request_id is always included. If the response has already started it does
// nothing (the status line can no longer change).
func RespondError(ctx *context.Context, status int, code, message string) {
	if ctx.ResponseWriter.Started {
		return
	}
	ids := EnsureTraceIDs(ctx)
	ctx.Output.SetStatus(status)
	_ = ctx.Output.JSON(ErrorBody{Error: ErrorDetail{
		Code:      code,
		Message:   message,
		RequestID: ids.RequestID,
	}}, false, false)
}

// RequestIDFrom returns the request ID of the current request ("" if none).
func RequestIDFrom(ctx *context.Context) string {
	return trace.RequestID(ctx.Request.Context())
}
