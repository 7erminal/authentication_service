package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "ChangeCustomerPassword",
            Router: `/change-customer-password/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "ChangePassword",
            Router: `/change-password/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "CheckCustomerTokenExpiry",
            Router: `/customer-token/check`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "VerifyCustomerToken",
            Router: `/customer-token/verify`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "ExpireCustomerToken",
            Router: `/expire-customer-token`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "Logout",
            Router: `/log-out`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

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
            Method: "RefreshCustomerAccessToken",
            Router: `/refresh/customer/token`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "RefreshAccessToken",
            Router: `/refresh/user/token`,
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
            Method: "ResetCustomerPassword",
            Router: `/reset-customer-password/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "ResetPasswordLink",
            Router: `/reset-password-link`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "ResetPassword",
            Router: `/reset-password/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "SendActivationCode",
            Router: `/send-activation-code`,
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
            Method: "VerifyToken",
            Router: `/token/verify`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "ValidateCustomerCredentialsToken",
            Router: `/validate-customer-credentials/token`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:AuthenticationController"],
        beego.ControllerComments{
            Method: "VerifyActivationCode",
            Router: `/verify-activation-code`,
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
