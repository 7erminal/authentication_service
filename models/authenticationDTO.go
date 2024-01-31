package models

import (
	"github.com/beego/beego/v2/client/orm"
)

type AuthenticationDTO struct {
	Password string `orm:"size(255)"`
	Username string `orm:"size(255)"`
}

func init() {
	orm.RegisterModel(new(AuthenticationDTO))
}
