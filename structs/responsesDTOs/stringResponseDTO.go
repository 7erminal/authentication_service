package responsesDTOs

import "authentication_service/models"

type TokenDestructureResponseDTO struct {
	Email  string
	RoleId string
}

type StringResponseDTO struct {
	StatusCode int
	Value      string
	StatusDesc string
}

type InviteHashDTO struct {
	Token *models.UserTokens
}

type InviteHashResponseDTO struct {
	StatusCode int
	Value      *InviteHashDTO
	StatusDesc string
}

type InviteDecodeResponseDTO struct {
	StatusCode int
	Value      *TokenDestructureResponseDTO
	StatusDesc string
}
