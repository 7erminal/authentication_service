package responsesDTOs

import "github.com/golang-jwt/jwt/v5"

type TokenResponseDTO struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int64
	PinSet       bool
}

type LoginDataResponseDTO struct {
	UserType string
	Token    *TokenResponseDTO
}

type LoginTokenResponseDTO struct {
	StatusCode int
	Result     *LoginDataResponseDTO
	StatusDesc string
}

type AccessTokenClaims struct {
	Username string
	Exp      int64
	jwt.RegisteredClaims
}

type CustomerAccessTokenClaims struct {
	CustomerID     string `json:"customerId"`
	Exp            int64  `json:"exp"`
	Username       string `json:"username"`
	CustomerNumber string `json:"customerNumber"`
	PhoneNumber    string `json:"phoneNumber"`
	Email          string `json:"email"`
	FullName       string `json:"fullName"`
	ImagePath      string `json:"imagePath"`
	Location       string `json:"location"`
	Active         int    `json:"active"`
	jwt.RegisteredClaims
}
