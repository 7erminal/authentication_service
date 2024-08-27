package responsesDTOs

import "authentication_service/models"

type UserResponseDTO struct {
	StatusCode int
	User       *models.Users
	StatusDesc string
}
