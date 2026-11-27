package controllers

import (
	"authentication_service/internal/observability"
	"authentication_service/internal/trace"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
)

// tracedController preserves existing application response contracts. Only an
// additive request_id is introduced on object-shaped failures; string responses
// are deliberately not converted into a different JSON type.
type tracedController struct{ beego.Controller }

func (c *tracedController) respondError(status int, code, message string) {
	observability.RespondError(c.Ctx.ResponseWriter, c.Ctx.Request, status, code, message)
}

func (c *tracedController) ServeJSON(encoding ...bool) error {
	data, err := json.Marshal(c.Data["json"])
	if err == nil {
		var object map[string]json.RawMessage
		if json.Unmarshal(data, &object) == nil && object != nil {
			var code int
			var valid *bool
			_ = json.Unmarshal(object["StatusCode"], &code)
			_ = json.Unmarshal(object["IsValid"], &valid)
			if code >= 300 || c.Ctx.Output.Status >= 400 || valid != nil && !*valid {
				trace.Logger(c.Ctx.Request.Context()).Warn("application error", "code", code, "http_status", c.Ctx.Output.Status)
				object["request_id"], _ = json.Marshal(trace.FromContext(c.Ctx.Request.Context()).RequestID)
				c.Data["json"] = object
			}
		}
	}
	return c.Controller.ServeJSON(encoding...)
}
