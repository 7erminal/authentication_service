package controllers

import (
	"net/http"

	"authentication_service/filters"

	beego "github.com/beego/beego/v2/server/web"
)

// ErrorController replaces Beego's default HTML error pages (unknown route,
// method not allowed, ...) with the JSON error envelope including request_id.
// Registered in main.go via beego.ErrorController.
type ErrorController struct {
	beego.Controller
}

// Prepare disables template rendering: these handlers write JSON directly,
// and Beego would otherwise look for views/errorcontroller/errorXXX.tpl when
// autorender is on (e.g. if conf/app.conf is missing).
func (c *ErrorController) Prepare() { c.EnableRender = false }

func (c *ErrorController) write(status int, code string) {
	filters.RespondError(c.Ctx, status, code, http.StatusText(status))
}

// Error401 ...
func (c *ErrorController) Error401() { c.write(http.StatusUnauthorized, "unauthorized") }

// Error403 ...
func (c *ErrorController) Error403() { c.write(http.StatusForbidden, "forbidden") }

// Error404 ...
func (c *ErrorController) Error404() { c.write(http.StatusNotFound, "not_found") }

// Error405 ...
func (c *ErrorController) Error405() { c.write(http.StatusMethodNotAllowed, "method_not_allowed") }

// Error413 ...
func (c *ErrorController) Error413() {
	c.write(http.StatusRequestEntityTooLarge, "payload_too_large")
}

// Error500 ...
func (c *ErrorController) Error500() {
	c.write(http.StatusInternalServerError, "internal_error")
}

// Error501 ...
func (c *ErrorController) Error501() { c.write(http.StatusNotImplemented, "not_implemented") }

// Error502 ...
func (c *ErrorController) Error502() { c.write(http.StatusBadGateway, "bad_gateway") }

// Error503 ...
func (c *ErrorController) Error503() {
	c.write(http.StatusServiceUnavailable, "service_unavailable")
}

// Error504 ...
func (c *ErrorController) Error504() { c.write(http.StatusGatewayTimeout, "gateway_timeout") }
