package models

import (
	"github.com/beego/beego/v2/client/orm"
)

type AuthenticationDTO struct {
	Username string `orm:"size(255)"`
	Password string `orm:"size(255)"`
}

func init() {
	orm.RegisterModel(new(AuthenticationDTO))
}
