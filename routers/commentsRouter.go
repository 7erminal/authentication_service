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
            Method: "VerifyOTP",
            Router: `/verify-otp`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:UsersController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:UsersController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:UsersController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:UsersController"],
        beego.ControllerComments{
            Method: "SignUp2",
            Router: `/2/sign-up`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:UsersController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:UsersController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:UsersController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:UsersController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:UsersController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:UsersController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:UsersController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:UsersController"],
        beego.ControllerComments{
            Method: "VerifyUsername",
            Router: `/:username`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:UsersController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:UsersController"],
        beego.ControllerComments{
            Method: "ResendOtp",
            Router: `/resend-otp/:username`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["authentication_service/controllers:UsersController"] = append(beego.GlobalControllerRouter["authentication_service/controllers:UsersController"],
        beego.ControllerComments{
            Method: "SignUp",
            Router: `/sign-up`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
