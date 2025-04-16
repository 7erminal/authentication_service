package controllers

import (
	"authentication_service/controllers/functions"
	"authentication_service/models"
	"authentication_service/structs/requestsDTOs"
	"authentication_service/structs/responsesDTOs"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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
	c.Mapping("LoginToken", c.LoginToken)
	c.Mapping("CheckTokenExpiry", c.CheckTokenExpiry)
	c.Mapping("GenerateInviteToken", c.GenerateInviteToken)
	c.Mapping("VerifyInviteToken", c.VerifyInviteToken)
	c.Mapping("ChangePassword", c.ChangePassword)
	c.Mapping("ResetPassword", c.ResetPassword)
	c.Mapping("SendActivationCode", c.SendActivationCode)
	c.Mapping("VerifyActivationCode", c.VerifyActivationCode)
	c.Mapping("ResetPasswordLink", c.ResetPasswordLink)
	c.Mapping("Logout", c.Logout)
	c.Mapping("VerifyToken", c.VerifyToken)
}

// Post ...
// @Title Login
// @Description Login User
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
			// cust, err := models.GetCustomersByUser(a)

			// if err != nil {
			// 	c.Data["json"] = err.Error()

			// 	var resp = responsesDTOs.UserResponseDTO{StatusCode: 601, User: nil, StatusDesc: "Error verifying user"}
			// 	c.Data["json"] = resp
			// } else {
			// 	logs.Info("Getting the customer ", cust.Branch.Country.DefaultCurrency.CurrencyId)

			// 	userResp := responsesDTOs.UserResp{
			// 		UserId:        a.UserId,
			// 		ImagePath:     a.ImagePath,
			// 		UserType:      a.UserType,
			// 		FullName:      a.FullName,
			// 		Username:      a.Username,
			// 		Password:      a.Password,
			// 		Email:         a.Email,
			// 		PhoneNumber:   a.PhoneNumber,
			// 		Gender:        a.Gender,
			// 		Dob:           a.Dob,
			// 		Address:       a.Address,
			// 		IdType:        a.IdType,
			// 		IdNumber:      a.IdNumber,
			// 		MaritalStatus: a.MaritalStatus,
			// 		Active:        a.Active,
			// 		Role:          a.Role,
			// 		IsVerified:    a.IsVerified,
			// 		DateCreated:   a.DateCreated,
			// 		DateModified:  a.DateModified,
			// 		CreatedBy:     a.CreatedBy,
			// 		ModifiedBy:    a.ModifiedBy,
			// 		Branch:        cust.Branch,
			// 	}
			// 	c.Ctx.Output.SetStatus(200)

			// 	var resp = responsesDTOs.UserResponseDTO{StatusCode: 200, User: &userResp, StatusDesc: "User has been authenticated"}
			// 	c.Data["json"] = resp
			// }

			c.Ctx.Output.SetStatus(200)

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

// Post ...
// @Title Login User
// @Description Login
// @Param	body		body 	models.AuthenticationDTO	true		"body for Authentication content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /login/token [post]
func (c *AuthenticationController) LoginToken() {
	var v models.AuthenticationDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	logs.Info("Received ", v.Password, v.Username)

	if a, err := models.GetUsersByUsername(v.Username); err == nil {
		// Compare the stored hashed password, with the hashed version of the password that was received
		logs.Info("User role is ", a.Role.Role)
		if a.Active == 1 {
			if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(v.Password)); err != nil {
				// If the two passwords don't match, return a 401 status
				c.Data["json"] = err.Error()

				logs.Error(err.Error())

				var resp = responsesDTOs.UserResponseDTO{StatusCode: 605, User: nil, StatusDesc: "Incorrect password"}
				c.Data["json"] = resp

			} else {
				c.Ctx.Output.SetStatus(200)

				token, expiryTime, err := functions.CreateToken(v.Username)

				logs.Info("Token created is ", token)

				if err != nil {
					logs.Error("Error updating token. ", err.Error())
					var resp = responsesDTOs.StringResponseDTO{StatusCode: 301, Value: "", StatusDesc: "Error generating token"}
					c.Data["json"] = resp
				} else {
					updateToken := models.AccessTokens{User: a, Revoked: true}
					if err := models.UpdateAccessTokensByUserId(&updateToken); err != nil {
						t := time.Unix(expiryTime, 0)
						tokenObj := models.AccessTokens{User: a, Token: token, ExpiresAt: t, DateCreated: time.Now()}
						if _, err := models.AddAccessTokens(&tokenObj); err == nil {
							var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: token, StatusDesc: "User has been authenticated"}
							c.Data["json"] = resp
						} else {
							logs.Error("Error adding token. ", err.Error())
							var resp = responsesDTOs.StringResponseDTO{StatusCode: 301, Value: "", StatusDesc: "Error generating token"}
							c.Data["json"] = resp
						}
					} else {
						t := time.Unix(expiryTime, 0)
						tokenObj := models.AccessTokens{User: a, Token: token, ExpiresAt: t, DateCreated: time.Now()}
						if _, err := models.AddAccessTokens(&tokenObj); err == nil {
							var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: token, StatusDesc: "User has been authenticated"}
							c.Data["json"] = resp
						} else {
							logs.Error("Error adding token. ", err.Error())
							var resp = responsesDTOs.StringResponseDTO{StatusCode: 301, Value: "", StatusDesc: "Error generating token"}
							c.Data["json"] = resp
						}
						// logs.Error("Error updating token. ")
						// var resp = responsesDTOs.StringResponseDTO{StatusCode: 301, Value: "", StatusDesc: "Error generating token"}
						// c.Data["json"] = resp
					}
				}
			}
		} else {
			logs.Error("User is not active ", a.Active)
			var resp = responsesDTOs.UserResponseDTO{StatusCode: 607, User: nil, StatusDesc: "Inactive user"}
			c.Data["json"] = resp
		}
	} else {
		logs.Error(err.Error())
		var resp = responsesDTOs.UserResponseDTO{StatusCode: 605, User: nil, StatusDesc: "Unidentified user"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// ChangePassword ...
// @Title Change Password
// @Description Change user password
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	requestsDTOs.ChangePassword	true		"body for Change password content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /change-password/:id [put]
func (c *AuthenticationController) ChangePassword() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)

	var v requestsDTOs.ChangePassword
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	logs.Info("Received ", v.OldPassword, v.NewPassword)

	if a, err := models.GetUsersById(id); err == nil {
		// Compare the stored hashed password, with the hashed version of the password that was received
		if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(v.OldPassword)); err != nil {
			// If the two passwords don't match, return a 401 status
			c.Data["json"] = err.Error()

			logs.Error(err.Error())

			var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Old password does not match"}
			c.Data["json"] = resp

		} else {
			hashedPassword, errr := bcrypt.GenerateFromPassword([]byte(v.NewPassword), 8)

			if errr == nil {
				logs.Debug(hashedPassword)

				a.Password = string(hashedPassword)

				logs.Debug("Sending", v.NewPassword)

				// models.Agents{AgentName: v.AgentName, BranchId: v.BranchId, IdType: v.IdType, IdNumber: v.IdNumber, IsVerified: false, Active: 1, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: c_by, ModifiedBy: c_by}
			} else {
				logs.Error("Error hashing password ", errr.Error())
			}

			if err := models.UpdateUsersById(a); err == nil {
				c.Ctx.Output.SetStatus(200)

				var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "Successfully changed password", StatusDesc: "User password has been changed successfully"}
				c.Data["json"] = resp
			} else {
				var resp = responsesDTOs.StringResponseDTO{StatusCode: 608, Value: "", StatusDesc: "User password change failed. " + err.Error()}
				c.Data["json"] = resp
			}
		}
	} else {
		logs.Error(err.Error())
		var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Unidentified user"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// Reset Password ...
// @Title Reset Password
// @Description Reset user password
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	requestsDTOs.ResetPassword	true		"body for Change password content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /reset-password/:id [put]
func (c *AuthenticationController) ResetPassword() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)

	var v requestsDTOs.ResetPassword
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	logs.Info("Received ", v.NewPassword)

	logs.Info("About to decrypt token")

	if a, err := models.GetUsersById(id); err == nil {
		// Compare the stored hashed password, with the hashed version of the password that was received

		hashedPassword, errr := bcrypt.GenerateFromPassword([]byte(v.NewPassword), 8)

		if errr == nil {
			logs.Debug(hashedPassword)

			a.Password = string(hashedPassword)

			logs.Debug("Sending", v.NewPassword)

			// models.Agents{AgentName: v.AgentName, BranchId: v.BranchId, IdType: v.IdType, IdNumber: v.IdNumber, IsVerified: false, Active: 1, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: c_by, ModifiedBy: c_by}
		} else {
			logs.Error("Error hashing password ", errr.Error())
		}

		if err := models.UpdateUsersById(a); err == nil {
			c.Ctx.Output.SetStatus(200)

			var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "Successfully changed password", StatusDesc: "User password has been changed successfully"}
			c.Data["json"] = resp
		} else {
			var resp = responsesDTOs.StringResponseDTO{StatusCode: 608, Value: "", StatusDesc: "User password change failed. " + err.Error()}
			c.Data["json"] = resp
		}

	} else {
		logs.Error(err.Error())
		var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Unidentified user"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// Reset Password Link...
// @Title Reset Password Link
// @Description Reset user password link
// @Param	body		body 	requestsDTOs.ResetPasswordLink	true		"body for Change password content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /reset-password-link [post]
func (c *AuthenticationController) ResetPasswordLink() {
	var v requestsDTOs.ResetPasswordLink
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	logs.Info("Received ", v.Email)

	if a, err := models.GetUsersByUsername(v.Email); err == nil {
		// Compare the stored hashed password, with the hashed version of the password that was received

		// hashedPassword, errr := bcrypt.GenerateFromPassword([]byte(v.NewPassword), 8)

		// logs.Debug(hashedPassword)
		fmt.Printf("Value of v: %+v\n", a)

		rawString := v.Email + "___" + a.Role.Role

		// ikey, _ := functions.GenerateKey()

		if token, nonce, err := functions.GetAESEncrypted(rawString); err == nil {

			logs.Info("Token generated is ", token)
			logs.Info("Nonce generated is ", nonce)
			// logs.Info("Key generated is ", string(ikey[:]))

			var userToken models.UserTokens = models.UserTokens{Token: token, Nonce: nonce, ExpiryDate: time.Now().Add(4 * time.Hour), Active: 1, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: 1, ModifiedBy: 1}

			_, err := models.AddUserTokens(&userToken)

			if err != nil {
				return
			}

			logs.Debug("Message is ", v.Message)
			logs.Debug("Subject is ", v.Subject)
			logs.Debug("Links are ", v.Links)
			logs.Debug("Sender is ", a.FullName)

			message_ := strings.Replace(v.Message, "[SENDER_NAME_ID]", a.FullName, -1)
			logs.Info("Message with name is ", message_)

			for i, link := range v.Links {
				iStr := strconv.Itoa(i)
				placeholder := "[LINK_" + iStr + "_ID]"
				formattedLink := *link + token
				message_ = strings.Replace(message_, placeholder, formattedLink, -1)
				logs.Info("Message with link is ", message_)
			}

			logs.Debug("Sending", message_)

			go functions.SendEmailNew(a.Email, v.Subject, message_)

		} else {
			logs.Error("Error validating token...", err.Error())
			var resp = responsesDTOs.StringResponseDTO{StatusCode: 703, Value: "", StatusDesc: "Error validating token"}
			c.Data["json"] = resp
		}

		// models.Agents{AgentName: v.AgentName, BranchId: v.BranchId, IdType: v.IdType, IdNumber: v.IdNumber, IsVerified: false, Active: 1, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: c_by, ModifiedBy: c_by}

		var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "Successfully sent reset link", StatusDesc: "User password link sent"}
		c.Data["json"] = resp
		// if err := models.UpdateUsersById(a); err == nil {
		// 	c.Ctx.Output.SetStatus(200)

		// 	var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "Successfully changed password", StatusDesc: "User password has been changed successfully"}
		// 	c.Data["json"] = resp
		// } else {
		// 	var resp = responsesDTOs.StringResponseDTO{StatusCode: 608, Value: "", StatusDesc: "User password change failed. " + err.Error()}
		// 	c.Data["json"] = resp
		// }

	} else {
		logs.Error(err.Error())
		var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Unidentified user"}
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
						var resp = responsesDTOs.UserResponseDTO{StatusCode: 407, User: nil, StatusDesc: "OTP has already been used."}
						c.Data["json"] = resp
					} else {
						otp.Status = 1

						if err := models.UpdateUserOtpById(otp); err == nil {
							var resp = responsesDTOs.UserResponseDTO{StatusCode: 200, User: nil, StatusDesc: "OTP Verified successfully"}
							c.Data["json"] = resp
						} else {
							logs.Error("Error is ", err.Error())
							var resp = responsesDTOs.UserResponseDTO{StatusCode: 403, User: nil, StatusDesc: "Error occurred inserting record."}
							c.Data["json"] = resp
						}
					}
				} else {
					logs.Debug("OTP has expired. Time to enter OTP of 5 mins exeeded.")
					var resp = responsesDTOs.UserResponseDTO{StatusCode: 403, User: nil, StatusDesc: "OTP Expired"}
					c.Data["json"] = resp
				}
			} else {
				logs.Debug("OTPs do not match ")
				var resp = responsesDTOs.UserResponseDTO{StatusCode: 402, User: nil, StatusDesc: "OTP Verification failed"}
				c.Data["json"] = resp
			}
		} else {
			logs.Debug("Error: ", err.Error(), " User not in OTP Table ")
			var resp = responsesDTOs.UserResponseDTO{StatusCode: 403, User: nil, StatusDesc: "OTP Expired"}
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

			var resp = responsesDTOs.UserResponseDTO{StatusCode: 200, User: nil, StatusDesc: "Email sent successfully"}
			c.Data["json"] = resp
		} else {
			logs.Error("Error inserting OTP...", err.Error())
			var resp = responsesDTOs.UserResponseDTO{StatusCode: 703, User: nil, StatusDesc: "Error sending email"}
			c.Data["json"] = resp
		}
	}
	c.ServeJSON()
}

// Send Activation Code ...
// @Title Send Activation Code
// @Description Send Activation Code
// @Param	body		body 	requestsDTOs.SendActivationCode	true		"body for SignUp content"
// @Success 201 {object} responsesDTOs.StringResponseDTO
// @Failure 403 body is empty
// @router /send-activation-code [post]
func (c *AuthenticationController) SendActivationCode() {
	// username := c.Ctx.Input.Param(":username")
	var q requestsDTOs.SendActivationCode
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	// Generate random number
	randNum := functions.EncodeToString(6)
	// set OTP to 1111 for tests
	randNum = "1111"
	logs.Debug("Random number generated is ", randNum)
	logs.Debug("Mobile number in request is ", q.MobileNumber)
	proceed := false

	if ac, err := models.GetActivationCodesByNumber(q.MobileNumber); err == nil {
		// fmt.Printf("Value of v: %+v\n", ac)
		for _, suc := range ac {
			fmt.Printf("Value of v: %+v\n", suc)
			singleAc, err := suc.(models.ActivationCodes)
			logs.Debug("Activation...")
			fmt.Printf("Type of v: %T\n", singleAc)
			fmt.Printf("Value of v: %+v\n", singleAc)
			logs.Debug(singleAc)
			if !err {
				singleAc.ExpiryDate = time.Now()
				models.UpdateActivationCodesById(&singleAc)
			}
		}
		proceed = true
	} else {
		proceed = true
	}

	if proceed {
		expiryDate := time.Now().Local().Add(time.Hour*time.Duration(1) + time.Minute*time.Duration(0) + time.Second*time.Duration(0))

		otpModel := models.ActivationCodes{Code: randNum, Number: q.MobileNumber, DateCreated: time.Now(), DateModified: time.Now(), ExpiryDate: expiryDate, Active: 1}

		if _, err := models.AddActivationCodes(&otpModel); err == nil {
			// Function to send Code via sms
			// functions.SendEmail(v.Email, randNum)

			var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "SUCCESS", StatusDesc: "Email sent successfully"}
			c.Data["json"] = resp
		} else {
			logs.Error("Error inserting Activation code...", err.Error())
			var resp = responsesDTOs.StringResponseDTO{StatusCode: 703, Value: "FAILED", StatusDesc: "Error sending email"}
			c.Data["json"] = resp
		}
	} else {
		logs.Error("Unable to perform send code due to failure...")
		var resp = responsesDTOs.StringResponseDTO{StatusCode: 703, Value: "FAILED", StatusDesc: "An error occurred when sending activation code. Please try again."}
		c.Data["json"] = resp
	}

	c.ServeJSON()
}

// Verify Activation Code ...
// @Title Verify Activation Code
// @Description Verify Activation code
// @Param	body		body 	requestsDTOs.VerifyActivationCodeDTO	true		"body for Verify OTP content"
// @Success 201 {object} responsesDTOs.StringResponseDTO
// @Failure 403 body is empty
// @router /verify-activation-code [post]
func (c *AuthenticationController) VerifyActivationCode() {
	// username := c.Ctx.Input.Param(":username")
	var q requestsDTOs.VerifyActivationCodeDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	logs.Debug("Got request with mobile number ", q.MobileNumber, " and pin ", q.Password)

	if v, err := models.GetActivationCodeByNumber(q.MobileNumber); err != nil {
		logs.Error("Code cannot be found. It either does not exist or has expired :: ", err.Error())
		var resp = responsesDTOs.StringResponseDTO{StatusCode: 604, Value: "FAILED", StatusDesc: "Code cannot be found. It either does not exist or has expired"}
		// c.Data["json"] = err.Error()
		c.Data["json"] = resp
	} else {
		// Get OTP

		if q.Password == v.Code {
			logs.Debug("OTP Passed")
			logs.Debug("About to compare OTP expiry date...", v.ExpiryDate, " with date now ", time.Now())
			if v.ExpiryDate.After(time.Now()) {

				v.ExpiryDate = time.Now()
				logs.Debug("Expiry date is now ", v.ExpiryDate, " and ID is ", v.ActivationCodeId)
				if err := models.UpdateActivationCodesById(v); err == nil {
					var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "SUCCESS", StatusDesc: "OTP Verified successfully"}
					c.Data["json"] = resp
				} else {
					logs.Error("Error occurred updating record ", err.Error())
					var resp = responsesDTOs.StringResponseDTO{StatusCode: 403, Value: "FAILED", StatusDesc: "Error occurred updating record."}
					c.Data["json"] = resp
				}

			} else {
				logs.Debug("OTP has expired. Time to enter OTP of 5 mins exeeded.")
				var resp = responsesDTOs.StringResponseDTO{StatusCode: 403, Value: "FAILED", StatusDesc: "OTP Expired"}
				c.Data["json"] = resp
			}
		} else {
			logs.Debug("OTPs do not match ")
			var resp = responsesDTOs.StringResponseDTO{StatusCode: 402, Value: "FAILED", StatusDesc: "OTP Verification failed"}
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

// Post ...
// @Title Check token expiry
// @Description Check Token Expiry
// @Param	body		body 	requestsDTOs.StringRequestDTO	true		"body for Authentication content"
// @Success 200 {object} responsesDTOs.StringResponseDTO
// @Failure 403 body is empty
// @router /token/check [post]
func (c *AuthenticationController) CheckTokenExpiry() {
	var q requestsDTOs.StringRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	logs.Info("About to verify token ", q.Value)

	if token, err := functions.CheckTokenExpiry(q.Value); err == nil {
		if token.IsValid {
			logs.Info("Token is still valid. User is ", token.User)
			logs.Info("User role is ", token.User.Role.Role)
			// cust, err := models.GetCustomersByUser(token.User)

			// if err != nil {
			// 	c.Data["json"] = err.Error()

			// 	var resp = responsesDTOs.UserResponseDTO{StatusCode: 601, User: nil, StatusDesc: "Error verifying user"}
			// 	c.Data["json"] = resp
			// } else {
			// 	logs.Info("Getting the customer ", cust.Branch.Country.DefaultCurrency.CurrencyId)

			// 	userResp := responsesDTOs.UserResp{
			// 		UserId:        token.User.UserId,
			// 		ImagePath:     token.User.ImagePath,
			// 		UserType:      token.User.UserType,
			// 		FullName:      token.User.FullName,
			// 		Username:      token.User.Username,
			// 		Password:      token.User.Password,
			// 		Email:         token.User.Email,
			// 		PhoneNumber:   token.User.PhoneNumber,
			// 		Gender:        token.User.Gender,
			// 		Dob:           token.User.Dob,
			// 		Address:       token.User.Address,
			// 		IdType:        token.User.IdType,
			// 		IdNumber:      token.User.IdNumber,
			// 		MaritalStatus: token.User.MaritalStatus,
			// 		Active:        token.User.Active,
			// 		Role:          token.User.Role,
			// 		IsVerified:    token.User.IsVerified,
			// 		DateCreated:   token.User.DateCreated,
			// 		DateModified:  token.User.DateModified,
			// 		CreatedBy:     token.User.CreatedBy,
			// 		ModifiedBy:    token.User.ModifiedBy,
			// 		Branch:        cust.Branch,
			// 	}

			// 	var resp = responsesDTOs.UserResponseDTO{StatusCode: 200, User: &userResp, StatusDesc: "Token is valid"}
			// 	c.Data["json"] = resp
			// }

			var resp = responsesDTOs.UserResponseDTO{StatusCode: 200, User: token.User, StatusDesc: "Token is valid"}
			c.Data["json"] = resp
		} else {
			var resp = responsesDTOs.UserResponseDTO{StatusCode: 605, User: nil, StatusDesc: "Invalid token"}
			c.Data["json"] = resp
		}
	} else {
		logs.Error("Error validating token...", err.Error())
		var resp = responsesDTOs.UserResponseDTO{StatusCode: 703, User: nil, StatusDesc: "Error validating token"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// VerifyToken ...
// @Title Verify token
// @Description Verify token
// @Param	body		body 	requestsDTOs.StringRequestDTO	true		"body for Authentication content"
// @Success 200 {object} responsesDTOs.StringResponseDTO
// @Failure 403 body is empty
// @router /token/verify [post]
func (c *AuthenticationController) VerifyToken() {
	var q requestsDTOs.VerifyTokenReq
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	logs.Info("About to verify token ", q.Token)

	statusCode := 608
	message := "Unable to verify token"

	if tokenObj, err := models.GetUserTokensByToken(q.Token); err == nil {
		if plainText, err := functions.GetAESDecrypted(tokenObj.Token, tokenObj.Nonce); err == nil {
			logs.Info("Decrypted token is ", plainText)

			splitText := strings.Split(plainText, "__")

			email := ""
			if len(splitText) > 1 {
				email = splitText[0]
			} else {
				email = splitText[0]
			}

			if user, err := models.GetUsersByUsername(email); err == nil {
				statusCode = 200
				message = "Successfully validated"
				var resp = responsesDTOs.UserResponseDTO{StatusCode: statusCode, User: user, StatusDesc: message}
				c.Data["json"] = resp
			} else {
				statusCode = 708
				message = "user not found"
				var resp = responsesDTOs.UserResponseDTO{StatusCode: statusCode, User: nil, StatusDesc: message}
				c.Data["json"] = resp
			}
		} else {
			logs.Error("Error validating token...", err.Error())
			statusCode = 708
			var resp = responsesDTOs.UserResponseDTO{StatusCode: statusCode, User: nil, StatusDesc: message}
			c.Data["json"] = resp
		}
	} else {
		logs.Error("Error validating token...", err.Error())
		statusCode = 703
		var resp = responsesDTOs.UserResponseDTO{StatusCode: statusCode, User: nil, StatusDesc: message}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// Post ...
// @Title Generate invite token
// @Description Generate invite Token
// @Param	body		body 	requestsDTOs.EncryptInviteRequestDTO	true		"body for Authentication content"
// @Success 200 {object} responsesDTOs.InviteHashResponseDTO
// @Failure 403 body is empty
// @router /token/invite [post]
func (c *AuthenticationController) GenerateInviteToken() {
	var q requestsDTOs.EncryptInviteRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	logs.Info("About to generate token ", q.Email, " and ", q.Role)
	rawString := q.Email + "___" + q.Role

	// ikey, _ := functions.GenerateKey()

	if token, nonce, err := functions.GetAESEncrypted(rawString); err == nil {

		logs.Info("Token generated is ", token)
		logs.Info("Nonce generated is ", nonce)
		// logs.Info("Key generated is ", string(ikey[:]))

		var userToken models.UserTokens = models.UserTokens{Token: token, Nonce: nonce, ExpiryDate: time.Now().Add(4 * time.Hour), Active: 1, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: 1, ModifiedBy: 1}

		_, err := models.AddUserTokens(&userToken)

		if err != nil {
			return
		}

		var encryptResp responsesDTOs.InviteHashDTO = responsesDTOs.InviteHashDTO{Token: &userToken}
		var resp = responsesDTOs.InviteHashResponseDTO{StatusCode: 200, Value: &encryptResp, StatusDesc: "Encrypted"}
		c.Data["json"] = resp

	} else {
		logs.Error("Error validating token...", err.Error())
		var resp = responsesDTOs.InviteHashResponseDTO{StatusCode: 703, Value: nil, StatusDesc: "Error validating token"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// Post ...
// @Title Check invite token validity
// @Description Check Token Expiry
// @Param	body		body 	requestsDTOs.DecryptRequestDTO	true		"body for Authentication content"
// @Success 200 {object} responsesDTOs.InviteDecodeResponseDTO
// @Failure 403 body is empty
// @router /token/invite/verify [post]
func (c *AuthenticationController) VerifyInviteToken() {
	var q requestsDTOs.DecryptRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	logs.Info("About to verify token ", q.Token)

	// ikey, _ := functions.GenerateKey()

	if token, err := functions.GetAESDecrypted(q.Token, q.Nonce); err == nil {
		logs.Info("Token is ", string(token))
		splitToken := strings.Split(string(token), "___")

		logs.Info("Split Token is ", splitToken[0], " and ", splitToken[1])

		var token responsesDTOs.TokenDestructureResponseDTO = responsesDTOs.TokenDestructureResponseDTO{Email: splitToken[0], RoleId: splitToken[1]}

		// if splitToken[0] == q.Email {
		var resp = responsesDTOs.InviteDecodeResponseDTO{StatusCode: 200, Value: &token, StatusDesc: "Successfully verified token."}
		c.Data["json"] = resp
		// } else {
		// 	logs.Error("Error validating token...")
		// 	var resp = responsesDTOs.StringResponseDTO{StatusCode: 703, Value: "", StatusDesc: "Error validating token"}
		// 	c.Data["json"] = resp
		// }

		// if verified, err := functions.VerifyToken(splitToken[1]); err == nil {
		// 	if verified {
		// 		logs.Error("Error validating token...", err.Error())
		// 		var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: splitToken[0], StatusDesc: "Successfully verified token."}
		// 		c.Data["json"] = resp
		// 	} else {
		// 		logs.Error("Error validating token...", err.Error())
		// 		var resp = responsesDTOs.StringResponseDTO{StatusCode: 703, Value: "", StatusDesc: "Error validating token"}
		// 		c.Data["json"] = resp
		// 	}
		// } else {
		// 	logs.Error("Error validating token...", err.Error())
		// 	var resp = responsesDTOs.StringResponseDTO{StatusCode: 703, Value: "", StatusDesc: "Error validating token"}
		// 	c.Data["json"] = resp
		// }
	} else {
		logs.Error("Error validating token...", err.Error())
		var resp = responsesDTOs.InviteDecodeResponseDTO{StatusCode: 703, Value: nil, StatusDesc: "Error validating token"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// Log Out ...
// @Title Log Out
// @Description Logout User
// @Param	body		body 	requestsDTOs.TokenDTO	true		"body for Authentication content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /log-out [post]
func (c *AuthenticationController) Logout() {
	var v requestsDTOs.TokenDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	logs.Info("Received ", v.Token)

	if a, err := models.GetAccessTokensByToken(v.Token); err == nil {
		a.ExpiresAt = time.Now()
		a.Revoked = true
		if err := models.UpdateAccessTokensById(a); err == nil {
			c.Ctx.Output.SetStatus(200)

			var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "", StatusDesc: "User log out complete"}
			c.Data["json"] = resp
		}
	} else {
		logs.Error(err.Error())
		var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Invalid request"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}
