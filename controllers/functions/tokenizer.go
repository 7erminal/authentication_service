package functions

import (
	"authentication_service/models"
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

func CreateToken(username string) (string, int64, error) {
	expiryTime := time.Now().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      expiryTime,
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiryTime, nil
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

func CheckTokenExpiry(token_ string) (bool, error) {

	if token, err := VerifyToken(token_); err == nil {
		if token {
			logs.Info("Valid token...")
			if tokenObj, err := models.GetAccessTokensByToken(token_); err == nil {
				logs.Info("Token fetched is ", tokenObj.Token)
				logs.Info("Token expiry is ", tokenObj.ExpiresAt)
				logs.Info("Time now is ", time.Now())

				if tokenObj.ExpiresAt.After(time.Now()) {
					logs.Info("Token is valid")
					return true, nil
				} else {
					logs.Info("Token has expired")
					return false, nil
				}
			} else {
				logs.Error("Token does not exist...", err.Error())
				return false, err
			}
		} else {
			logs.Error("Token is invalid...")
			return false, nil
		}
	} else {
		logs.Error("Error validating token...", err.Error())
		return false, err
	}
}