package controllers

import (
	"context"

	"authentication_service/filters"
	"authentication_service/structs/responsesDTOs"

	beego "github.com/beego/beego/v2/server/web"
)

// BaseController adds tracing-aware helpers to beego.Controller. Controllers
// embed it instead of beego.Controller; behaviour is otherwise unchanged.
type BaseController struct {
	beego.Controller
}

// reqCtx returns the request context carrying the tracing IDs. Pass it to
// slog.*Context, model ...Ctx functions and outbound calls.
func (c *BaseController) reqCtx() context.Context {
	return c.Ctx.Request.Context()
}

// respondError writes {"error":{"code","message","request_id"}} with the
// given HTTP status. Use for NEW error paths only: existing endpoints keep
// their legacy {StatusCode, StatusDesc, ...} bodies because the gateways
// parse them.
func (c *BaseController) respondError(status int, code, message string) {
	filters.RespondError(c.Ctx, status, code, message)
}

// ServeJSON shadows beego.Controller.ServeJSON. If the payload is a legacy
// response DTO describing an error (StatusCode != 200) it adds the additive
// top-level "request_id" field before serialising. Success bodies and
// non-DTO payloads are untouched.
func (c *BaseController) ServeJSON(encoding ...bool) {
	if dto, ok := c.Data["json"].(responsesDTOs.RequestIDCarrier); ok {
		c.Data["json"] = dto.WithRequestIDOnError(filters.RequestIDFrom(c.Ctx))
	}
	c.Controller.ServeJSON(encoding...)
}
