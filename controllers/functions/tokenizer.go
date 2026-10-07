package functions

import (
	"authentication_service/structs/requestsDTOs"
	"authentication_service/structs/responsesDTOs"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("my32digitkey12345678901234567890")

// Generate a random AES key
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32) // 256-bit key
	_, err := rand.Read(key)
	return key, err
}

func CreateAccessToken(userid string, username string, roleid string, rolename string, permssions []responsesDTOs.UserPermission) (string, int64, string, error) {
	logs.Info("Creating access token for username: ", username, " and time now: ", time.Now())
	jti := uuid.NewString()
	expiryTime := time.Now().UTC().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userid":      userid,
			"username":    username,
			"roleid":      roleid,
			"rolename":    rolename,
			"exp":         expiryTime,
			"permissions": permssions,
			"jti":         jti,
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", 0, "", err
	}

	return tokenString, expiryTime, jti, nil
}

func CreateCustomerAccessToken(username string, category string) (string, int64, string, error) {
	logs.Info("Creating access token for username: ", username, " and time now: ", time.Now())
	jti := uuid.NewString()
	expiryTime := time.Now().UTC().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"customerId": username,
			"username":   username,
			"exp":        expiryTime,
			"category":   category,
			"number":     username,
			"jti":        jti,
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", 0, "", err
	}

	return tokenString, expiryTime, jti, nil
}

func CreateRefreshToken(username string) (string, int64, string, error) {
	jti := uuid.NewString()
	expiryTime := time.Now().UTC().Add(time.Hour * 24 * 7).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      expiryTime,
			"jti":      jti,
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", 0, "", err
	}
	return tokenString, expiryTime, jti, nil
	// return tokenString, expiryTime, nil
}

func CreateCustomerRefreshToken(username string) (string, int64, string, error) {
	jti := uuid.NewString()
	expiryTime := time.Now().UTC().Add(time.Hour * 24 * 7).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      expiryTime,
			"jti":      jti,
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", 0, "", err
	}
	return tokenString, expiryTime, jti, nil
	// return tokenString, expiryTime, nil
}

func VerifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, fmt.Errorf("invalid token")
	}

	return true, nil
}

// GetAccessTokenClaims decodes access token claims into a strongly typed struct.
func GetAccessTokenClaims(tokenString string) (*responsesDTOs.AccessTokenClaims, error) {
	claims := &responsesDTOs.AccessTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// GetCustomerAccessTokenClaims decodes customer token claims into a strongly typed struct.
func GetCustomerAccessTokenClaims(tokenString string) (*responsesDTOs.CustomerAccessTokenClaims, error) {
	claims := &responsesDTOs.CustomerAccessTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func CheckTokenExpiry(token_ string) (responsesDTOs.UserTokenResponseDTO, error) {

	if token, err := VerifyToken(token_); err == nil {
		if token {
			logs.Info("Valid token...")
			if claims, claimsErr := GetAccessTokenClaims(token_); claimsErr == nil {
				logs.Info("Decoded access token claims: username=", claims.Username, " roleid=", claims.RoleID, " rolename=", claims.RoleName, " permissions=", claims.Permissions)
				logs.Info("Token claims successfully decoded")

				authUser := &responsesDTOs.AuthenticatedUser{
					UserID:      claims.UserID,
					Username:    claims.Username,
					RoleID:      claims.RoleID,
					RoleName:    claims.RoleName,
					Permissions: claims.Permissions,
				}
				if claims.ExpiresAt.Unix() > time.Now().UTC().Unix() {
					logs.Info("Token is valid")
					resp := responsesDTOs.UserTokenResponseDTO{IsValid: true, User: authUser}
					return resp, nil
				} else {
					logs.Info("Token has expired")
					resp := responsesDTOs.UserTokenResponseDTO{IsValid: false, User: nil}
					return resp, nil
				}

			} else {
				logs.Error("Error decoding access token claims: ", claimsErr.Error())
				resp := responsesDTOs.UserTokenResponseDTO{IsValid: false, User: nil}
				return resp, claimsErr
			}
		} else {
			logs.Error("Token is invalid...")
			resp := responsesDTOs.UserTokenResponseDTO{IsValid: false, User: nil}
			return resp, nil
		}
	} else {
		logs.Error("Error validating token...", err.Error())
		resp := responsesDTOs.UserTokenResponseDTO{IsValid: false, User: nil}
		return resp, err
	}
}

func GetUserFromBearerToken(authorizationHeader string) (*responsesDTOs.Users, error) {
	authorizationHeader = strings.TrimSpace(authorizationHeader)
	if authorizationHeader == "" {
		return nil, fmt.Errorf("authorization header is required")
	}

	tokenString := authorizationHeader
	parts := strings.SplitN(authorizationHeader, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		tokenString = parts[1]
	}

	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, fmt.Errorf("access token is missing")
	}

	tokenResp, err := CheckTokenExpiry(tokenString)
	if err != nil {
		return nil, err
	}

	if !tokenResp.IsValid || tokenResp.User == nil {
		return nil, fmt.Errorf("invalid access token")
	}

	logs.Info("Getting user")
	userReq := requestsDTOs.GetUserRequest{UserId: tokenResp.User.UserID}
	if user, err := GetUser(&beego.Controller{}, userReq); err == nil {
		if user.StatusCode == 200 {
			return user.Result, nil
		} else {
			return nil, fmt.Errorf("error fetching user: %s", user.StatusDesc)
		}
	} else {
		return nil, err
	}
}

func CheckCustomerTokenExpiry(token_ string) (responsesDTOs.CustomerTokenResponseDTO, error) {

	if token, err := VerifyToken(token_); err == nil {
		if token {
			logs.Info("Valid token...")
			if claims, claimsErr := GetCustomerAccessTokenClaims(token_); claimsErr == nil {
				logs.Info("Decoded customer token claims: username=", claims.Username, " category=", claims.Category)
				logs.Info("Checking token expiry against current time: ", time.Now().UTC().Unix())
				if claims.ExpiresAt.Unix() > time.Now().UTC().Unix() {
					logs.Info("Token is valid")
					authCustomer := &responsesDTOs.AuthenticatedCustomer{
						CustomerId:       claims.CustomerId,
						Username:         claims.Username,
						Number:           claims.Number,
						CustomerCategory: claims.Category,
						ExpiryTime:       claims.ExpiresAt.Unix(),
					}
					resp := responsesDTOs.CustomerTokenResponseDTO{IsValid: true, Customer: authCustomer}
					return resp, nil
				} else {
					logs.Info("Token has expired")
					resp := responsesDTOs.CustomerTokenResponseDTO{IsValid: false, Customer: nil}
					return resp, nil
				}
			} else {
				logs.Error("Error decoding customer token claims: ", claimsErr.Error())
				resp := responsesDTOs.CustomerTokenResponseDTO{IsValid: false, Customer: nil}
				return resp, claimsErr
			}
		} else {
			logs.Error("Token is invalid...")
			resp := responsesDTOs.CustomerTokenResponseDTO{IsValid: false, Customer: nil}
			return resp, nil
		}
	} else {
		logs.Error("Error validating token...", err.Error())
		resp := responsesDTOs.CustomerTokenResponseDTO{IsValid: false, Customer: nil}
		return resp, err
	}
}

// var (
// 	// We're using a 32 byte long secret key.
// 	// This is probably something you generate first
// 	// then put into and environment variable.
// 	secretKey string = "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
// )

func GetAESDecrypted(encrypted string, nonce string) (string, error) {
	// key := "my32digitkey12345678901234567890"
	// iv := "my16digitIvKey12"

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	decodedCipherText, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	decodedNonce, err := base64.StdEncoding.DecodeString(nonce)
	if err != nil {
		return "", err
	}

	plainText, err := aesGCM.Open(nil, decodedNonce, decodedCipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

// PKCS5UnPadding  pads a certain blob of data with necessary data to be used in AES block cipher
func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}

// GetAESEncrypted encrypts given text in AES 256 CBC
func GetAESEncrypted(plaintext string) (string, string, error) {
	// key := "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
	// iv := "my16digitIvKey12"
	for {
		block, err := aes.NewCipher(secretKey)
		if err != nil {
			return "", "", err
		}

		aesGCM, err := cipher.NewGCM(block)
		if err != nil {
			return "", "", err
		}

		nonce := make([]byte, aesGCM.NonceSize())
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			return "", "", err
		}

		cipherText := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)

		encryptedString := base64.StdEncoding.EncodeToString(cipherText)

		if !strings.Contains(encryptedString, "/") {
			logs.Info("returning hash")
			return base64.StdEncoding.EncodeToString(cipherText), base64.StdEncoding.EncodeToString(nonce), nil
		}
	}
}

func EncryptInfo(plaintext string) (string, error) {
	for {
		// Generate a bcrypt hash
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plaintext), 8)
		if err != nil {
			return "", err
		}

		logs.Info("hashed password generated ", string(hashedPassword))

		// Check if hash contains a slash
		if !strings.Contains(string(hashedPassword), "/") {
			logs.Info("returning hash")
			return string(hashedPassword), nil
		}

		logs.Info("password contains slash. Regenerating hash")
	}
}
