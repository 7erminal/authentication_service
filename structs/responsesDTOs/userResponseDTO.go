package responsesDTOs

import (
	"authentication_service/models"
	"time"
)

type UserResp struct {
	UserId        int64
	ImagePath     string
	UserType      int
	FullName      string
	Username      string
	Password      string
	Email         string
	PhoneNumber   string
	Gender        string
	Dob           time.Time
	Address       string
	IdType        string
	IdNumber      string
	MaritalStatus string
	Active        int
	Role          *models.Roles
	IsVerified    bool
	DateCreated   time.Time
	DateModified  time.Time
	CreatedBy     int
	ModifiedBy    int
	Branch        *models.Branches
}

type UserResponseDTO struct {
	StatusCode int
	User       *models.Users
	StatusDesc string
}

type UserTokenResponseDTO struct {
	IsValid bool
	User    *models.Users
}

type CustomerTokenResponseDTO struct {
	IsValid  bool
	Customer *models.Customers
}

type CustomerResponseDTO struct {
	StatusCode int
	Result     *models.Customers
	StatusDesc string
}
