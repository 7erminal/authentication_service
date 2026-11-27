package controllers

import (
	"authentication_service/controllers/functions"
	"authentication_service/internal/logger"
	"authentication_service/internal/trace"
	"authentication_service/models"
	"authentication_service/structs/requestsDTOs"
	"authentication_service/structs/responsesDTOs"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"golang.org/x/crypto/bcrypt"
)

// AuthenticationController operations for Authentication
type AuthenticationController struct {
	BaseController
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
	c.Mapping("ValidateCustomerCredentialsToken", c.ValidateCustomerCredentialsToken)
	c.Mapping("ExpireCustomerToken", c.ExpireCustomerToken)
	c.Mapping("CheckCustomerTokenExpiry", c.CheckCustomerTokenExpiry)
	c.Mapping("ChangeCustomerPassword", c.ChangeCustomerPassword)
	c.Mapping("ResetCustomerPassword", c.ResetCustomerPassword)
	c.Mapping("RefreshAccessToken", c.RefreshAccessToken)
	c.Mapping("RefreshCustomerAccessToken", c.RefreshCustomerAccessToken)
	c.Mapping("PinValidation", c.PinValidation)
}

// Login ...
// @Title Login
// @Description Login User
// @Param	body		body 	models.AuthenticationDTO	true		"body for Authentication content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /login [post]
func (c *AuthenticationController) Login() {
	var v models.AuthenticationDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	logs.Info("Login request received")

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

// LoginToken ...
// @Title Login User
// @Description Login
// @Param	body		body 	models.AuthenticationDTO	true		"body for Authentication content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /login/token [post]
func (c *AuthenticationController) LoginToken() {
	ctx := c.reqCtx()
	var v models.AuthenticationDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	slog.InfoContext(ctx, "user login requested")

	ipAddress := c.Ctx.Request.RemoteAddr

	statusCode := 400
	statusMessage := "Failed"
	accessTokenObj := &models.AccessTokens{}
	refreshTokenObj := &models.RefreshTokens{}
	// userAgent := c.Ctx.Request.UserAgent()

	if a, err := models.GetUsersByUsernameCtx(ctx, strings.Trim(v.Username, " ")); err == nil {
		// Compare the stored hashed password, with the hashed version of the password that was received
		slog.InfoContext(ctx, "user credentials fetched", "user_id", a.UserId)
		if a.Active == 1 {
			if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(v.Password)); err != nil {
				// If the two passwords don't match, return a 401 status
				slog.InfoContext(ctx, "user login failed: incorrect password", "user_id", a.UserId)
				c.Data["json"] = err.Error()

				statusCode = 605
				statusMessage = "Incorrect password"
				var resp = responsesDTOs.UserResponseDTO{StatusCode: statusCode, User: nil, StatusDesc: statusMessage}
				c.Data["json"] = resp

			} else {
				c.Ctx.Output.SetStatus(200)

				slog.InfoContext(ctx, "user credentials validated", "user_id", a.UserId)

				// Create access token (15 minutes expiry)
				token, expiryTime, err := functions.CreateAccessToken(v.Username)

				if err != nil {
					slog.ErrorContext(ctx, "error generating access token", "user_id", a.UserId, logger.Err(err))
					statusCode = 301
					statusMessage = "Error generating token"
					var resp = responsesDTOs.StringResponseDTO{StatusCode: statusCode, Value: "", StatusDesc: statusMessage}
					c.Data["json"] = resp
				} else {
					// Revoke old tokens for this user
					updateToken := models.AccessTokens{User: a, Revoked: true}
					if err := models.UpdateAccessTokensByUserIdCtx(ctx, &updateToken); err != nil {
						slog.WarnContext(ctx, "error revoking old access tokens", "user_id", a.UserId, logger.Err(err))
					}

					t := time.Unix(expiryTime, 0)
					tokenObj := models.AccessTokens{
						User:         a,
						Token:        token,
						ExpiresAt:    t,
						DateCreated:  time.Now(),
						DateModified: time.Now(),
					}
					if _, err := models.AddAccessTokensCtx(ctx, &tokenObj); err == nil {
						statusCode = 200
						statusMessage = "Access token generated successfully"
						accessTokenObj = &tokenObj
						slog.InfoContext(ctx, "access token issued", "user_id", a.UserId, "access_token_id", tokenObj.AccessTokenId)

						// Create refresh token (7 days)
						refreshToken, refreshExpiryTime, err := functions.CreateRefreshToken(v.Username)
						if err != nil {
							slog.ErrorContext(ctx, "error generating refresh token", "user_id", a.UserId, logger.Err(err))
							c.Data["json"] = err.Error()
							c.ServeJSON()
							return
						}

						// Store refresh token
						refreshTokenObj = &models.RefreshTokens{
							User:         a,
							Token:        refreshToken,
							ExpiresAt:    time.Unix(refreshExpiryTime, 0),
							IPAddress:    ipAddress,
							UserAgent:    "", //userAgent,
							AccessToken:  accessTokenObj,
							DateCreated:  time.Now(),
							DateModified: time.Now(),
						}

						if _, err := models.AddRefreshTokensCtx(ctx, refreshTokenObj); err != nil {
							slog.ErrorContext(ctx, "error saving refresh token", "user_id", a.UserId, logger.Err(err))
							c.Data["json"] = err.Error()
							c.ServeJSON()
							return
						}
						slog.InfoContext(ctx, "refresh token issued", "user_id", a.UserId, "refresh_token_id", refreshTokenObj.RefreshTokenId)
					} else {
						slog.ErrorContext(ctx, "error saving access token", "user_id", a.UserId, logger.Err(err))
						statusCode = 301
						statusMessage = "Error generating token"
					}

				}
			}
		} else {
			slog.WarnContext(ctx, "user login rejected: inactive user", "user_id", a.UserId)
			statusCode = 607
			statusMessage = "Inactive user"
		}
	} else {
		slog.InfoContext(ctx, "user login failed: user not found", logger.Err(err))
		statusCode = 605
		statusMessage = "Unidentified user"
	}
	var tokenResponse = responsesDTOs.TokenResponseDTO{
		AccessToken:  accessTokenObj.Token,
		RefreshToken: refreshTokenObj.Token,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}

	userType := "USER"

	result := responsesDTOs.LoginDataResponseDTO{
		UserType: userType,
		Token:    &tokenResponse,
	}

	slog.InfoContext(ctx, "user login completed", "status_code", statusCode)

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
	ctx := c.reqCtx()
	refreshToken := c.Ctx.Input.Header("RefreshToken")

	slog.InfoContext(ctx, "user refresh token request received")

	// Validate refresh token
	if refreshTokenObj, err := models.GetRefreshTokensByTokenCtx(ctx, refreshToken); err == nil {
		slog.InfoContext(ctx, "refresh token found", "refresh_token_id", refreshTokenObj.RefreshTokenId)
		if refreshTokenObj.ExpiresAt.After(time.Now().UTC()) && !refreshTokenObj.Revoked {
			// Create new access token
			accessToken, accessExpiryTime, err := functions.CreateAccessToken(refreshTokenObj.User.Username)
			if err != nil {
				slog.ErrorContext(ctx, "error generating access token", "refresh_token_id", refreshTokenObj.RefreshTokenId, logger.Err(err))
				c.Data["json"] = err.Error()
				c.ServeJSON()
				return
			}

			accessTokenObj := models.AccessTokens{
				User:        refreshTokenObj.User,
				Token:       accessToken,
				ExpiresAt:   time.Unix(accessExpiryTime, 0),
				IPAddress:   c.Ctx.Request.RemoteAddr,
				DateCreated: time.Now(),
			}

			if _, err := models.AddAccessTokensCtx(ctx, &accessTokenObj); err != nil {
				slog.ErrorContext(ctx, "error saving access token", "refresh_token_id", refreshTokenObj.RefreshTokenId, logger.Err(err))
				c.Data["json"] = err.Error()
				c.ServeJSON()
				return
			}
			slog.InfoContext(ctx, "access token refreshed", "refresh_token_id", refreshTokenObj.RefreshTokenId, "access_token_id", accessTokenObj.AccessTokenId)

			var tokenResponse = responsesDTOs.TokenResponseDTO{
				AccessToken:  accessToken,
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
			slog.InfoContext(ctx, "refresh token expired or revoked", "refresh_token_id", refreshTokenObj.RefreshTokenId)
			c.Ctx.Output.SetStatus(401)
			var resp = responsesDTOs.StringResponseDTO{
				StatusCode: 605,
				Value:      "",
				StatusDesc: "Refresh token expired or revoked",
			}

			c.Data["json"] = resp
		}
	} else {
		slog.InfoContext(ctx, "invalid refresh token provided", logger.Err(err))
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
	ctx := c.reqCtx()
	var v models.AuthenticationDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	ipAddress := c.Ctx.Request.RemoteAddr

	clientCode := c.Ctx.Input.Header("clientCode")
	slog.InfoContext(ctx, "customer login requested", "client_code", clientCode)

	trimUsername := strings.TrimSpace(v.Username)

	statusCode := 400
	statusMessage := "Failed"
	accessTokenObj := &models.Customer_access_tokens{}
	refreshTokenObj := &models.CustomerRefreshTokens{}
	pinSet := false

	if clientCode != "" {
		if a, err := models.GetCustomer_credentialsByCustomerUsernameAndClientCodeCtx(ctx, trimUsername, clientCode); err == nil {
			// Compare the stored hashed password, with the hashed version of the password that was received
			slog.InfoContext(ctx, "customer credentials fetched", "customer_credential_id", a.Id)
			if a.Active == 1 {
				if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(v.Password)); err != nil {
					// If the two passwords don't match, return a 401 status
					slog.InfoContext(ctx, "customer login failed: incorrect password", "customer_credential_id", a.Id)
					c.Data["json"] = err.Error()

					statusCode = 605
					statusMessage = "Incorrect password"
				} else {
					c.Ctx.Output.SetStatus(200)

					if customer, err := models.GetCustomersByIdCtx(ctx, a.Customer.CustomerId); err != nil {
						slog.ErrorContext(ctx, "error fetching customer", "customer_id", a.Customer.CustomerId, logger.Err(err))
					} else {

						token, expiryTime, err := functions.CreateCustomerAccessToken(v.Username, *customer)

						if a.Pin != "1234" {
							pinSet = true
						}

						if err != nil {
							slog.ErrorContext(ctx, "error generating customer access token", "customer_id", customer.CustomerId, logger.Err(err))
							var resp = responsesDTOs.StringResponseDTO{StatusCode: 301, Value: "", StatusDesc: "Error generating token"}
							c.Data["json"] = resp
						} else {
							updateToken := models.Customer_access_tokens{Customer: a.Customer, Revoked: true}
							if err := models.UpdateCustomer_access_tokensByCustomerCtx(ctx, &updateToken); err != nil {
								statusMessage = fmt.Sprintf("Error revoking old tokens. %s", err.Error())
								slog.WarnContext(ctx, "error revoking old customer access tokens", "customer_id", customer.CustomerId, logger.Err(err))
							}
							t := time.Unix(expiryTime, 0)
							accessTokenObj = &models.Customer_access_tokens{Customer: a.Customer, Token: token, ExpiresAt: t, DateCreated: time.Now()}
							if _, err := models.AddCustomer_access_tokensCtx(ctx, accessTokenObj); err == nil {
								slog.InfoContext(ctx, "customer access token issued", "customer_id", customer.CustomerId, "customer_access_token_id", accessTokenObj.Id)
								statusCode = 200
								statusMessage = "Access token generated successfully"

								// Create refresh token (7 days)
								refreshToken, refreshExpiryTime, err := functions.CreateCustomerRefreshToken(v.Username)
								if err != nil {
									slog.ErrorContext(ctx, "error generating customer refresh token", "customer_id", customer.CustomerId, logger.Err(err))
									statusCode = 301
									statusMessage = "Error generating token"
								} else {
									// Store refresh token
									refreshTokenObj = &models.CustomerRefreshTokens{
										Customer:     a.Customer,
										Token:        refreshToken,
										ExpiresAt:    time.Unix(refreshExpiryTime, 0),
										IPAddress:    ipAddress,
										UserAgent:    clientCode,
										AccessToken:  accessTokenObj,
										DateCreated:  time.Now(),
										DateModified: time.Now(),
									}

									if _, err := models.AddCustomerRefreshTokensCtx(ctx, refreshTokenObj); err != nil {
										slog.ErrorContext(ctx, "error saving customer refresh token", "customer_id", customer.CustomerId, logger.Err(err))
										statusCode = 301
										statusMessage = "Error generating token"
									} else {
										slog.InfoContext(ctx, "customer refresh token issued", "customer_id", customer.CustomerId, "refresh_token_id", refreshTokenObj.RefreshTokenId)
										statusCode = 200
										statusMessage = "Tokens generated successfully"
									}
								}
							} else {
								slog.ErrorContext(ctx, "error saving customer access token", "customer_id", customer.CustomerId, logger.Err(err))
								statusCode = 301
								statusMessage = "Error generating token"
							}

						}
					}
				}
			} else {
				slog.WarnContext(ctx, "customer login rejected: inactive customer", "customer_credential_id", a.Id)
				var resp = responsesDTOs.StringResponseDTO{StatusCode: 607, Value: "", StatusDesc: "Inactive Customer"}
				c.Data["json"] = resp
			}
		} else {
			slog.InfoContext(ctx, "customer not found for client code, retrying with username only", "client_code", clientCode, logger.Err(err))
			if a, err := models.GetCustomer_credentialsByCustomerUsernameCtx(ctx, trimUsername); err == nil {
				slog.InfoContext(ctx, "customer found by username only", "customer_credential_id", a.Id)
				if a.ClientCode == "" {
					statusCode = 213
					statusMessage = "User Exists but not mapped to a client"
				}
			} else {
				slog.InfoContext(ctx, "customer login failed: customer not found", logger.Err(err))
				statusCode = 605
				statusMessage = "Unidentified customer"
			}
		}
	} else {
		slog.InfoContext(ctx, "customer login rejected: client code is empty")
		statusCode = 301
		statusMessage = "Client code is empty"
	}

	var tokenResponse = responsesDTOs.TokenResponseDTO{
		AccessToken:  accessTokenObj.Token,
		RefreshToken: refreshTokenObj.Token,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		PinSet:       pinSet,
	}

	userType := "CUSTOMER"

	result := responsesDTOs.LoginDataResponseDTO{
		UserType: userType,
		Token:    &tokenResponse,
	}
	slog.InfoContext(ctx, "customer login completed", "status_code", statusCode, "pin_set", pinSet)

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
	ctx := c.reqCtx()
	refreshToken := c.Ctx.Input.Header("RefreshToken")

	var v requestsDTOs.RefreshTokenRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	slog.InfoContext(ctx, "customer refresh token request received")

	// Validate refresh token
	if refreshTokenObj, err := models.GetCustomerRefreshTokensByTokenCtx(ctx, refreshToken); err == nil {
		if cc, err := models.GetCustomer_credentialsByCustomerUsernameAndClientCodeCtx(ctx, v.Username, refreshTokenObj.UserAgent); err == nil {

			if refreshTokenObj.ExpiresAt.After(time.Now().UTC()) && !refreshTokenObj.Revoked {
				// Create new access token
				accessToken, accessExpiryTime, err := functions.CreateCustomerAccessToken(cc.Username, *refreshTokenObj.Customer)
				if err != nil {
					slog.ErrorContext(ctx, "error generating customer access token", "refresh_token_id", refreshTokenObj.RefreshTokenId, logger.Err(err))
					c.Data["json"] = err.Error()
					c.ServeJSON()
					return
				}

				accessTokenObj := models.Customer_access_tokens{
					Customer:    refreshTokenObj.Customer,
					Token:       accessToken,
					ExpiresAt:   time.Unix(accessExpiryTime, 0),
					Revoked:     false,
					DateCreated: time.Now(),
				}

				if _, err := models.AddCustomer_access_tokensCtx(ctx, &accessTokenObj); err != nil {
					slog.ErrorContext(ctx, "error saving customer access token", "refresh_token_id", refreshTokenObj.RefreshTokenId, logger.Err(err))
					c.Data["json"] = err.Error()
					c.ServeJSON()
					return
				}
				slog.InfoContext(ctx, "customer access token refreshed", "refresh_token_id", refreshTokenObj.RefreshTokenId, "customer_access_token_id", accessTokenObj.Id)

				var data = responsesDTOs.TokenResponseDTO{
					AccessToken:  accessToken,
					RefreshToken: refreshToken, // Return same refresh token
					TokenType:    "Bearer",
					ExpiresIn:    900,
					PinSet:       false,
				}

				userType := "CUSTOMER"

				result := responsesDTOs.LoginDataResponseDTO{
					UserType: userType,
					Token:    &data,
				}

				var resp = responsesDTOs.LoginTokenResponseDTO{StatusCode: 200, StatusDesc: "Token generated", Result: &result}
				c.Data["json"] = resp
			} else {
				slog.InfoContext(ctx, "customer refresh token expired or revoked", "refresh_token_id", refreshTokenObj.RefreshTokenId)
				c.Ctx.Output.SetStatus(401)
				var resp = responsesDTOs.StringResponseDTO{
					StatusCode: 605,
					Value:      "",
					StatusDesc: "Refresh token expired or revoked",
				}
				c.Data["json"] = resp
			}
		} else {
			slog.InfoContext(ctx, "customer for refresh token not found", "refresh_token_id", refreshTokenObj.RefreshTokenId, logger.Err(err))
			c.Ctx.Output.SetStatus(200)
			var resp = responsesDTOs.StringResponseDTO{
				StatusCode: 608,
				Value:      "",
				StatusDesc: "User not found",
			}
			c.Data["json"] = resp
		}
	} else {
		slog.InfoContext(ctx, "invalid customer refresh token provided", logger.Err(err))
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
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)

	var v requestsDTOs.ChangePassword
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	logs.Info("Change password request received")

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

				a.Password = string(hashedPassword)

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

	logs.Info("Reset password request received")

	if a, err := models.GetUsersById(id); err == nil {
		// Compare the stored hashed password, with the hashed version of the password that was received

		hashedPassword, errr := bcrypt.GenerateFromPassword([]byte(v.NewPassword), 8)

		if errr == nil {

			a.Password = string(hashedPassword)

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

// Reset Customer Password ...
// @Title Reset Customer Password
// @Description Reset customer password
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	requestsDTOs.ResetPassword	true		"body for Change password content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /reset-customer-password/:id [put]
func (c *AuthenticationController) ResetCustomerPassword() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)

	var v requestsDTOs.ResetPassword
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	logs.Info("Reset customer password request received")

	if a, err := models.GetCustomersById(id); err == nil {
		// Compare the stored hashed password, with the hashed version of the password that was received

		if custCred, err := models.GetCustomer_credentialsByCustomerId(*a); err == nil {
			logs.Info("Customer credentials found")
			hashedPassword, errr := bcrypt.GenerateFromPassword([]byte(v.NewPassword), 8)

			if errr == nil {
				logs.Info("Successfully encrypted password")

				custCred.Password = string(hashedPassword)
				logs.Info("Password changed")

				// models.Agents{AgentName: v.AgentName, BranchId: v.BranchId, IdType: v.IdType, IdNumber: v.IdNumber, IsVerified: false, Active: 1, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: c_by, ModifiedBy: c_by}
			} else {
				logs.Error("Error hashing password ", errr.Error())
			}

			if err := models.UpdateCustomer_credentialsById(custCred); err == nil {
				logs.Info("Password change updated")
				c.Ctx.Output.SetStatus(200)

				var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "Successfully changed password", StatusDesc: "Customer password has been changed successfully"}
				c.Data["json"] = resp
			} else {
				var resp = responsesDTOs.StringResponseDTO{StatusCode: 608, Value: "", StatusDesc: "Customer password reset failed. " + err.Error()}
				c.Data["json"] = resp
			}
		} else {
			logs.Error(err.Error())
			var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Unidentified customer"}
			c.Data["json"] = resp
		}

	} else {
		logs.Error("Customer not found")
		logs.Error(err.Error())
		var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Unidentified Customer"}
		c.Data["json"] = resp
	}

	c.ServeJSON()
}

// Change Customer Password ...
// @Title Change Customer Password
// @Description Change customer password
// @Param	username	path 	string	true		"The username for which to change password"
// @Param	body		body 	requestsDTOs.ChangePassword	true		"body for Change password content"
// @Success 201 {object} models.UserResponseDTO
// @Failure 403 body is empty
// @router /change-customer-password/:username [put]
func (c *AuthenticationController) ChangeCustomerPassword() {
	ctx := c.reqCtx()
	username := c.Ctx.Input.Param(":username")

	// Try both ClientCode and clientCode header variants for compatibility
	clientCode := c.Ctx.Input.Header("ClientCode")
	if clientCode == "" {
		clientCode = c.Ctx.Input.Header("clientCode")
	}

	var v requestsDTOs.ChangePassword
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	slog.InfoContext(ctx, "customer password change requested", "client_code", clientCode)

	statusCode := 605
	statusMessage := "Internal server error"
	// Compare the stored hashed password, with the hashed version of the password that was received

	if custCred, err := models.GetCustomer_credentialsByCustomerUsernameAndClientCodeCtx(ctx, username, clientCode); err == nil {
		if err := bcrypt.CompareHashAndPassword([]byte(custCred.Password), []byte(v.OldPassword)); err != nil {
			// If the two passwords don't match, return a 401 status
			c.Data["json"] = err.Error()

			slog.InfoContext(ctx, "customer password change failed: old password does not match", "customer_credential_id", custCred.Id)
			statusCode = 601
			statusMessage = "Old password does not match"

		} else {
			hashedPassword, errr := bcrypt.GenerateFromPassword([]byte(v.NewPassword), 8)

			if errr == nil {
				custCred.Password = string(hashedPassword)
				custCred.DateModified = time.Now()

				// models.Agents{AgentName: v.AgentName, BranchId: v.BranchId, IdType: v.IdType, IdNumber: v.IdNumber, IsVerified: false, Active: 1, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: c_by, ModifiedBy: c_by}
			} else {
				slog.ErrorContext(ctx, "error hashing password", "customer_credential_id", custCred.Id, logger.Err(errr))
			}

			if err := models.UpdateCustomer_credentialsByIdCtx(ctx, custCred); err == nil {
				c.Ctx.Output.SetStatus(200)

				slog.InfoContext(ctx, "customer password changed", "customer_credential_id", custCred.Id)
				statusCode = 200
				statusMessage = "Customer password has been changed successfully"
			} else {
				slog.ErrorContext(ctx, "error updating customer credentials", "customer_credential_id", custCred.Id, logger.Err(err))
				statusCode = 608
				statusMessage = "Customer password reset failed. " + err.Error()
			}
		}
	} else {
		slog.InfoContext(ctx, "customer not found for client code, retrying with username only", "client_code", clientCode, logger.Err(err))
		statusCode = 604
		statusMessage = "Customer not found with username: " + username

		// 	query map[string]string, fields []string, sortby []string, order []string,
		// offset int64, limit int64
		query := map[string]string{"username": username}
		// TODO(tracing): GetAllCustomer_credentials has no ctx variant yet.
		if custCreds, err := models.GetAllCustomer_credentials(query, nil, nil, nil, 0, 1); err == nil {
			if len(custCreds) > 0 {
				for _, custCred := range custCreds {
					// Process each customer credential as needed
					custModel := custCred.(models.Customer_credentials)
					slog.InfoContext(ctx, "customer credentials found by username", "customer_credential_id", custModel.Id)
					if err := bcrypt.CompareHashAndPassword([]byte(custModel.Password), []byte(v.OldPassword)); err != nil {
						// If the two passwords don't match, return a 401 status
						c.Data["json"] = err.Error()
						slog.InfoContext(ctx, "customer password change failed: old password does not match", "customer_credential_id", custModel.Id)
						statusCode = 601
						statusMessage = "Old password does not match"

					} else {
						hashedPassword, errr := bcrypt.GenerateFromPassword([]byte(v.NewPassword), 8)
						if errr == nil {
							custModel.Password = string(hashedPassword)
							custModel.DateModified = time.Now()
							custModel.ClientCode = clientCode

							if err := models.UpdateCustomer_credentialsByIdCtx(ctx, &custModel); err != nil {
								slog.ErrorContext(ctx, "error updating customer credentials", "customer_credential_id", custModel.Id, logger.Err(err))
								statusCode = 608
								statusMessage = "Customer password reset failed. " + err.Error()
							} else {
								slog.InfoContext(ctx, "customer password changed", "customer_credential_id", custModel.Id)
								statusCode = 200
								statusMessage = "Customer password has been changed successfully"
								break
							}
						} else {
							slog.ErrorContext(ctx, "error hashing password", "customer_credential_id", custModel.Id, logger.Err(errr))
							statusCode = 609
							statusMessage = "Customer password reset failed. " + errr.Error()
						}
					}
				}
			} else {
				slog.InfoContext(ctx, "customer password change failed: no customer credentials found")
				statusCode = 604
				statusMessage = "Customer credentials not found for username: " + username
			}
		} else {
			slog.ErrorContext(ctx, "error retrieving customer credentials", logger.Err(err))
			statusCode = 607
			statusMessage = "Error retrieving customer credentials for username: " + username + ". " + err.Error()
		}
	}

	var resp = responsesDTOs.StringResponseDTO{StatusCode: statusCode, Value: "", StatusDesc: statusMessage}
	c.Data["json"] = resp
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
	ctx := c.reqCtx()
	var v requestsDTOs.ResetPasswordLink
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	slog.InfoContext(ctx, "password reset link requested", "link_count", len(v.Links))

	if a, err := models.GetUsersByUsernameCtx(ctx, v.Email); err == nil {
		// Compare the stored hashed password, with the hashed version of the password that was received

		// hashedPassword, errr := bcrypt.GenerateFromPassword([]byte(v.NewPassword), 8)

		rawString := v.Email + "___" + a.Role.Role

		// ikey, _ := functions.GenerateKey()

		if token, nonce, err := functions.GetAESEncrypted(rawString); err == nil {

			var userToken models.UserTokens = models.UserTokens{Token: token, Nonce: nonce, ExpiryDate: time.Now().Add(4 * time.Hour), Active: 1, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: 1, ModifiedBy: 1}

			_, err := models.AddUserTokens(&userToken)

			if err != nil {
				slog.ErrorContext(ctx, "error saving password reset token", "user_id", a.UserId, logger.Err(err))
				return
			}
			slog.InfoContext(ctx, "password reset token created", "user_id", a.UserId, "user_token_id", userToken.UserTokenId)

			namePlaceHolder := strings.Split(a.FullName, " | ")
			name := strings.Join(namePlaceHolder, " ")
			message_ := strings.Replace(v.Message, "[SENDER_NAME_ID]", name, -1)

			for i, link := range v.Links {
				iStr := strconv.Itoa(i)
				placeholder := "[LINK_" + iStr + "_ID]"
				formattedLink := *link + token
				message_ = strings.Replace(message_, placeholder, formattedLink, -1)
			}

			// Background job: detach from the request's cancellation but keep
			// the tracing IDs so the email worker logs correlate.
			go functions.SendEmailNewCtx(trace.Detach(ctx), a.Email, v.Subject, message_)

		} else {
			slog.ErrorContext(ctx, "error generating password reset token", "user_id", a.UserId, logger.Err(err))
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
		slog.InfoContext(ctx, "password reset link failed: user not found", logger.Err(err))
		var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Unidentified user"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// Pin Validation ...
// @Title Pin Validation
// @Description Pin validation using username
// @Param	body		body 	requestsDTOs.PinAuthenticationDTO	true		"body for Verify OTP content"
// @Success 201 {object} responsesDTOs.StringResponseDTO
// @Failure 403 body is empty
// @router /customer-pin-validation [post]
func (c *AuthenticationController) PinValidation() {
	ctx := c.reqCtx()
	var v requestsDTOs.PinAuthenticationDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	clientCode := c.Ctx.Input.Header("clientCode")
	slog.InfoContext(ctx, "customer pin validation requested", "client_code", clientCode)

	trimUsername := strings.TrimSpace(v.Username)

	statusCode := 400
	statusMessage := "Failed"
	value := ""

	if clientCode != "" {
		if a, err := models.GetCustomer_credentialsByCustomerUsernameAndClientCodeCtx(ctx, trimUsername, clientCode); err == nil {
			// Compare the stored hashed password, with the hashed version of the password that was received
			slog.InfoContext(ctx, "customer credentials fetched", "customer_credential_id", a.Id)
			if a.Active == 1 {
				if err := bcrypt.CompareHashAndPassword([]byte(a.Pin), []byte(v.Pin)); err != nil {
					// If the two pins don't match, return a 401 status
					slog.InfoContext(ctx, "customer pin validation failed: incorrect pin", "customer_credential_id", a.Id)
					c.Data["json"] = err.Error()

					statusCode = 605
					statusMessage = "Incorrect password"
				} else {
					c.Ctx.Output.SetStatus(200)
					slog.InfoContext(ctx, "customer pin validated", "customer_credential_id", a.Id)
					statusCode = 200
					value = "Pin validation successful"
					statusMessage = "Successfully validated pin"
				}
			} else {
				slog.WarnContext(ctx, "customer pin validation rejected: inactive customer", "customer_credential_id", a.Id)
				statusCode = 607
				statusMessage = "Inactive Customer"
			}
		} else {
			slog.InfoContext(ctx, "customer pin validation failed: customer not found", logger.Err(err))
			statusCode = 605
			statusMessage = "Unidentified customer"
		}
	} else {
		slog.InfoContext(ctx, "customer pin validation rejected: client code is empty")
		statusCode = 301
		statusMessage = "Client code is empty"
	}

	var resp = responsesDTOs.StringResponseDTO{StatusCode: statusCode, Value: value, StatusDesc: statusMessage}
	c.Data["json"] = resp
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

	if err != nil {
		var resp = responsesDTOs.UserResponseDTO{StatusCode: 604, User: nil, StatusDesc: "User cannot be found"}
		// c.Data["json"] = err.Error()
		c.Data["json"] = resp
	} else {
		// Get OTP
		logs.Debug("Got user. Now checking for user in OTP table ", v.UserId)
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
	proceed := false

	if ac, err := models.GetActivationCodesByNumber(q.MobileNumber); err == nil {
		// fmt.Printf("Value of v: %+v\n", ac)
		for _, suc := range ac {
			singleAc, err := suc.(models.ActivationCodes)
			logs.Debug("Activation...")
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

	logs.Debug("Verify activation code request received")

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

// CheckTokenExpiry ...
// @Title Check token expiry
// @Description Check Token Expiry
// @Param	body		body 	requestsDTOs.StringRequestDTO	true		"body for Authentication content"
// @Success 200 {object} responsesDTOs.StringResponseDTO
// @Failure 403 body is empty
// @router /token/check [post]
func (c *AuthenticationController) CheckTokenExpiry() {
	ctx := c.reqCtx()
	var q requestsDTOs.StringRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	slog.InfoContext(ctx, "user token check requested")

	if token, err := functions.CheckTokenExpiryCtx(ctx, q.Value); err == nil {
		if token.IsValid {
			slog.InfoContext(ctx, "user token is valid", "user_id", token.User.UserId)

			var resp = responsesDTOs.UserResponseDTO{StatusCode: 200, User: token.User, StatusDesc: "Token is valid"}
			c.Data["json"] = resp
		} else {
			slog.InfoContext(ctx, "user token is invalid or expired")
			var resp = responsesDTOs.UserResponseDTO{StatusCode: 605, User: nil, StatusDesc: "Invalid token"}
			c.Data["json"] = resp
		}
	} else {
		slog.InfoContext(ctx, "error validating user token", logger.Err(err))
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
	ctx := c.reqCtx()
	var q requestsDTOs.TokenDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	slog.InfoContext(ctx, "refresh token verification requested")

	statusCode := 608
	message := "Unable to verify token"

	authorization := c.Ctx.Input.Header("Authorization")
	requestUser, err := functions.GetUserFromBearerTokenCtx(ctx, authorization)
	if err != nil {
		slog.InfoContext(ctx, "refresh token verification rejected: invalid or missing access token", logger.Err(err))
		c.Ctx.Output.SetStatus(401)
		var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Invalid or missing access token"}
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	if tokenObj, err := models.GetRefreshTokensByTokenCtx(ctx, q.Token); err == nil {
		slog.InfoContext(ctx, "refresh token found", "refresh_token_id", tokenObj.RefreshTokenId, "user_id", requestUser.UserId)

		if tokenObj.User == nil {
			statusCode = 703
			message = "Token is not linked to any user"
			var resp = responsesDTOs.UserResponseDTO{StatusCode: statusCode, User: nil, StatusDesc: message}
			c.Data["json"] = resp
			c.ServeJSON()
			return
		}

		// Optional hardening: bind to device/network
		if tokenObj.IPAddress != "" && tokenObj.IPAddress != c.Ctx.Request.RemoteAddr {
			slog.InfoContext(ctx, "refresh token verification rejected: ip mismatch", "refresh_token_id", tokenObj.RefreshTokenId)
			c.Ctx.Output.SetStatus(401)
			c.Data["json"] = responsesDTOs.StringResponseDTO{
				StatusCode: 605,
				Value:      "",
				StatusDesc: "Token IP mismatch",
			}
			c.ServeJSON()
			return
		}

		if tokenObj.Revoked || !tokenObj.ExpiresAt.After(time.Now().UTC()) {
			slog.InfoContext(ctx, "refresh token verification rejected: expired or revoked", "refresh_token_id", tokenObj.RefreshTokenId)
			c.Ctx.Output.SetStatus(401)
			c.Data["json"] = responsesDTOs.StringResponseDTO{
				StatusCode: 605,
				Value:      "",
				StatusDesc: "Refresh token expired or revoked",
			}
			c.ServeJSON()
			return
		}

		if tokenObj.User.UserId != requestUser.UserId {
			slog.WarnContext(ctx, "refresh token verification rejected: token belongs to another user", "refresh_token_id", tokenObj.RefreshTokenId, "user_id", requestUser.UserId)
			c.Ctx.Output.SetStatus(401)
			c.Data["json"] = responsesDTOs.StringResponseDTO{
				StatusCode: 605,
				Value:      "",
				StatusDesc: "Token does not belong to this user",
			}
			c.ServeJSON()
			return
		}

		statusCode = 200
		message = "Successfully validated"
		var resp = responsesDTOs.UserResponseDTO{StatusCode: statusCode, User: tokenObj.User, StatusDesc: message}
		c.Data["json"] = resp
	} else {
		slog.InfoContext(ctx, "error validating refresh token", logger.Err(err))
		statusCode = 703
		var resp = responsesDTOs.UserResponseDTO{StatusCode: statusCode, User: nil, StatusDesc: message}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// CheckCustomerTokenExpiry ...
// @Title Check Customer token expiry
// @Description Check Token Expiry
// @Param	body		body 	requestsDTOs.StringRequestDTO	true		"body for Authentication content"
// @Success 200 {object} responsesDTOs.StringResponseDTO
// @Failure 403 body is empty
// @router /customer-token/check [post]
func (c *AuthenticationController) CheckCustomerTokenExpiry() {
	ctx := c.reqCtx()
	var q requestsDTOs.StringRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	slog.InfoContext(ctx, "customer token check requested")

	if token, err := functions.GetCustomerAccessTokenClaims(q.Value); err == nil {
		currentTime := time.Now().Unix()

		if token.Exp > currentTime {
			slog.InfoContext(ctx, "customer token is valid", "customer_id", token.CustomerID)

			var resp = responsesDTOs.CustomerAccessTokenResponseDTO{StatusCode: 200, Result: &responsesDTOs.CustomerAccessTokenData{
				CustomerID:     token.CustomerID,
				CustomerNumber: token.CustomerNumber,
				Exp:            token.Exp,
				Username:       token.Username,
				PhoneNumber:    token.PhoneNumber,
				Email:          token.Email,
				FullName:       token.FullName,
				ImagePath:      token.ImagePath,
				Location:       token.Location,
				Active:         token.Active,
			}, StatusDesc: "Token is valid"}
			c.Data["json"] = resp
		} else {
			slog.InfoContext(ctx, "customer token has expired", "customer_id", token.CustomerID, "expired_at", token.Exp)
			var resp = responsesDTOs.CustomerAccessTokenResponseDTO{StatusCode: 605, Result: nil, StatusDesc: "Token has expired"}
			c.Data["json"] = resp
		}

	} else {
		slog.InfoContext(ctx, "error validating customer token", logger.Err(err))
		var resp = responsesDTOs.CustomerAccessTokenResponseDTO{StatusCode: 703, Result: nil, StatusDesc: "Error validating token"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// VerifyCustomerToken ...
// @Title Verify Customer token
// @Description Verify Customer token
// @Param	body		body 	requestsDTOs.StringRequestDTO	true		"body for Authentication content"
// @Success 200 {object} responsesDTOs.StringResponseDTO
// @Failure 403 body is empty
// @router /customer-token/verify [post]
func (c *AuthenticationController) VerifyCustomerToken() {
	var q requestsDTOs.VerifyTokenReq
	json.Unmarshal(c.Ctx.Input.RequestBody, &q)

	logs.Info("Verify customer token request received")

	statusCode := 608
	message := "Unable to verify token"

	if tokenObj, err := models.GetUserTokensByToken(q.Token); err == nil {
		if plainText, err := functions.GetAESDecrypted(tokenObj.Token, tokenObj.Nonce); err == nil {

			splitText := strings.Split(plainText, "__")

			email := ""
			if len(splitText) > 1 {
				email = splitText[0]
			} else {
				email = splitText[0]
			}

			if user, err := models.GetUsersByUsername(email); err == nil {
				logs.Info("User was created on ", user.DateCreated)
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

	logs.Info("Generate invite token request received")
	rawString := q.Email + "___" + q.Role

	// ikey, _ := functions.GenerateKey()

	if token, nonce, err := functions.GetAESEncrypted(rawString); err == nil {

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

	logs.Info("Verify invite token request received")

	// ikey, _ := functions.GenerateKey()

	if token, err := functions.GetAESDecrypted(q.Token, q.Nonce); err == nil {
		splitToken := strings.Split(string(token), "___")

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
	ctx := c.reqCtx()
	var v requestsDTOs.TokenDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	slog.InfoContext(ctx, "user logout requested")

	if a, err := models.GetAccessTokensByTokenCtx(ctx, v.Token); err == nil {
		a.ExpiresAt = time.Now()
		a.Revoked = true
		if err := models.UpdateAccessTokensByIdCtx(ctx, a); err == nil {
			slog.InfoContext(ctx, "user logged out", "access_token_id", a.AccessTokenId)
			c.Ctx.Output.SetStatus(200)

			var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "", StatusDesc: "User log out complete"}
			c.Data["json"] = resp
		} else {
			slog.ErrorContext(ctx, "error revoking access token on logout", "access_token_id", a.AccessTokenId, logger.Err(err))
		}
	} else {
		slog.InfoContext(ctx, "user logout failed: access token not found", logger.Err(err))
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
	ctx := c.reqCtx()
	var v requestsDTOs.TokenDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	slog.InfoContext(ctx, "customer token expiry requested")

	if a, err := models.GetCustomer_access_tokensByTokenCtx(ctx, v.Token); err == nil {
		a.ExpiresAt = time.Now()
		a.Revoked = true
		if err := models.UpdateCustomer_access_tokensByIdCtx(ctx, a); err == nil {
			slog.InfoContext(ctx, "customer token expired", "customer_access_token_id", a.Id)
			c.Ctx.Output.SetStatus(200)

			var resp = responsesDTOs.StringResponseDTO{StatusCode: 200, Value: "", StatusDesc: "Customer expired"}
			c.Data["json"] = resp
		} else {
			slog.ErrorContext(ctx, "error expiring customer token", "customer_access_token_id", a.Id, logger.Err(err))
		}
	} else {
		slog.InfoContext(ctx, "customer token expiry failed: token not found", logger.Err(err))
		var resp = responsesDTOs.StringResponseDTO{StatusCode: 605, Value: "", StatusDesc: "Invalid request"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}
