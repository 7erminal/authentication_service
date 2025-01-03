package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type PasswordResetToken struct {
	Email     string    `orm:"size(255)"`
	Token     string    `orm:"size(255)"`
	CreatedAt time.Time `orm:"type(datetime)"`
}

func init() {
	orm.RegisterModel(new(PasswordResetToken))
}

func (u *PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

// AddToken insert a new User Token into database and returns
// last inserted Id on success.
func AddToken(m *PasswordResetToken) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}
