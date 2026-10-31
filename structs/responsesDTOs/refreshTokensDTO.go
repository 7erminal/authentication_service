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
	IssuedAt int64
	jwt.RegisteredClaims
}
