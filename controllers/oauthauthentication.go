package controllers

import (
	"authentication_service/api"
	"authentication_service/structs/responsesDTOs"
	"encoding/json"
	"io"

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

// Get ...
// @Title Create
// @Description create OAuthAuthentication
// @Param	body		body 	models.OAuthAuthentication	true		"body for OAuthAuthentication content"
// @Success 201 {object} models.OAuthAuthentication
// @Failure 403 body is empty
// @router /google/authorize [get]
func (c *OAuthAuthenticationController) OAuthGoogle() {
	host, err := beego.AppConfig.String("host")

	logs.Info("Callback Received ", c.Ctx.Input.Query("code"))

	var code string = c.Ctx.Input.Query("code")

	request := api.NewRequest(
		"https://oauth2.googleapis.com",
		"/token",
		api.POST)
	request.Params["code"] = code
	request.Params["client_id"] = "1027199556532-m0a6r4sb74dd8oah3bnoo0igeahgvvis.apps.googleusercontent.com"
	request.Params["client_secret"] = "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"
	request.Params["redirect_uri"] = host + "/v1/oauth/google/authorize"
	request.Params["grant_type"] = "authorization_code"
	client := api.Client{
		Request: request,
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		c.Data["json"] = err.Error()
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		c.Data["json"] = err.Error()
	}

	logs.Info("Raw response received is ", res)
	// data := map[string]interface{}{}
	var data responsesDTOs.GoogleOAuthRespDTO
	json.Unmarshal(read, &data)
	c.Data["json"] = data

	// logs.Info("Response received ", c.Data["json"])
	// logs.Info("Access token ", data["access_token"])
	// logs.Info("Expires in ", data["expires_in"])
	// logs.Info("Scope is ", data["scope"])
	// logs.Info("Token Type is ", data["token_type"])
	logs.Info("Response received ", c.Data["json"])
	logs.Info("Access token ", data.Access_token)
	logs.Info("Expires in ", data.Expires_in)
	logs.Info("Scope is ", data.Scope)
	logs.Info("Token Type is ", data.Token_type)

	Register(data.Access_token)

	c.ServeJSON()
}

// Get ...
// @Title Create
// @Description create OAuthAuthentication
// @Param	body		body 	models.OAuthAuthentication	true		"body for OAuthAuthentication content"
// @Success 201 {object} models.OAuthAuthentication
// @Failure 403 body is empty
// @router /third-party/login [get]
func (c *OAuthAuthenticationController) OAuthThirdPartyLogin() {
	host, _ := beego.AppConfig.String("host")

	logs.Info("Callback Received ", c.Ctx.Input.Query("third-party"))
	// c.Ctx.Input.Param(":id")

	var url_ string = ""

	if c.Ctx.Input.Query("third-party") == "google" {
		url_ = "https://accounts.google.com/o/oauth2/v2/auth?client_id=1027199556532-m0a6r4sb74dd8oah3bnoo0igeahgvvis.apps.googleusercontent.com&redirect_uri=" + host + "/v1/oauth/google/authorize&response_type=code&scope=openid%20profile%20email&access_type=offline&include_granted_scopes=true"
	}

	c.Data["json"] = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: url_, StatusDesc: "Successful"}

	c.ServeJSON()

}

func Register(Token_ string) {

	// Get user details from google and register
	request := api.NewRequest(
		"https://www.googleapis.com",
		"/oauth2/v3/userinfo",
		api.GET)
	request.HeaderField["Authorization"] = Token_
	request.HeaderField["Host"] = "www.googleapis.com"
	client := api.Client{
		Request: request,
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		// c.Data["json"] = err.Error()
	}
	defer res.Body.Close()
	// read, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	logs.Error("client.Error: %v", err)
	// }

	logs.Info("Raw response for user info received is ", res)

	// var data responsesDTOs.GoogleOAuthRespDTO
	// json.Unmarshal(read, &data)
}
