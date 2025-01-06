package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/login`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "LoginToken",
            Router: `/login/token`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "ResendOTP",
            Router: `/resend-otp`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "CheckTokenExpiry",
            Router: `/token/check`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "GenerateInviteToken",
            Router: `/token/invite`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "VerifyInviteToken",
            Router: `/token/invite/verify`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "VerifyOTP",
            Router: `/verify-otp`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:OAuthAuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:OAuthAuthenticationController"],
        beego.ControllerComments{
            Method: "OAuthGoogle",
            Router: `/google/authorize`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:OAuthAuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:OAuthAuthenticationController"],
        beego.ControllerComments{
            Method: "OAuthThirdPartyLogin",
            Router: `/third-party/login`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
