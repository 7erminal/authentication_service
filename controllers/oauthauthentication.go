package controllers

import (
	"authentication_service/api"
	"authentication_service/internal/logger"
	"authentication_service/structs/responsesDTOs"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/url"
	"os"

	beego "github.com/beego/beego/v2/server/web"
)

// Google OAuth credentials are read from the environment; they must not be
// hardcoded in source.
const (
	EnvGoogleClientID     = "GOOGLE_CLIENT_ID"
	EnvGoogleClientSecret = "GOOGLE_CLIENT_SECRET"
)

func googleOAuthCredentials() (clientID, clientSecret string) {
	return os.Getenv(EnvGoogleClientID), os.Getenv(EnvGoogleClientSecret)
}

// GoogleOAuthConfigured reports whether both Google OAuth env vars are set.
func GoogleOAuthConfigured() bool {
	id, secret := googleOAuthCredentials()
	return id != "" && secret != ""
}

// OAuthAuthenticationController operations for OAuthAuthentication
type OAuthAuthenticationController struct {
	BaseController
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
	ctx := c.reqCtx()
	host, err := beego.AppConfig.String("host")

	// Never log the authorization code or anything Google returns (tokens).
	slog.InfoContext(ctx, "google oauth callback received")

	var code string = c.Ctx.Input.Query("code")

	clientID, clientSecret := googleOAuthCredentials()
	if clientID == "" || clientSecret == "" {
		slog.ErrorContext(ctx, "google oauth credentials not configured", "env", EnvGoogleClientID+"/"+EnvGoogleClientSecret)
	}

	request := api.NewRequest(
		"https://oauth2.googleapis.com",
		"/token",
		api.POST)
	request.Params["code"] = code
	request.Params["client_id"] = clientID
	request.Params["client_secret"] = clientSecret
	request.Params["redirect_uri"] = host + "/v1/oauth/google/authorize"
	request.Params["grant_type"] = "authorization_code"
	// Third-party call: ctx for cancellation/logging, but internal tracing
	// headers are NOT forwarded to Google (api.NewRequest, not
	// NewInternalRequest).
	client := api.Client{
		Request: request,
		Ctx:     ctx,
	}
	res, err := client.SendRequest()
	if err != nil {
		slog.ErrorContext(ctx, "google token exchange failed", logger.Err(err))
		c.Data["json"] = err.Error()
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		c.Data["json"] = err.Error()
	}

	slog.InfoContext(ctx, "google token exchange response received", "status", res.StatusCode)
	// data := map[string]interface{}{}
	var data responsesDTOs.GoogleOAuthRespDTO
	json.Unmarshal(read, &data)
	c.Data["json"] = data

	// logs.Info("Response received ", c.Data["json"])
	// logs.Info("Access token ", data["access_token"])
	// logs.Info("Expires in ", data["expires_in"])
	// logs.Info("Scope is ", data["scope"])
	// logs.Info("Token Type is ", data["token_type"])
	slog.InfoContext(ctx, "google token exchange completed", "expires_in", data.Expires_in, "token_type", data.Token_type)

	RegisterCtx(ctx, data.Access_token)

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

	slog.InfoContext(c.reqCtx(), "third party login url requested")
	// c.Ctx.Input.Param(":id")

	var url_ string = ""

	if c.Ctx.Input.Query("third-party") == "google" {
		clientID, _ := googleOAuthCredentials()
		if clientID == "" {
			slog.ErrorContext(c.reqCtx(), "google oauth client id not configured", "env", EnvGoogleClientID)
		}
		url_ = "https://accounts.google.com/o/oauth2/v2/auth?client_id=" + url.QueryEscape(clientID) + "&redirect_uri=" + host + "/v1/oauth/google/authorize&response_type=code&scope=openid%20profile%20email&access_type=offline&include_granted_scopes=true"
	}

	c.Data["json"] = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: url_, StatusDesc: "Successful"}

	c.ServeJSON()

}

func Register(Token_ string) {
	RegisterCtx(context.Background(), Token_)
}

// RegisterCtx is Register with request context propagation.
func RegisterCtx(ctx context.Context, Token_ string) {

	// Get user details from google and register
	request := api.NewRequest(
		"https://www.googleapis.com",
		"/oauth2/v3/userinfo",
		api.GET)
	request.HeaderField["Authorization"] = Token_
	request.HeaderField["Host"] = "www.googleapis.com"
	client := api.Client{
		Request: request,
		Ctx:     ctx,
	}
	res, err := client.SendRequest()
	if err != nil {
		slog.ErrorContext(ctx, "google userinfo request failed", logger.Err(err))
		// c.Data["json"] = err.Error()
	}
	defer res.Body.Close()
	// read, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	logs.Error("client.Error: %v", err)
	// }

	slog.InfoContext(ctx, "google userinfo response received", "status", res.StatusCode)

	// var data responsesDTOs.GoogleOAuthRespDTO
	// json.Unmarshal(read, &data)
}
