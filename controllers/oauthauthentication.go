package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// OAuthAuthenticationController operations for OAuthAuthentication
type OAuthAuthenticationController struct {
	beego.Controller
}

// URLMapping ...
func (c *OAuthAuthenticationController) URLMapping() {
	c.Mapping("GoogleAuth", c.OAuthGoogle)
}

// Post ...
// @Title Create
// @Description create OAuthAuthentication
// @Param	body		body 	models.OAuthAuthentication	true		"body for OAuthAuthentication content"
// @Success 201 {object} models.OAuthAuthentication
// @Failure 403 body is empty
// @router /oauthgoogle/authorize [post]
func (c *OAuthAuthenticationController) OAuthGoogle() {
	logs.Info("Callback Received ")
}
