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
	RequestId  string `json:"request_id,omitempty"`
}

type InviteHashDTO struct {
	Token *models.UserTokens
}

type InviteHashResponseDTO struct {
	StatusCode int
	Value      *InviteHashDTO
	StatusDesc string
	RequestId  string `json:"request_id,omitempty"`
}

type InviteDecodeResponseDTO struct {
	StatusCode int
	Value      *TokenDestructureResponseDTO
	StatusDesc string
	RequestId  string `json:"request_id,omitempty"`
}
