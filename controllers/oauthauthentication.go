package controllers

import (
	"authentication_service/api"
	"authentication_service/structs/responsesDTOs"
	"context"
	"encoding/json"
	"io"
	"net/url"
	"os"

	beego "github.com/beego/beego/v2/server/web"
)

// OAuthAuthenticationController operations for OAuthAuthentication
type OAuthAuthenticationController struct {
	tracedController
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
	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		c.respondError(500, "internal_error", "Internal server error")
		return
	}

	host, err := beego.AppConfig.String("host")

	var code string = c.Ctx.Input.Query("code")

	request := api.NewRequest(
		"https://oauth2.googleapis.com",
		"/token",
		api.POST)
	request.Params["code"] = code
	request.Params["client_id"] = clientID
	request.Params["client_secret"] = clientSecret
	request.Params["redirect_uri"] = host + "/v1/oauth/google/authorize"
	request.Params["grant_type"] = "authorization_code"
	client := api.Client{
		Request: request,
	}
	res, err := client.SendRequestContext(c.Ctx.Request.Context())
	if err != nil {

		c.respondError(500, "internal_error", "Internal server error")
		return
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		c.Data["json"] = err.Error()
	}

	// data := map[string]interface{}{}
	var data responsesDTOs.GoogleOAuthRespDTO
	json.Unmarshal(read, &data)
	c.Data["json"] = data

	if err := RegisterContext(c.Ctx.Request.Context(), data.Access_token); err != nil {
		c.respondError(500, "internal_error", "Internal server error")
		return
	}

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

	// c.Ctx.Input.Param(":id")

	var url_ string = ""

	if c.Ctx.Input.Query("third-party") == "google" {
		clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
		if clientID == "" {
			c.respondError(500, "internal_error", "Internal server error")
			return
		}

		url_ = "https://accounts.google.com/o/oauth2/v2/auth?client_id=" + url.QueryEscape(clientID) + "&redirect_uri=" + host + "/v1/oauth/google/authorize&response_type=code&scope=openid%20profile%20email&access_type=offline&include_granted_scopes=true"
	}

	c.Data["json"] = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: url_, StatusDesc: "Successful"}

	c.ServeJSON()

}

// Register preserves the legacy non-request entry point and its failure signal.
func Register(Token_ string) {
	if err := RegisterContext(context.Background(), Token_); err != nil {
		panic(err)
	}
}

func RegisterContext(ctx context.Context, Token_ string) error {

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
	res, err := client.SendRequestContext(ctx)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// read, err := io.ReadAll(res.Body)
	// if err != nil {
	// }

	// var data responsesDTOs.GoogleOAuthRespDTO
	// json.Unmarshal(read, &data)
	return nil
}
