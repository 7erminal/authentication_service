package controllers

import (
	"authentication_service/models"
	"encoding/json"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
)

// AuthenticationController operations for Authentication
type AuthenticationController struct {
	beego.Controller
}

// URLMapping ...
func (c *AuthenticationController) URLMapping() {
	c.Mapping("Login", c.Login)
	c.Mapping("VerifyOTP", c.VerifyOTP)
}

// Post ...
// @Title Create
// @Description Login
// @Param	body		body 	models.AuthenticationDTO	true		"body for Authentication content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /login [post]
func (c *AuthenticationController) Login() {
	var v models.AuthenticationDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	logs.Info("Received ", v.Password, v.Username)

	if a, err := models.GetUsersByUsername(v.Username); err == nil {
		// Compare the stored hashed password, with the hashed version of the password that was received
		if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(v.Password)); err != nil {
			// If the two passwords don't match, return a 401 status
			c.Data["json"] = err.Error()

			logs.Error(err.Error())

			var resp = models.UserResponseDTO{StatusCode: 605, User: nil, StatusDesc: "Incorrect password"}
			c.Data["json"] = resp

		} else {
			c.Ctx.Output.SetStatus(201)

			var resp = models.UserResponseDTO{StatusCode: 200, User: a, StatusDesc: "User has been authenticated"}
			c.Data["json"] = resp
		}
	} else {
		logs.Error(err.Error())
		var resp = models.UserResponseDTO{StatusCode: 605, User: nil, StatusDesc: "Unidentified user"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// Verify OTP ...
// @Title Verify OTP
// @Description Login
// @Param	body		body 	models.AuthenticationDTO	true		"body for Authentication content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /verify-otp [post]
func (c *AuthenticationController) VerifyOTP() {
	var v models.AuthenticationDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	logs.Info("Received OTP verification request: ", v.Password, v.Username)

	if a, err := models.GetUsersByUsername(v.Username); err == nil {
		// Verify OTP
		if uo, err := models.VerifyUserOTP(a.UserId); err == nil {
			logs.Info("Returned OTP model is ", uo)

			// Compare the stored hashed password, with the hashed version of the password that was received
			if err := bcrypt.CompareHashAndPassword([]byte(uo.OneTimePin), []byte(v.Password)); err != nil {
				// If the two passwords don't match, return a 401 status
				c.Data["json"] = err.Error()

				logs.Error(err.Error())

				var resp = models.UserResponseDTO{StatusCode: 605, User: nil, StatusDesc: "Incorrect password"}
				c.Data["json"] = resp

			} else {
				c.Ctx.Output.SetStatus(201)

				var resp = models.UserResponseDTO{StatusCode: 200, User: a, StatusDesc: "User OTP matches"}
				c.Data["json"] = resp
			}
		} else {
			logs.Error(err.Error())
			var resp = models.UserResponseDTO{StatusCode: 605, User: nil, StatusDesc: "No OTP found for this user"}
			c.Data["json"] = resp
		}
	} else {
		logs.Error(err.Error())
		var resp = models.UserResponseDTO{StatusCode: 605, User: nil, StatusDesc: "Unidentified user"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}
