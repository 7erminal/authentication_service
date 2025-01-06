package responsesDTOs

import "authentication_service/models"

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
