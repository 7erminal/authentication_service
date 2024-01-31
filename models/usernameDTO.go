package models

import (
	"github.com/beego/beego/v2/client/orm"
)

type UsernameDTO struct {
	Username string `orm:"size(255)"`
}

func init() {
	orm.RegisterModel(new(UsernameDTO))
}
