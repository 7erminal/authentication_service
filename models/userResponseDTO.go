package models

// import (
// 	"github.com/beego/beego/v2/client/orm"
// )

type UserResponseDTO struct {
	StatusCode int    `orm: "omitempty"`
	User       *Users `orm: "omitempty"`
	StatusDesc string `orm:"size(255)"`
}
