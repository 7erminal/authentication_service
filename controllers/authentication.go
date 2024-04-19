package controllers

import (
	"authentication_service/controllers/functions"
	"authentication_service/models"
	"authentication_service/structs/requestsDTOs"
	"authentication_service/structs/responsesDTOs"
	"encoding/json"
	"time"

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
	c.Mapping("ResendOTP", c.ResendOTP)
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

			var resp = responsesDTOs.UserResponseDTO{StatusCode: 605, User: nil, StatusDesc: "Incorrect password"}
			c.Data["json"] = resp

		} else {
			c.Ctx.Output.SetStatus(201)

			var resp = responsesDTOs.UserResponseDTO{StatusCode: 200, User: a, StatusDesc: "User has been authenticated"}
			c.Data["json"] = resp
		}
	} else {
		logs.Error(err.Error())
		var resp = responsesDTOs.UserResponseDTO{StatusCode: 605, User: nil, StatusDesc: "Unidentified user"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// Verify OTP ...
// @Title Verify OTP
// @Description Verify OTP using username
// @Param	body		body 	requestsDTOs.VerifyOtpDTO	true		"body for Verify OTP content"
// @Success 201 {object} requestsDTOs.UserResponseDTO
// @Failure 403 body is empty
// @router /verify-otp [post]
func (c *AuthenticationController) VerifyOTP() {
	// username := c.Ctx.Input.Param(":username")
	var q requestsDTOs.VerifyOtpDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	v, err := models.GetUsersByUsername(q.Username)
	logs.Debug("Checking for username ", q.Username)

	if err != nil {
		var resp = responsesDTOs.UserResponseDTO{StatusCode: 604, User: nil, StatusDesc: "User cannot be found"}
		// c.Data["json"] = err.Error()
		c.Data["json"] = resp
	} else {
		// Get OTP
		logs.Debug("Got user. Now checking for user in OTP table ", v.UserId, v.Email, v.FullName)
		otp, err := models.VerifyUserOTP(v.UserId)

		logs.Debug("User in OTP table ")

		if err == nil {
			if q.Password == otp.Code {
				logs.Debug("OTP Passed")
				logs.Debug("About to compare OTP expiry date...", otp.ExpiryDate, " with date now ", time.Now())
				if otp.ExpiryDate.After(time.Now()) {
					if otp.Status == 1 {
						logs.Debug("OTP has been used already.")
						var resp = responsesDTOs.UserResponseDTO{StatusCode: 407, User: v, StatusDesc: "OTP has already been used."}
						c.Data["json"] = resp
					} else {
						otp.Status = 1
						if err := models.UpdateUserOtpById(otp); err == nil {
							var resp = responsesDTOs.UserResponseDTO{StatusCode: 200, User: v, StatusDesc: "OTP Verified successfully"}
							c.Data["json"] = resp
						} else {
							logs.Error("Error is ", err.Error())
							var resp = responsesDTOs.UserResponseDTO{StatusCode: 403, User: v, StatusDesc: "Error occurred inserting record."}
							c.Data["json"] = resp
						}
					}
				} else {
					logs.Debug("OTP has expired. Time to enter OTP of 5 mins exeeded.")
					var resp = responsesDTOs.UserResponseDTO{StatusCode: 403, User: v, StatusDesc: "OTP Expired"}
					c.Data["json"] = resp
				}
			} else {
				logs.Debug("OTPs do not match ")
				var resp = responsesDTOs.UserResponseDTO{StatusCode: 402, User: v, StatusDesc: "OTP Verification failed"}
				c.Data["json"] = resp
			}
		} else {
			logs.Debug("Error: ", err.Error(), " User not in OTP Table ")
			var resp = responsesDTOs.UserResponseDTO{StatusCode: 403, User: v, StatusDesc: "OTP Expired"}
			c.Data["json"] = resp
		}
		// Generate random number
		// randNum := functions.EncodeToString(6)
		// logs.Debug("Random number generated is ", randNum)

		// expiryDate := time.Now().Local().Add(time.Hour*time.Duration(0) + time.Minute*time.Duration(5) + time.Second*time.Duration(0))

		// otpModel := models.User_otps{Code: randNum, User: v.UserId, Status: 2, DateCreated: time.Now(), DateModified: time.Now(), DateGenerated: time.Now(), ExpiryDate: expiryDate, Active: 1}

		// if _, err := models.AddUser_otps(&otpModel); err == nil {
		// 	functions.SendEmail(v.Email, randNum)

		// 	var resp = models.UserResponseDTO{StatusCode: 200, User: v, StatusDesc: "Email sent successfully"}
		// 	c.Data["json"] = resp
		// } else {
		// 	var resp = models.UserResponseDTO{StatusCode: 703, User: v, StatusDesc: "Error sending email"}
		// 	c.Data["json"] = resp
		// }
	}
	c.ServeJSON()
}

// Resend OTP ...
// @Title Resend OTP
// @Description Resend OTP using username
// @Param	body		body 	requestsDTOs.UsernameDTO	true		"body for SignUp content"
// @Success 201 {object} responsesDTOs.UserResponseDTO
// @Failure 403 body is empty
// @router /resend-otp [post]
func (c *AuthenticationController) ResendOTP() {
	// username := c.Ctx.Input.Param(":username")
	var q requestsDTOs.UsernameDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	v, err := models.GetUsersByUsername(q.Username)

	if err != nil {
		var resp = responsesDTOs.UserResponseDTO{StatusCode: 604, User: nil, StatusDesc: "User cannot be found"}
		// c.Data["json"] = err.Error()
		c.Data["json"] = resp
	} else {
		// Generate random number
		randNum := functions.EncodeToString(6)
		logs.Debug("Random number generated is ", randNum)

		expiryDate := time.Now().Local().Add(time.Hour*time.Duration(0) + time.Minute*time.Duration(5) + time.Second*time.Duration(0))

		otpModel := models.UserOtps{Code: randNum, UserId: v.UserId, Status: 2, DateCreated: time.Now(), DateModified: time.Now(), DateGenerated: time.Now(), ExpiryDate: expiryDate, Active: 1}

		if _, err := models.AddUserOtp(&otpModel); err == nil {
			functions.SendEmail(v.Email, randNum)

			var resp = responsesDTOs.UserResponseDTO{StatusCode: 200, User: v, StatusDesc: "Email sent successfully"}
			c.Data["json"] = resp
		} else {
			logs.Error("Error inserting OTP...", err.Error())
			var resp = responsesDTOs.UserResponseDTO{StatusCode: 703, User: v, StatusDesc: "Error sending email"}
			c.Data["json"] = resp
		}
	}
	c.ServeJSON()
}
