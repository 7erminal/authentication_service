package responsesDTOs

import (
	"github.com/golang-jwt/jwt/v5"
)

type TokenResponseDTO struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int64
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
	UserID      string           `json:"userid"`
	Username    string           `json:"username"`
	RoleID      string           `json:"roleid"`
	RoleName    string           `json:"rolename"`
	Permissions []UserPermission `json:"permissions"`
	BranchID    string           `json:"branchid"`
	Shop        string           `json:"shop"`
	jwt.RegisteredClaims
}

type CustomerAccessTokenClaims struct {
	CustomerId string `json:"customerid"`
	Username   string `json:"username"`
	Category   string `json:"category"`
	Number     string `json:"number"`
	ExpiryTime int64  `json:"expirytime"`
	BranchID   string `json:"branchid"`
	Shop       string `json:"shop"`
	jwt.RegisteredClaims
}

type AccessVerificationResponseDTO struct {
	StatusCode int
	Result     *AuthenticatedUser
	StatusDesc string
}

type CustomerAccessVerificationResponseDTO struct {
	StatusCode int
	Result     *AuthenticatedCustomer
	StatusDesc string
}
