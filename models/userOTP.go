package models

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"golang.org/x/crypto/bcrypt"
)

type UserOtps struct {
	UserId       int64     `orm:"size(255)"`
	OneTimePin   string    `orm:"size(255)"`
	Active       int       `orm: "omitempty"`
	DateCreated  time.Time `orm:"type(datetime)" orm: "omitempty"`
	DateModified time.Time `orm:"type(datetime)" orm: "omitempty"`
	CreatedBy    int       `orm: "omitempty"`
	ModifiedBy   int       `orm: "omitempty"`
}

func init() {
	orm.RegisterModel(new(UserOtps))
}

// Generate OTP
func GenerateUserOTP(id int64) (v *UserOtps, err error) {
	o := orm.NewOrm()
	generatedOTP := EncodeToString(6)

	logs.Info("OTP for user is ", generatedOTP)

	hashedPassword, errr := bcrypt.GenerateFromPassword([]byte(generatedOTP), 8)

	if errr == nil {
		v = &UserOtps{UserId: id, OneTimePin: string(hashedPassword), Active: 0, DateCreated: time.Now()}
		v = &UserOtps{UserId: id, OneTimePin: fmt.Sprint(111111), Active: 0, DateCreated: time.Now()}

		id, err = o.Insert(v)
		return nil, err
	} else {
		return nil, errr
	}
}

// VeiryUserOTP retrieves Users by Id. Returns error if
// Id doesn't exist
func VerifyUserOTP(id int64) (v *UserOtps, err error) {
	o := orm.NewOrm()
	v = &UserOtps{UserId: id}
	if err = o.QueryTable(new(UserOtps)).Filter("UserId", id).RelatedSel().OrderBy("-DateCreated").Limit(1).One(v); err == nil {
		return v, nil
	}
	return nil, err
}

func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)

	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
