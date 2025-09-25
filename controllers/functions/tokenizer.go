package functions

import (
	"authentication_service/models"
	"authentication_service/structs/responsesDTOs"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("my32digitkey12345678901234567890")

// Generate a random AES key
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32) // 256-bit key
	_, err := rand.Read(key)
	return key, err
}

func CreateToken(username string) (string, int64, error) {
	expiryTime := time.Now().Add(time.Hour * 4).Unix()
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

func CheckTokenExpiry(token_ string) (responsesDTOs.UserTokenResponseDTO, error) {

	if token, err := VerifyToken(token_); err == nil {
		if token {
			logs.Info("Valid token...")
			if tokenObj, err := models.GetAccessTokensByToken(token_); err == nil {
				logs.Info("Token fetched is ", tokenObj.Token)
				logs.Info("Token expiry is ", tokenObj.ExpiresAt)
				logs.Info("Time now is ", time.Now())

				if tokenObj.ExpiresAt.After(time.Now()) {
					logs.Info("Token is valid")
					resp := responsesDTOs.UserTokenResponseDTO{IsValid: true, User: tokenObj.User}
					return resp, nil
				} else {
					logs.Info("Token has expired")
					resp := responsesDTOs.UserTokenResponseDTO{IsValid: false, User: nil}
					return resp, nil
				}
			} else {
				logs.Error("Token does not exist...", err.Error())
				resp := responsesDTOs.UserTokenResponseDTO{IsValid: false, User: nil}
				return resp, err
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

func CheckCustomerTokenExpiry(token_ string) (responsesDTOs.CustomerTokenResponseDTO, error) {

	if token, err := VerifyToken(token_); err == nil {
		if token {
			logs.Info("Valid token...")
			if tokenObj, err := models.GetCustomer_access_tokensByToken(token_); err == nil {
				logs.Info("Token fetched is ", tokenObj.Token)
				logs.Info("Token expiry is ", tokenObj.ExpiresAt)
				logs.Info("Time now is ", time.Now())
				customerJson, err := json.Marshal(tokenObj.Customer)
				if err != nil {
					logs.Error("Error marshalling customer to JSON: ", err.Error())
				} else {
					logs.Info("Customer for token is ", string(customerJson))
				}

				if tokenObj.ExpiresAt.After(time.Now()) {
					logs.Info("Token is valid")
					resp := responsesDTOs.CustomerTokenResponseDTO{IsValid: true, Customer: tokenObj.Customer}
					return resp, nil
				} else {
					logs.Info("Token has expired")
					resp := responsesDTOs.CustomerTokenResponseDTO{IsValid: false, Customer: nil}
					return resp, nil
				}
			} else {
				logs.Error("Token does not exist...", err.Error())
				resp := responsesDTOs.CustomerTokenResponseDTO{IsValid: false, Customer: nil}
				return resp, err
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
