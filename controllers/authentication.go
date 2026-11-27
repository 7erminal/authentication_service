package controllers

import (
	"authentication_service/controllers/functions"
	"authentication_service/internal/trace"
	"authentication_service/models"
	"authentication_service/structs/requestsDTOs"
	"authentication_service/structs/responsesDTOs"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AuthenticationController operations for Authentication
type AuthenticationController struct {
	tracedController
}

// URLMapping ...
func (c *AuthenticationController) URLMapping() {
	c.Mapping("Login", c.Login)
	c.Mapping("VerifyOTP", c.VerifyOTP)
	c.Mapping("ResendOTP", c.ResendOTP)
	c.Mapping("LoginToken", c.LoginToken)
	c.Mapping("GetUserFromToken", c.GetUserFromToken)
	c.Mapping("GenerateInviteToken", c.GenerateInviteToken)
	c.Mapping("VerifyInviteToken", c.VerifyInviteToken)
	c.Mapping("ChangePassword", c.ChangePassword)
	c.Mapping("ResetPassword", c.ResetPassword)
	c.Mapping("SendActivationCode", c.SendActivationCode)
	c.Mapping("VerifyActivationCode", c.VerifyActivationCode)
	c.Mapping("ResetPasswordLink", c.ResetPasswordLink)
	c.Mapping("Logout", c.Logout)
	c.Mapping("VerifyToken", c.VerifyToken)
	c.Mapping("ValidateCustomerCredentialsToken", c.ValidateCustomerCredentialsToken)
	c.Mapping("ExpireCustomerToken", c.ExpireCustomerToken)
	c.Mapping("VerifyCustomerToken", c.VerifyCustomerToken)
	c.Mapping("ChangeCustomerPassword", c.ChangeCustomerPassword)
	c.Mapping("ResetCustomerPassword", c.ResetCustomerPassword)
	c.Mapping("RefreshAccessToken", c.RefreshAccessToken)
	c.Mapping("RefreshCustomerAccessToken", c.RefreshCustomerAccessToken)
}

// Login ...
// @Title Login
// @Description Login User
// @Param	body		body 	models.AuthenticationDTO	true		"body for Authentication content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /login [post]
func (c *AuthenticationController) Login() {
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "Login")
	var v models.AuthenticationDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	if a, err := functions.GetUserWithUsername(&c.Controller, requestsDTOs.GetUserWithUsernameRequest{Username: v.Username}); err == nil {
		// Compare the stored hashed password, with the hashed version of the password that was received
		if a.StatusCode == 200 {
			if err := bcrypt.CompareHashAndPassword([]byte(a.Result.Password), []byte(v.Password)); err != nil {
				// If the two passwords don't match, return a 401 status
				c.Data["json"] = err.Error()

				var resp = responsesDTOs.UserResponseDTO{StatusCode: 605, Result: nil, StatusDesc: "Incorrect password"}
				c.Data["json"] = resp

			} else {
				// cust, err := models.GetCustomersByUser(a)

				// if err != nil {
				// 	c.Data["json"] = err.Error()

				// 	var resp = responsesDTOs.UserResponseDTO{StatusCode: 601, Result: nil, StatusDesc: "Error verifying user"}
				// 	c.Data["json"] = resp
				// } else {

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

				// 	var resp = responsesDTOs.UserResponseDTO{StatusCode: 200, Result: &userResp, StatusDesc: "User has been authenticated"}
				// 	c.Data["json"] = resp
				// }

				c.Ctx.Output.SetStatus(200)

				var resp = responsesDTOs.UserResponseDTO{StatusCode: 200, Result: a.Result, StatusDesc: "User has been authenticated"}
				c.Data["json"] = resp
			}
		} else {
			var resp = responsesDTOs.UserResponseDTO{StatusCode: 302, Result: a.Result, StatusDesc: a.StatusDesc}
			c.Data["json"] = resp
		}
	} else {

		var resp = responsesDTOs.UserResponseDTO{StatusCode: 605, Result: nil, StatusDesc: "Unidentified user"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// LoginToken ...
// @Title Login User
// @Description Login
// @Param	body		body 	models.AuthenticationDTO	true		"body for Authentication content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /login/token [post]
func (c *AuthenticationController) LoginToken() {
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "LoginToken")
	var v models.AuthenticationDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	ipAddress := c.Ctx.Request.RemoteAddr

	statusCode := 400
	statusMessage := "Failed"
	accessTokenObj := &models.AccessTokens{}
	refreshTokenObj := &models.RefreshTokens{}
	// userAgent := c.Ctx.Request.UserAgent()
	accessToken := ""
	refreshToken := ""

	if a, err := functions.GetUserWithUsername(&c.Controller, requestsDTOs.GetUserWithUsernameRequest{Username: strings.Trim(v.Username, " ")}); err == nil && a.StatusCode == 200 {
		// Compare the stored hashed password, with the hashed version of the password that was received

		if a.Result.Active == 1 {
			if err := bcrypt.CompareHashAndPassword([]byte(a.Result.Password), []byte(v.Password)); err != nil {
				// If the two passwords don't match, return a 401 status

				c.Data["json"] = err.Error()

				statusCode = 605
				statusMessage = "Incorrect password"
				var resp = responsesDTOs.UserResponseDTO{StatusCode: statusCode, Result: nil, StatusDesc: statusMessage}
				c.Data["json"] = resp

			} else {
				c.Ctx.Output.SetStatus(200)

				roleStr := a.Result.Role.RoleId
				if roleResp, err := functions.GetRole(&c.Controller, roleStr); err == nil && roleResp.StatusCode == 200 {

					permissions := []responsesDTOs.UserPermission{}

					for _, perm := range roleResp.Role.RolePermissions {
						permissions = append(permissions, responsesDTOs.UserPermission{PermissionCode: perm.Permission.PermissionCode, ActionCode: perm.Action.Action})

					}

					// Create access token (15 minutes expiry)
					useridstr := a.Result.UserId
					branchId := ""
					if a.Result.UserDetails.Branch != nil {
						branchId = a.Result.UserDetails.Branch.BranchId
					}
					token, expiryTime, jti, err := functions.CreateAccessToken(useridstr, v.Username, roleStr, a.Result.Role.Role, permissions, branchId, a.Result.UserDetails.Shop.ShopId)

					if err != nil {

						statusCode = 301
						statusMessage = "Error generating token"
						var resp = responsesDTOs.StringResponseDTO{StatusCode: statusCode, Value: "", StatusDesc: statusMessage}
						c.Data["json"] = resp
					} else {
						// Revoke old tokens for this user
						userID, _ := strconv.ParseInt(a.Result.UserId, 10, 64)
						updateToken := models.AccessTokens{User: userID, Revoked: true}
						if err := models.UpdateAccessTokensByUserIdContext(c.Ctx.Request.Context(), &updateToken); err != nil {

						}

						// time.Now().UTC().Add(time.Hour * 1).Unix()
						t := time.Unix(expiryTime, 0)

						tokenObj := models.AccessTokens{
							User:         userID,
							Token:        jti,
							ExpiresAt:    t,
							DateCreated:  time.Now().UTC(),
							DateModified: time.Now().UTC(),
							LastUsedAt:   time.Now().UTC(),
						}
						if _, err := models.AddAccessTokensContext(c.Ctx.Request.Context(), &tokenObj); err == nil {
							statusCode = 200
							statusMessage = "Access token generated successfully"
							accessTokenObj = &tokenObj
							accessToken = token

							// Create refresh token (7 days)
							refreshToken_, refreshExpiryTime := functions.CreateRefreshToken(v.Username)
							// if err != nil {
							// 	c.Data["json"] = err.Error()
							// 	c.ServeJSON()
							// 	return
							// }

							// Store refresh token
							refreshTokenObj = &models.RefreshTokens{
								User:         userID,
								Token:        functions.HashToken(refreshToken_),
								ExpiresAt:    time.Unix(refreshExpiryTime, 0),
								IPAddress:    ipAddress,
								UserAgent:    "", //userAgent,
								AccessToken:  accessTokenObj,
								DateCreated:  time.Now().UTC(),
								DateModified: time.Now().UTC(),
							}

							if _, err := models.AddRefreshTokensContext(c.Ctx.Request.Context(), refreshTokenObj); err != nil {

								c.Data["json"] = err.Error()
								c.ServeJSON()
								return
							}
							refreshToken = refreshToken_

						} else {

							statusCode = 301
							statusMessage = "Error generating token"
						}
					}
				} else {

					statusCode = 604
					statusMessage = "Role not found"
				}
			}
		} else {

			statusCode = 607
			statusMessage = "Inactive user"
		}
	} else {

		statusCode = 605
		statusMessage = "Unidentified user"
	}
	var tokenResponse = responsesDTOs.TokenResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}

	userType := "USER"

	result := responsesDTOs.LoginDataResponseDTO{
		UserType: userType,
		Token:    &tokenResponse,
	}

	var resp = responsesDTOs.LoginTokenResponseDTO{StatusCode: statusCode, StatusDesc: statusMessage, Result: &result}
	c.Data["json"] = resp
	c.ServeJSON()
}

// Refresh access token ...
// @Title Refresh Access Token
// @Description Refresh Access Token
// @Param	body		body 	models.AuthenticationDTO	true		"body for Authentication content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /refresh/user/token [post]
func (c *AuthenticationController) RefreshAccessToken() {
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "RefreshAccessToken")
	refreshToken := c.Ctx.Input.Header("RefreshToken")

	// Validate refresh token
	if refreshTokenObj, err := models.GetRefreshTokensByTokenContext(c.Ctx.Request.Context(), functions.HashToken(refreshToken)); err == nil {

		if refreshTokenObj.ExpiresAt.Unix() > time.Now().UTC().Unix() && !refreshTokenObj.Revoked {
			if userResp, err := functions.GetUser(&c.Controller, requestsDTOs.GetUserRequest{UserId: strconv.Itoa(int(refreshTokenObj.User))}); err == nil {
				if userResp.StatusCode == 200 {
					// Create new access token

					username := userResp.Result.Username
					if username == "" {

						username = userResp.Result.Email

					}
					roleStr := userResp.Result.Role.RoleId
					if roleResp, err := functions.GetRole(&c.Controller, roleStr); err == nil && roleResp.StatusCode == 200 {

						permissions := []responsesDTOs.UserPermission{}

						for _, perm := range roleResp.Role.RolePermissions {
							permissions = append(permissions, responsesDTOs.UserPermission{PermissionCode: perm.Permission.PermissionCode, ActionCode: perm.Action.Action})

						}
						useridstr := userResp.Result.UserId
						accessToken_, accessExpiryTime, jti, err := functions.CreateAccessToken(useridstr, username, roleStr, roleResp.Role.Role, permissions, userResp.Result.UserDetails.Branch.BranchId, userResp.Result.UserDetails.Shop.ShopId)
						if err != nil {
							c.Data["json"] = err.Error()
							c.ServeJSON()
							return
						}

						accessTokenObj := models.AccessTokens{
							User:         refreshTokenObj.User,
							Token:        jti,
							ExpiresAt:    time.Unix(accessExpiryTime, 0).UTC(),
							IPAddress:    c.Ctx.Request.RemoteAddr,
							DateCreated:  time.Now().UTC(),
							DateModified: time.Now().UTC(),
							LastUsedAt:   time.Now().UTC(),
						}

						if _, err := models.AddAccessTokensContext(c.Ctx.Request.Context(), &accessTokenObj); err != nil {

							c.Data["json"] = err.Error()
							c.ServeJSON()
							return
						}

						var tokenResponse = responsesDTOs.TokenResponseDTO{
							AccessToken:  accessToken_,
							RefreshToken: refreshToken, // Return same refresh token
							TokenType:    "Bearer",
							ExpiresIn:    900,
						}

						userType := "USER"

						result := responsesDTOs.LoginDataResponseDTO{
							UserType: userType,
							Token:    &tokenResponse,
						}

						var resp = responsesDTOs.LoginTokenResponseDTO{StatusCode: 200, StatusDesc: "Access token generated successfully", Result: &result}

						c.Data["json"] = resp
					} else {

						var resp = responsesDTOs.StringResponseDTO{
							StatusCode: 608,
							Value:      "",
							StatusDesc: "Role not found",
						}
						c.Data["json"] = resp
					}
				} else {

					c.Ctx.Output.SetStatus(401)
					var resp = responsesDTOs.StringResponseDTO{
						StatusCode: 605,
						Value:      "",
						StatusDesc: "User not found",
					}
					c.Data["json"] = resp
				}
			} else {

				c.Ctx.Output.SetStatus(401)
				var resp = responsesDTOs.StringResponseDTO{
					StatusCode: 605,
					Value:      "",
					StatusDesc: "Error fetching user details",
				}
				c.Data["json"] = resp
			}
		} else {

			c.Ctx.Output.SetStatus(401)
			var resp = responsesDTOs.StringResponseDTO{
				StatusCode: 605,
				Value:      "",
				StatusDesc: "Refresh token expired or revoked",
			}

			c.Data["json"] = resp
		}
	} else {

		c.Ctx.Output.SetStatus(401)
		var resp = responsesDTOs.StringResponseDTO{
			StatusCode: 605,
			Value:      "",
			StatusDesc: "Invalid refresh token",
		}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// ValidateCustomerCredentials ...
// @Title Validate Customer Credentials
// @Description Validate Customer Credentials
// @Param	body		body 	models.AuthenticationDTO	true		"body for Authentication content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /validate-customer-credentials/token [post]
func (c *AuthenticationController) ValidateCustomerCredentialsToken() {
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "ValidateCustomerCredentialsToken")
	var v models.AuthenticationDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	ipAddress := c.Ctx.Request.RemoteAddr

	trimUsername := strings.TrimSpace(v.Username)

	statusCode := 400
	statusMessage := "Failed"
	accessTokenObj := &models.Customer_access_tokens{}
	refreshTokenObj := &models.CustomerRefreshTokens{}
	accessToken := ""
	refreshToken := ""

	if a, err := models.GetCustomer_credentialsByCustomerUsernameContext(c.Ctx.Request.Context(), trimUsername); err == nil {
		// Compare the stored hashed password, with the hashed version of the password that was received

		if a.Active == 1 {
			if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(v.Password)); err != nil {
				// If the two passwords don't match, return a 401 status

				c.Data["json"] = err.Error()

				statusCode = 605
				statusMessage = "Incorrect password"
			} else {

				custId := requestsDTOs.GetCustomerRequest{CustomerId: strconv.FormatInt(a.Customer, 10)}
				if customer, err := functions.GetCustomer(&c.Controller, custId); err != nil {

					statusCode = 401
					statusMessage = "Customer not found"
				} else {

					c.Ctx.Output.SetStatus(200)

					token, expiryTime, jti, err := functions.CreateCustomerAccessToken(strconv.FormatInt(a.Customer, 10), v.Username, "customer", customer.Result.Branch.BranchId, customer.Result.Shop.ShopId)

					if err != nil {

						var resp = responsesDTOs.StringResponseDTO{StatusCode: 301, Value: "", StatusDesc: "Error generating token"}
						c.Data["json"] = resp
					} else {
						updateToken := models.Customer_access_tokens{Customer: a.Customer, Revoked: true}
						if err := models.UpdateCustomer_access_tokensByCustomerContext(c.Ctx.Request.Context(), &updateToken); err != nil {
							statusMessage = fmt.Sprintf("Error revoking old tokens. %s", err.Error())

						}

						t := time.Unix(expiryTime, 0)
						accessTokenObj = &models.Customer_access_tokens{Customer: a.Customer, Token: jti, ExpiresAt: t, DateCreated: time.Now(), LastUsedAt: time.Now()}
						if _, err := models.AddCustomer_access_tokensContext(c.Ctx.Request.Context(), accessTokenObj); err == nil {

							statusCode = 200
							statusMessage = "Access token generated successfully"
							accessToken = token

							// Create refresh token (7 days)
							refreshToken_, refreshExpiryTime := functions.CreateRefreshToken(v.Username)

							// Store refresh token
							refreshTokenObj = &models.CustomerRefreshTokens{
								Customer:     a.Customer,
								Token:        functions.HashToken(refreshToken_),
								ExpiresAt:    time.Unix(refreshExpiryTime, 0),
								IPAddress:    ipAddress,
								UserAgent:    "", //userAgent,
								AccessToken:  accessTokenObj,
								DateCreated:  time.Now(),
								DateModified: time.Now(),
							}

							if _, err := models.AddCustomerRefreshTokensContext(c.Ctx.Request.Context(), refreshTokenObj); err != nil {

								statusCode = 301
								statusMessage = "Error generating token"
							} else {

								refreshToken = refreshToken_
								statusCode = 200
								statusMessage = "Tokens generated successfully"
							}

						} else {

							statusCode = 301
							statusMessage = "Error generating token"
						}

					}
				}
			}
		} else {

			var resp = responsesDTOs.StringResponseDTO{StatusCode: 607, Value: "", StatusDesc: "Inactive Customer"}
			c.Data["json"] = resp
		}
	} else {

		statusCode = 605
		statusMessage = "Unidentified customer"
	}

	var tokenResponse = responsesDTOs.TokenResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}

	userType := "CUSTOMER"

	result := responsesDTOs.LoginDataResponseDTO{
		UserType: userType,
		Token:    &tokenResponse,
	}

	var resp = responsesDTOs.LoginTokenResponseDTO{StatusCode: statusCode, StatusDesc: statusMessage, Result: &result}
	c.Data["json"] = resp
	c.ServeJSON()
}

// Refresh customer access token ...
// @Title Refresh Customer Access Token
// @Description Refresh Access Token
// @Param	body		body 	models.AuthenticationDTO	true		"body for Authentication content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /refresh/customer/token [post]
func (c *AuthenticationController) RefreshCustomerAccessToken() {
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "RefreshCustomerAccessToken")
	refreshToken := c.Ctx.Input.Header("RefreshToken")

	// Validate refresh token
	if refreshTokenObj, err := models.GetCustomerRefreshTokensByTokenContext(c.Ctx.Request.Context(), refreshToken); err == nil {
		if refreshTokenObj.ExpiresAt.After(time.Now().UTC()) && !refreshTokenObj.Revoked {
			// Create new access token
			if customerResp, err := functions.GetCustomer(&c.Controller, requestsDTOs.GetCustomerRequest{CustomerId: strconv.Itoa(int(refreshTokenObj.Customer))}); err == nil {
				if customerResp.StatusCode == 200 {

					accessToken, accessExpiryTime, jti, err := functions.CreateCustomerAccessToken(customerResp.Result.CustomerNumber, "customer", customerResp.Result.CustomerCategory.CustomerCategoryId, customerResp.Result.Branch.BranchId, customerResp.Result.Shop.ShopId)
					if err != nil {
						c.Data["json"] = err.Error()
						c.ServeJSON()
						return
					}

					accessTokenObj := models.Customer_access_tokens{
						Customer:    refreshTokenObj.Customer,
						Token:       jti,
						ExpiresAt:   time.Unix(accessExpiryTime, 0),
						Revoked:     false,
						DateCreated: time.Now(),
						LastUsedAt:  time.Now(),
					}

					if _, err := models.AddCustomer_access_tokensContext(c.Ctx.Request.Context(), &accessTokenObj); err != nil {
						c.Data["json"] = err.Error()
						c.ServeJSON()
						return
					}

					var resp = responsesDTOs.TokenResponseDTO{
						AccessToken:  accessToken,
						RefreshToken: refreshToken, // Return same refresh token
						TokenType:    "Bearer",
						ExpiresIn:    900,
					}

					c.Data["json"] = resp
				} else {

					c.Ctx.Output.SetStatus(200)

					var resp = responsesDTOs.StringResponseDTO{
						StatusCode: 605,
						Value:      "",
						StatusDesc: "Customer not found",
					}
					c.Data["json"] = resp
				}
			} else {

				c.Ctx.Output.SetStatus(200)

				var resp = responsesDTOs.StringResponseDTO{
					StatusCode: 605,
					Value:      "",
					StatusDesc: "Customer not found",
				}
				c.Data["json"] = resp
			}

		} else {
			c.Ctx.Output.SetStatus(200)
			var resp = responsesDTOs.StringResponseDTO{
				StatusCode: 605,
				Value:      "",
				StatusDesc: "Refresh token expired or revoked",
			}
			c.Data["json"] = resp
		}
	} else {
		c.Ctx.Output.SetStatus(401)
		var resp = responsesDTOs.StringResponseDTO{
			StatusCode: 605,
			Value:      "",
			StatusDesc: "Invalid refresh token",
		}
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
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "ChangePassword")
	idStr := c.Ctx.Input.Param(":id")

	var v requestsDTOs.ChangePassword
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	if a, err := functions.GetUser(&c.Controller, requestsDTOs.GetUserRequest{UserId: idStr}); err == nil {
		if a.StatusCode == 200 {
			// Compare the stored hashed password, with the hashed version of the password that was received
			if err := bcrypt.CompareHashAndPassword([]byte(a.Result.Password), []byte(v.OldPassword)); err != nil {
				// If the two passwords don't match, return a 401 status
				c.Data["json"] = err.Error()

				var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Old password does not match"}
				c.Data["json"] = resp

			} else {
				hashedPassword, errr := bcrypt.GenerateFromPassword([]byte(v.NewPassword), 8)

				if errr == nil {

					a.Result.Password = string(hashedPassword)

					// models.Agents{AgentName: v.AgentName, BranchId: v.BranchId, IdType: v.IdType, IdNumber: v.IdNumber, IsVerified: false, Active: 1, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: c_by, ModifiedBy: c_by}
				} else {

				}

				if _, err := functions.UpdateUserPassword(&c.Controller, requestsDTOs.UpdateUserPasswordRequest{
					UserId:      idStr,
					OldPassword: v.OldPassword,
					NewPassword: v.NewPassword,
				}); err == nil {
					c.Ctx.Output.SetStatus(200)

					var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "Successfully changed password", StatusDesc: "User password has been changed successfully"}
					c.Data["json"] = resp
				} else {
					var resp = responsesDTOs.StringResponseDTO{StatusCode: 608, Value: "", StatusDesc: "User password change failed. " + err.Error()}
					c.Data["json"] = resp
				}
			}
		} else {

			var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Unidentified user"}
			c.Data["json"] = resp
		}
	} else {

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
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "ResetPassword")
	idStr := c.Ctx.Input.Param(":id")
	// id, _ := strconv.ParseInt(idStr, 0, 64)

	var v requestsDTOs.ResetPassword
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	if a, err := functions.GetUser(&c.Controller, requestsDTOs.GetUserRequest{UserId: idStr}); err == nil {
		// Compare the stored hashed password, with the hashed version of the password that was received

		if a.StatusCode == 200 {
			hashedPassword, errr := bcrypt.GenerateFromPassword([]byte(v.NewPassword), 8)

			if errr == nil {

				a.Result.Password = string(hashedPassword)

				// models.Agents{AgentName: v.AgentName, BranchId: v.BranchId, IdType: v.IdType, IdNumber: v.IdNumber, IsVerified: false, Active: 1, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: c_by, ModifiedBy: c_by}
			} else {

			}

			if _, err := functions.UpdateUserPassword(&c.Controller, requestsDTOs.UpdateUserPasswordRequest{
				UserId:      idStr,
				OldPassword: a.Result.Password,
				NewPassword: v.NewPassword,
			}); err == nil {
				c.Ctx.Output.SetStatus(200)

				var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "Successfully changed password", StatusDesc: "User password has been changed successfully"}
				c.Data["json"] = resp
			} else {
				var resp = responsesDTOs.StringResponseDTO{StatusCode: 608, Value: "", StatusDesc: "User password change failed. " + err.Error()}
				c.Data["json"] = resp
			}
		} else {

			var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Unidentified user"}
			c.Data["json"] = resp
		}

	} else {

		var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Unidentified user"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// Reset Customer Password ...
// @Title Reset Customer Password
// @Description Reset customer password
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	requestsDTOs.ResetPassword	true		"body for Change password content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /reset-customer-password/:id [put]
func (c *AuthenticationController) ResetCustomerPassword() {
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "ResetCustomerPassword")
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)

	var v requestsDTOs.ResetPassword
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	req := requestsDTOs.GetCustomerRequest{CustomerId: strconv.Itoa(int(id))}
	if a, err := functions.GetCustomer(&c.Controller, req); err == nil {
		// Compare the stored hashed password, with the hashed version of the password that was received

		customerID, _ := strconv.ParseInt(a.Result.CustomerId, 10, 64)
		if custCred, err := models.GetCustomer_credentialsByCustomerIdContext(c.Ctx.Request.Context(), customerID); err == nil {

			hashedPassword, errr := bcrypt.GenerateFromPassword([]byte(v.NewPassword), 8)

			if errr == nil {

				custCred.Password = string(hashedPassword)

				// models.Agents{AgentName: v.AgentName, BranchId: v.BranchId, IdType: v.IdType, IdNumber: v.IdNumber, IsVerified: false, Active: 1, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: c_by, ModifiedBy: c_by}
			} else {

			}

			if err := models.UpdateCustomer_credentialsByIdContext(c.Ctx.Request.Context(), custCred); err == nil {

				c.Ctx.Output.SetStatus(200)

				var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "Successfully changed password", StatusDesc: "Customer password has been changed successfully"}
				c.Data["json"] = resp
			} else {
				var resp = responsesDTOs.StringResponseDTO{StatusCode: 608, Value: "", StatusDesc: "Customer password reset failed. " + err.Error()}
				c.Data["json"] = resp
			}
		} else {

			var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Unidentified customer"}
			c.Data["json"] = resp
		}

	} else {

		var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Unidentified Customer"}
		c.Data["json"] = resp
	}

	c.ServeJSON()
}

// Change Customer Password ...
// @Title Change Customer Password
// @Description Change customer password
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	requestsDTOs.ChangePassword	true		"body for Change password content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /change-customer-password/:id [put]
func (c *AuthenticationController) ChangeCustomerPassword() {
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "ChangeCustomerPassword")
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)

	var v requestsDTOs.ChangePassword
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	if a, err := functions.GetCustomer(&c.Controller, requestsDTOs.GetCustomerRequest{CustomerId: strconv.Itoa(int(id))}); err == nil {
		// Compare the stored hashed password, with the hashed version of the password that was received

		if a.StatusCode == 200 {
			customerID, _ := strconv.ParseInt(a.Result.CustomerId, 10, 64)
			if custCred, err := models.GetCustomer_credentialsByCustomerIdContext(c.Ctx.Request.Context(), customerID); err == nil {
				if err := bcrypt.CompareHashAndPassword([]byte(custCred.Password), []byte(v.OldPassword)); err != nil {
					// If the two passwords don't match, return a 401 status
					c.Data["json"] = err.Error()

					var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Old password does not match"}
					c.Data["json"] = resp

				} else {
					hashedPassword, errr := bcrypt.GenerateFromPassword([]byte(v.NewPassword), 8)

					if errr == nil {

						custCred.Password = string(hashedPassword)

						// models.Agents{AgentName: v.AgentName, BranchId: v.BranchId, IdType: v.IdType, IdNumber: v.IdNumber, IsVerified: false, Active: 1, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: c_by, ModifiedBy: c_by}
					} else {

					}

					if err := models.UpdateCustomer_credentialsByIdContext(c.Ctx.Request.Context(), custCred); err == nil {
						c.Ctx.Output.SetStatus(200)

						var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "Successfully changed password", StatusDesc: "Customer password has been changed successfully"}
						c.Data["json"] = resp
					} else {
						var resp = responsesDTOs.StringResponseDTO{StatusCode: 608, Value: "", StatusDesc: "Customer password reset failed. " + err.Error()}
						c.Data["json"] = resp
					}
				}

			} else {

				var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Unidentified customer"}
				c.Data["json"] = resp
			}
		} else {

			var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Unidentified Customer"}
			c.Data["json"] = resp
		}

	} else {

		var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Unidentified Customer"}
		c.Data["json"] = resp
	}

	c.ServeJSON()
}

// AddCustomerCredential ...
// @Title Post
// @Description create Customer_credentials
// @Param	body		body 	requests.CustomerCredentialRequestDTO	true		"body for Customer_credentials content"
// @Success 201 {int} models.Customer_credentials
// @Failure 403 body is empty
// @router /add-customer-credential [post]
func (c *AuthenticationController) AddCustomerCredential() {
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "AddCustomerCredential")
	var v requestsDTOs.CustomerCredentialRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	statusCode := 401
	statusDesc := "Unauthorized"

	if customer, err := functions.GetCustomer(&c.Controller, requestsDTOs.GetCustomerRequest{CustomerId: v.CustomerId}); err == nil {
		if customer.StatusCode == 200 {
			hashedPassword, errr := bcrypt.GenerateFromPassword([]byte(v.Password), 8)
			if errr != nil {

				c.Data["json"] = &responsesDTOs.CustomerCredentialsResponseDTO{
					StatusCode: 500,
					StatusDesc: "Internal Server Error",
					Result:     "Error hashing password",
				}
				c.Ctx.Output.SetStatus(500)
				return
			}

			// hashedPin, errr := bcrypt.GenerateFromPassword([]byte(v.Pin), 8)
			// if errr != nil {
			// 	c.Data["json"] = &responses.CustomerCredentialsResponseDTO{
			// 		StatusCode: 500,
			// 		StatusDesc: "Internal Server Error",
			// 		Result:     "Error hashing password",
			// 	}
			// 	c.Ctx.Output.SetStatus(500)
			// 	return
			// }

			custIdInt, _ := strconv.ParseInt(v.CustomerId, 10, 64)
			ccredential := models.Customer_credentials{
				Customer:     custIdInt,
				Username:     v.Username,
				Password:     string(hashedPassword),
				Pin:          v.Pin,
				DateCreated:  time.Now(),
				DateModified: time.Now(),
				CreatedBy:    1,
				ModifiedBy:   1,
				Active:       1,
			}

			if _, err := models.AddCustomer_credentialsContext(c.Ctx.Request.Context(), &ccredential); err == nil {
				c.Ctx.Output.SetStatus(200)

				statusCode = 200
				statusDesc = "Customer credential added successfully"

				c.Data["json"] = &responsesDTOs.CustomerCredentialsResponseDTO{
					StatusCode: statusCode,
					StatusDesc: statusDesc,
					Result:     "Customer credential added successfully",
				}
				// c.Data["json"] = v
			} else {
				c.Data["json"] = err.Error()

				c.Ctx.Output.SetStatus(200)

				c.Data["json"] = &responsesDTOs.CustomerCredentialsResponseDTO{
					StatusCode: statusCode,
					StatusDesc: statusDesc,
					Result:     err.Error(),
				}
			}
		} else {

			c.Data["json"] = &responsesDTOs.CustomerCredentialsResponseDTO{
				StatusCode: statusCode,
				StatusDesc: statusDesc,
				Result:     "Customer not found",
			}
			c.Ctx.Output.SetStatus(200)
		}
	} else {

		c.Data["json"] = &responsesDTOs.CustomerCredentialsResponseDTO{
			StatusCode: statusCode,
			StatusDesc: statusDesc,
			Result:     "Error fetching customer",
		}
		c.Ctx.Output.SetStatus(200)
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Customer_credentials
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	requests.CustomerCredentialUpdateRequestDTO	true		"body for Customer_credentials content"
// @Success 200 {object} models.Customer_credentials
// @Failure 403 :customerId is not int
// @router /update-customer-credential/:customerId [put]
func (c *AuthenticationController) Put() {
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "Put")
	idStr := c.Ctx.Input.Param(":customerId")
	// id, _ := strconv.ParseInt(idStr, 0, 64)
	// v := models.Customer_credentials{Id: id}

	var v requestsDTOs.CustomerCredentialUpdateRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	statusCode := 401
	statusDesc := "Unauthorized"

	if customer, err := functions.GetCustomer(&c.Controller, requestsDTOs.GetCustomerRequest{CustomerId: idStr}); err == nil {

		if customer.StatusCode == 200 {

			// c.Data["json"] = v
		} else {

			c.Data["json"] = &responsesDTOs.CustomerCredentialsResponseDTO{
				StatusCode: statusCode,
				StatusDesc: statusDesc,
				Result:     "Customer not found",
			}
			c.Ctx.Output.SetStatus(200)
			c.ServeJSON()
			return
		}

		customerID, _ := strconv.ParseInt(customer.Result.CustomerId, 10, 64)
		if customerC, err := models.GetCustomer_credentialsByCustomerIdContext(c.Ctx.Request.Context(), customerID); err == nil {

			hashedPassword, errr := bcrypt.GenerateFromPassword([]byte(v.Password), 8)
			if errr != nil {

				c.Data["json"] = &responsesDTOs.CustomerCredentialsResponseDTO{
					StatusCode: 500,
					StatusDesc: "Internal Server Error",
					Result:     "Error hashing password",
				}
				c.Ctx.Output.SetStatus(500)
				return
			}

			// hashedPin, errr := bcrypt.GenerateFromPassword([]byte(v.Pin), 8)
			// if errr != nil {
			// 	c.Data["json"] = &responses.CustomerCredentialsResponseDTO{
			// 		StatusCode: 500,
			// 		StatusDesc: "Internal Server Error",
			// 		Result:     "Error hashing password",
			// 	}
			// 	c.Ctx.Output.SetStatus(500)
			// 	return
			// }

			ccredential := models.Customer_credentials{
				Id:           customerC.Id,
				Customer:     customerID,
				Username:     v.Username,
				Password:     string(hashedPassword),
				Pin:          v.Pin,
				DateCreated:  time.Now(),
				DateModified: time.Now(),
				CreatedBy:    1,
				ModifiedBy:   1,
				Active:       1,
			}

			// json.Unmarshal(c.Ctx.Input.RequestBody, &v)
			if err := models.UpdateCustomer_credentialsByIdContext(c.Ctx.Request.Context(), &ccredential); err == nil {
				c.Data["json"] = "OK"

				c.Ctx.Output.SetStatus(200)
				statusCode = 200
				statusDesc = "Customer credential updated successfully"

				c.Data["json"] = &responsesDTOs.CustomerCredentialsResponseDTO{
					StatusCode: statusCode,
					StatusDesc: statusDesc,
					Result:     "Customer credential updated successfully",
				}
			} else {
				// c.Data["json"] = err.Error()

				c.Ctx.Output.SetStatus(200)
				c.Data["json"] = &responsesDTOs.CustomerCredentialsResponseDTO{
					StatusCode: statusCode,
					StatusDesc: statusDesc,
					Result:     err.Error(),
				}
			}
		} else {

			c.Data["json"] = &responsesDTOs.CustomerCredentialsResponseDTO{
				StatusCode: statusCode,
				StatusDesc: statusDesc,
				Result:     "Error fetching customer credential",
			}
			c.Ctx.Output.SetStatus(200)
		}
	} else {

		c.Data["json"] = &responsesDTOs.CustomerCredentialsResponseDTO{
			StatusCode: statusCode,
			StatusDesc: statusDesc,
			Result:     "Error fetching customer",
		}
		c.Ctx.Output.SetStatus(200)
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
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "ResetPasswordLink")
	var v requestsDTOs.ResetPasswordLink
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	req := requestsDTOs.GetUserWithUsernameRequest{Username: v.Email}
	if a, err := functions.GetUserWithUsername(&c.Controller, req); err == nil {
		// Compare the stored hashed password, with the hashed version of the password that was received

		// hashedPassword, errr := bcrypt.GenerateFromPassword([]byte(v.NewPassword), 8)

		rawString := v.Email + "___" + a.Result.Role.Role

		// ikey, _ := functions.GenerateKey()

		if token, nonce, err := functions.GetAESEncrypted(rawString); err == nil {

			var userToken models.UserTokens = models.UserTokens{Token: token, Nonce: nonce, ExpiryDate: time.Now().Add(4 * time.Hour), Active: 1, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: 1, ModifiedBy: 1}

			_, err := models.AddUserTokensContext(c.Ctx.Request.Context(), &userToken)

			if err != nil {
				return
			}

			namePlaceHolder := strings.Split(a.Result.FullName, " | ")
			name := strings.Join(namePlaceHolder, " ")
			message_ := strings.Replace(v.Message, "[SENDER_NAME_ID]", name, -1)

			for i, link := range v.Links {
				iStr := strconv.Itoa(i)
				placeholder := "[LINK_" + iStr + "_ID]"
				formattedLink := *link + token
				message_ = strings.Replace(message_, placeholder, formattedLink, -1)

			}

			go functions.SendEmailNewContext(trace.Background(c.Ctx.Request.Context()), a.Result.Email, v.Subject, message_)

		} else {

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
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "VerifyOTP")
	// username := c.Ctx.Input.Param(":username")
	var q requestsDTOs.VerifyOtpDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	req := requestsDTOs.GetUserWithUsernameRequest{Username: q.Username}
	v, err := functions.GetUserWithUsername(&c.Controller, req)

	if err != nil {
		var resp = responsesDTOs.UserResponseDTO{StatusCode: 604, Result: nil, StatusDesc: "User cannot be found"}
		// c.Data["json"] = err.Error()
		c.Data["json"] = resp
	} else {
		// Get OTP
		if v.StatusCode == 200 {

			userID, _ := strconv.ParseInt(v.Result.UserId, 10, 64)
			otp, err := models.VerifyUserOTPContext(c.Ctx.Request.Context(), userID)

			if err == nil {
				if q.Password == otp.Code {

					if otp.ExpiryDate.After(time.Now()) {
						if otp.Status == 1 {

							var resp = responsesDTOs.UserResponseDTO{StatusCode: 407, Result: nil, StatusDesc: "OTP has already been used."}
							c.Data["json"] = resp
						} else {
							otp.Status = 1

							if err := models.UpdateUserOtpByIdContext(c.Ctx.Request.Context(), otp); err == nil {
								var resp = responsesDTOs.UserResponseDTO{StatusCode: 200, Result: nil, StatusDesc: "OTP Verified successfully"}
								c.Data["json"] = resp
							} else {

								var resp = responsesDTOs.UserResponseDTO{StatusCode: 403, Result: nil, StatusDesc: "Error occurred inserting record."}
								c.Data["json"] = resp
							}
						}
					} else {

						var resp = responsesDTOs.UserResponseDTO{StatusCode: 403, Result: nil, StatusDesc: "OTP Expired"}
						c.Data["json"] = resp
					}
				} else {

					var resp = responsesDTOs.UserResponseDTO{StatusCode: 402, Result: nil, StatusDesc: "OTP Verification failed"}
					c.Data["json"] = resp
				}
			} else {

				var resp = responsesDTOs.UserResponseDTO{StatusCode: 403, Result: nil, StatusDesc: "OTP Expired"}
				c.Data["json"] = resp
			}
		} else {

			var resp = responsesDTOs.UserResponseDTO{StatusCode: 403, Result: nil, StatusDesc: "OTP Expired"}
			c.Data["json"] = resp
		}
		// Generate random number
		// randNum := functions.EncodeToString(6)

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
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "ResendOTP")
	// username := c.Ctx.Input.Param(":username")
	var q requestsDTOs.UsernameDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	req := requestsDTOs.GetUserWithUsernameRequest{Username: q.Username}
	v, err := functions.GetUserWithUsername(&c.Controller, req)

	if err != nil {
		var resp = responsesDTOs.UserResponseDTO{StatusCode: 604, Result: nil, StatusDesc: "User cannot be found"}
		// c.Data["json"] = err.Error()
		c.Data["json"] = resp
	} else {
		// Generate random number
		if v.StatusCode == 200 {
			randNum := functions.EncodeToString(6)

			expiryDate := time.Now().Local().Add(time.Hour*time.Duration(0) + time.Minute*time.Duration(5) + time.Second*time.Duration(0))

			userID, _ := strconv.ParseInt(v.Result.UserId, 10, 64)
			otpModel := models.UserOtps{Code: randNum, UserId: userID, Status: 2, DateCreated: time.Now(), DateModified: time.Now(), DateGenerated: time.Now(), ExpiryDate: expiryDate, Active: 1}

			if _, err := models.AddUserOtpContext(c.Ctx.Request.Context(), &otpModel); err == nil {
				functions.SendEmailContext(c.Ctx.Request.Context(), v.Result.Email, randNum)

				var resp = responsesDTOs.UserResponseDTO{StatusCode: 200, Result: nil, StatusDesc: "Email sent successfully"}
				c.Data["json"] = resp
			} else {

				var resp = responsesDTOs.UserResponseDTO{StatusCode: 703, Result: nil, StatusDesc: "Error sending email"}
				c.Data["json"] = resp
			}
		} else {

			var resp = responsesDTOs.UserResponseDTO{StatusCode: 605, Result: nil, StatusDesc: "Unidentified user"}
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
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "SendActivationCode")
	// username := c.Ctx.Input.Param(":username")
	var q requestsDTOs.SendActivationCode
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	// Generate random number
	randNum := functions.EncodeToString(6)
	// set OTP to 1111 for tests
	randNum = "1111"

	proceed := false

	if ac, err := models.GetActivationCodesByNumberContext(c.Ctx.Request.Context(), q.MobileNumber); err == nil {
		for _, suc := range ac {

			singleAc, err := suc.(models.ActivationCodes)

			if !err {
				singleAc.ExpiryDate = time.Now()
				models.UpdateActivationCodesByIdContext(c.Ctx.Request.Context(), &singleAc)
			}
		}
		proceed = true
	} else {
		proceed = true
	}

	if proceed {
		expiryDate := time.Now().Local().Add(time.Hour*time.Duration(1) + time.Minute*time.Duration(0) + time.Second*time.Duration(0))

		otpModel := models.ActivationCodes{Code: randNum, Number: q.MobileNumber, DateCreated: time.Now(), DateModified: time.Now(), ExpiryDate: expiryDate, Active: 1}

		if _, err := models.AddActivationCodesContext(c.Ctx.Request.Context(), &otpModel); err == nil {
			// Function to send Code via sms
			// functions.SendEmail(v.Email, randNum)

			var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "SUCCESS", StatusDesc: "Email sent successfully"}
			c.Data["json"] = resp
		} else {

			var resp = responsesDTOs.StringResponseDTO{StatusCode: 703, Value: "FAILED", StatusDesc: "Error sending email"}
			c.Data["json"] = resp
		}
	} else {

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
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "VerifyActivationCode")
	// username := c.Ctx.Input.Param(":username")
	var q requestsDTOs.VerifyActivationCodeDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	if v, err := models.GetActivationCodeByNumberContext(c.Ctx.Request.Context(), q.MobileNumber); err != nil {

		var resp = responsesDTOs.StringResponseDTO{StatusCode: 604, Value: "FAILED", StatusDesc: "Code cannot be found. It either does not exist or has expired"}
		// c.Data["json"] = err.Error()
		c.Data["json"] = resp
	} else {
		// Get OTP

		if q.Password == v.Code {

			if v.ExpiryDate.After(time.Now()) {

				v.ExpiryDate = time.Now()

				if err := models.UpdateActivationCodesByIdContext(c.Ctx.Request.Context(), v); err == nil {
					var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "SUCCESS", StatusDesc: "OTP Verified successfully"}
					c.Data["json"] = resp
				} else {

					var resp = responsesDTOs.StringResponseDTO{StatusCode: 403, Value: "FAILED", StatusDesc: "Error occurred updating record."}
					c.Data["json"] = resp
				}

			} else {

				var resp = responsesDTOs.StringResponseDTO{StatusCode: 403, Value: "FAILED", StatusDesc: "OTP Expired"}
				c.Data["json"] = resp
			}
		} else {

			var resp = responsesDTOs.StringResponseDTO{StatusCode: 402, Value: "FAILED", StatusDesc: "OTP Verification failed"}
			c.Data["json"] = resp
		}

		// Generate random number
		// randNum := functions.EncodeToString(6)

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

// VerifyTokenExpiry ...
// @Title Verify token expiry
// @Description Verify Token Expiry
// @Param	body		body 	requestsDTOs.StringRequestDTO	true		"body for Authentication content"
// @Success 200 {object} responsesDTOs.StringResponseDTO
// @Failure 403 body is empty
// @router /token/verify [post]
func (c *AuthenticationController) VerifyToken() {
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "VerifyToken")

	var q requestsDTOs.StringRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	if token, err := functions.CheckTokenExpiry(q.Value); err == nil {
		if token.IsValid {

			var resp = responsesDTOs.AccessVerificationResponseDTO{StatusCode: 200, Result: token.User, StatusDesc: "Token is valid"}
			c.Data["json"] = resp
		} else {
			var resp = responsesDTOs.AccessVerificationResponseDTO{StatusCode: 605, Result: nil, StatusDesc: "Invalid token"}
			c.Data["json"] = resp
		}
	} else {

		var resp = responsesDTOs.AccessVerificationResponseDTO{StatusCode: 703, Result: nil, StatusDesc: "Error validating token"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// GetUserFromToken ...
// @Title Get user from token
// @Description Get user from token
// @Param	body		body 	requestsDTOs.StringRequestDTO	true		"body for Authentication content"
// @Success 200 {object} responsesDTOs.UserResponseDTO
// @Failure 403 body is empty
// @router /token/user [post]
func (c *AuthenticationController) GetUserFromToken() {
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "GetUserFromToken")
	var q requestsDTOs.TokenDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	statusCode := 608
	message := "Unable to verify token"

	authorization := c.Ctx.Input.Header("Authorization")
	requestUser, err := functions.GetUserFromBearerTokenContext(c.Ctx.Request.Context(), authorization)
	if err != nil {
		c.Ctx.Output.SetStatus(401)
		var resp = responsesDTOs.UserResponseDTO{StatusCode: 605, Result: nil, StatusDesc: "Invalid or missing access token"}
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	if tokenObj, err := models.GetRefreshTokensByTokenContext(c.Ctx.Request.Context(), q.Token); err == nil {

		if tokenObj.User == 0 {
			statusCode = 703
			message = "Token is not linked to any user"
			var resp = responsesDTOs.UserResponseDTO{StatusCode: statusCode, Result: nil, StatusDesc: message}
			c.Data["json"] = resp
			c.ServeJSON()
			return
		}

		// Optional hardening: bind to device/network
		if tokenObj.IPAddress != "" && tokenObj.IPAddress != c.Ctx.Request.RemoteAddr {
			c.Ctx.Output.SetStatus(401)
			c.Data["json"] = responsesDTOs.UserResponseDTO{
				StatusCode: 605,
				Result:     nil,
				StatusDesc: "Token IP mismatch",
			}
			c.ServeJSON()
			return
		}

		if tokenObj.Revoked || !tokenObj.ExpiresAt.After(time.Now().UTC()) {
			c.Ctx.Output.SetStatus(401)
			c.Data["json"] = responsesDTOs.UserResponseDTO{
				StatusCode: 605,
				Result:     nil,
				StatusDesc: "Refresh token expired or revoked",
			}
			c.ServeJSON()
			return
		}

		reqUserID, _ := strconv.ParseInt(requestUser.UserId, 10, 64)
		if tokenObj.User != reqUserID {
			c.Ctx.Output.SetStatus(401)
			c.Data["json"] = responsesDTOs.UserResponseDTO{
				StatusCode: 605,
				Result:     nil,
				StatusDesc: "Token does not belong to this user",
			}
			c.ServeJSON()
			return
		}

		statusCode = 200
		message = "Successfully validated"
		var resp = responsesDTOs.UserResponseDTO{StatusCode: statusCode, Result: requestUser, StatusDesc: message}
		c.Data["json"] = resp
	} else {

		statusCode = 703
		var resp = responsesDTOs.UserResponseDTO{StatusCode: statusCode, Result: nil, StatusDesc: message}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// VerifyCustomerTokenExpiry ...
// @Title Verify Customer token
// @Description Verify Customer token
// @Param	body		body 	requestsDTOs.StringRequestDTO	true		"body for Authentication content"
// @Success 200 {object} responsesDTOs.StringResponseDTO
// @Failure 403 body is empty
// @router /customer-token/verify [post]
func (c *AuthenticationController) VerifyCustomerToken() {
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "VerifyCustomerToken")
	var q requestsDTOs.StringRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	if token, err := functions.CheckCustomerTokenExpiry(q.Value); err == nil {
		if token.IsValid {

			var resp = responsesDTOs.CustomerAccessVerificationResponseDTO{StatusCode: 200, Result: token.Customer, StatusDesc: "Token is valid"}
			c.Data["json"] = resp
		} else {
			var resp = responsesDTOs.CustomerAccessVerificationResponseDTO{StatusCode: 605, Result: nil, StatusDesc: "Invalid token"}
			c.Data["json"] = resp
		}
	} else {

		var resp = responsesDTOs.CustomerAccessVerificationResponseDTO{StatusCode: 703, Result: nil, StatusDesc: "Error validating token"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// GetCustomerWithToken ...
// @Title Get Customer with token
// @Description Get Customer with token
// @Param	body		body 	requestsDTOs.StringRequestDTO	true		"body for Authentication content"
// @Success 200 {object} responsesDTOs.CustomerAccessVerificationResponseDTO
// @Failure 403 body is empty
// @router /customer-token/verify [post]
func (c *AuthenticationController) GetCustomerWithToken() {
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "GetCustomerWithToken")
	var q requestsDTOs.VerifyTokenReq
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	statusCode := 608
	message := "Unable to verify token"

	if tokenObj, err := models.GetUserTokensByTokenContext(c.Ctx.Request.Context(), q.Token); err == nil {

		if plainText, err := functions.GetAESDecrypted(tokenObj.Token, tokenObj.Nonce); err == nil {

			splitText := strings.Split(plainText, "__")

			email := ""
			if len(splitText) > 1 {
				email = splitText[0]
			} else {
				email = splitText[0]
			}

			req := requestsDTOs.GetUserWithUsernameRequest{Username: email}

			if user, err := functions.GetUserWithUsername(&c.Controller, req); err == nil {
				if user.StatusCode == 200 {

					statusCode = 200
					message = "Successfully validated"
				} else {
					statusCode = 708
					message = "user not found"
				}
				var resp = responsesDTOs.UserResponseDTO{StatusCode: statusCode, Result: user.Result, StatusDesc: message}
				c.Data["json"] = resp
			} else {
				statusCode = 708
				message = "user not found"
				var resp = responsesDTOs.UserResponseDTO{StatusCode: statusCode, Result: nil, StatusDesc: message}
				c.Data["json"] = resp
			}
		} else {

			statusCode = 708
			var resp = responsesDTOs.UserResponseDTO{StatusCode: statusCode, Result: nil, StatusDesc: message}
			c.Data["json"] = resp
		}
	} else {

		statusCode = 703
		var resp = responsesDTOs.UserResponseDTO{StatusCode: statusCode, Result: nil, StatusDesc: message}
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
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "GenerateInviteToken")
	var q requestsDTOs.EncryptInviteRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	rawString := q.Email + "___" + q.Role

	// ikey, _ := functions.GenerateKey()

	if token, nonce, err := functions.GetAESEncrypted(rawString); err == nil {

		var userToken models.UserTokens = models.UserTokens{Token: token, Nonce: nonce, ExpiryDate: time.Now().Add(4 * time.Hour), Active: 1, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: 1, ModifiedBy: 1}

		_, err := models.AddUserTokensContext(c.Ctx.Request.Context(), &userToken)

		if err != nil {
			return
		}

		var encryptResp responsesDTOs.InviteHashDTO = responsesDTOs.InviteHashDTO{Token: &userToken}
		var resp = responsesDTOs.InviteHashResponseDTO{StatusCode: 200, Value: &encryptResp, StatusDesc: "Encrypted"}
		c.Data["json"] = resp

	} else {

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
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "VerifyInviteToken")
	var q requestsDTOs.DecryptRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	// ikey, _ := functions.GenerateKey()

	if token, err := functions.GetAESDecrypted(q.Token, q.Nonce); err == nil {

		splitToken := strings.Split(string(token), "___")

		var token responsesDTOs.TokenDestructureResponseDTO = responsesDTOs.TokenDestructureResponseDTO{Email: splitToken[0], RoleId: splitToken[1]}

		// if splitToken[0] == q.Email {
		var resp = responsesDTOs.InviteDecodeResponseDTO{StatusCode: 200, Value: &token, StatusDesc: "Successfully verified token."}
		c.Data["json"] = resp
		// } else {
		// 	var resp = responsesDTOs.StringResponseDTO{StatusCode: 703, Value: "", StatusDesc: "Error validating token"}
		// 	c.Data["json"] = resp
		// }

		// if verified, err := functions.VerifyToken(splitToken[1]); err == nil {
		// 	if verified {
		// 		var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: splitToken[0], StatusDesc: "Successfully verified token."}
		// 		c.Data["json"] = resp
		// 	} else {
		// 		var resp = responsesDTOs.StringResponseDTO{StatusCode: 703, Value: "", StatusDesc: "Error validating token"}
		// 		c.Data["json"] = resp
		// 	}
		// } else {
		// 	var resp = responsesDTOs.StringResponseDTO{StatusCode: 703, Value: "", StatusDesc: "Error validating token"}
		// 	c.Data["json"] = resp
		// }
	} else {

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
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "Logout")
	var v requestsDTOs.TokenDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	if a, err := models.GetAccessTokensByTokenContext(c.Ctx.Request.Context(), v.Token); err == nil {
		a.ExpiresAt = time.Now()
		a.Revoked = true
		if err := models.UpdateAccessTokensByIdContext(c.Ctx.Request.Context(), a); err == nil {
			c.Ctx.Output.SetStatus(200)

			var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "", StatusDesc: "User log out complete"}
			c.Data["json"] = resp
		}
	} else {

		var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Invalid request"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// Log Customer Out ...
// @Title Expire Customer Token
// @Description Expire Customer Token
// @Param	body		body 	requestsDTOs.TokenDTO	true		"body for Authentication content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /expire-customer-token [post]
func (c *AuthenticationController) ExpireCustomerToken() {
	trace.Logger(c.Ctx.Request.Context()).Debug("authentication operation", "operation", "ExpireCustomerToken")
	var v requestsDTOs.TokenDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	if a, err := models.GetCustomer_access_tokensByTokenContext(c.Ctx.Request.Context(), v.Token); err == nil {
		a.ExpiresAt = time.Now()
		a.Revoked = true
		if err := models.UpdateCustomer_access_tokensByIdContext(c.Ctx.Request.Context(), a); err == nil {
			c.Ctx.Output.SetStatus(200)

			var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "", StatusDesc: "Customer expired"}
			c.Data["json"] = resp
		}
	} else {

		var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Invalid request"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}
