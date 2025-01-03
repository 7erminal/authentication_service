package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type PersonalAccessToken struct {
	Id             int64     `orm:"auto"`
	Tokenable_type string    `orm:"size(255)"`
	Tokenable_id   int       `orm:"size(255)"`
	Name           string    `orm:"size(255)"`
	Token          string    `orm:"size(255)"`
	Abilities      string    `orm:"size(255)"`
	LastUsedAt     time.Time `orm:"type(datetime)"`
	ExpiresAt      time.Time `orm:"type(datetime)"`
	CreatedAt      time.Time `orm:"type(datetime)"`
	updated_at     time.Time `orm:"type(datetime)"`
}

func init() {
	orm.RegisterModel(new(PersonalAccessToken))
}
