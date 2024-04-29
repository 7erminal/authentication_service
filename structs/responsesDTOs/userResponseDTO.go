package responsesDTOs

import "authentication_service/models"

type UserResponseDTO struct {
	StatusCode int           `orm: "omitempty"`
	User       *models.Users `orm: "omitempty"`
	StatusDesc string        `orm:"size(255)"`
}
