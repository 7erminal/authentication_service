package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type CustomerRefreshTokens struct {
	RefreshTokenId int64                   `orm:"auto;column(refresh_token_id)"`
	Token          string                  `orm:"size(255);unique"`
	Customer       *Customers              `orm:"rel(fk)"`
	AccessToken    *Customer_access_tokens `orm:"rel(fk);null"`
	DateCreated    time.Time               `orm:"type(datetime)"`
	DateModified   time.Time               `orm:"type(datetime)"`
	ExpiresAt      time.Time               `orm:"type(datetime)"`
	Revoked        bool                    `orm:"default(false)"`
	IPAddress      string                  `orm:"size(45);null"`
	UserAgent      string                  `orm:"size(255);null"`
	LastUsedAt     time.Time               `orm:"type(datetime);null"`
}

func init() {
	orm.RegisterModel(new(CustomerRefreshTokens))
}

// AddRefreshTokens insert a new RefreshTokens into database and returns
// last inserted Id on success.
func AddCustomerRefreshTokens(m *CustomerRefreshTokens) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetRefreshTokensById retrieves RefreshTokens by Id. Returns error if
// Id doesn't exist
func GetCustomerRefreshTokensById(id int64) (v *CustomerRefreshTokens, err error) {
	o := orm.NewOrm()
	v = &CustomerRefreshTokens{RefreshTokenId: id}
	if err = o.QueryTable(new(CustomerRefreshTokens)).Filter("RefreshTokenId", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetRefreshTokensByToken retrieves RefreshTokens by Token. Returns error if
// Token doesn't exist
func GetCustomerRefreshTokensByToken(token string) (v *CustomerRefreshTokens, err error) {
	o := orm.NewOrm()
	v = &CustomerRefreshTokens{Token: token}
	if err = o.QueryTable(new(CustomerRefreshTokens)).Filter("Token", token).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetRefreshTokensByUser retrieves all RefreshTokens for a specific user
func GetCustomerRefreshTokensByUser(userId int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(CustomerRefreshTokens)).Filter("Customer__CustomerId", userId).Filter("Revoked", false).RelatedSel()
	_, err = qs.All(&ml)
	if err == nil {
		return ml, nil
	}
	return nil, err
}

// GetAllRefreshTokens retrieves all RefreshTokens matches certain condition. Returns empty list if
// no records exist
func GetAllCustomerRefreshTokens(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(CustomerRefreshTokens))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		}
	} else {
		if len(fields) != 0 {
			// default order by primary key
			qs = qs.OrderBy("-RefreshTokenId")
		}
	}
	qs = qs.RelatedSel()
	if limit == 0 {
		limit = 1000
	}
	_, err = qs.Limit(limit, offset).All(&ml, fields...)
	if err == nil {
		return ml, nil
	}
	return nil, err
}

// UpdateRefreshTokensById updates RefreshTokens by Id and returns error if
// the record to be updated doesn't exist
func UpdateCustomerRefreshTokensById(m *CustomerRefreshTokens) (err error) {
	o := orm.NewOrm()
	v := CustomerRefreshTokens{RefreshTokenId: m.RefreshTokenId}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); !errors.Is(err, nil) {
			return
		} else if num != 1 {
			return errors.New("number of rows affected by update: " + fmt.Sprint(num))
		}
		return
	}
	return
}

// DeleteRefreshTokens deletes RefreshTokens by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCustomerRefreshTokens(id int64) (err error) {
	o := orm.NewOrm()
	v := CustomerRefreshTokens{RefreshTokenId: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&v); !errors.Is(err, nil) {
			return
		} else if num != 1 {
			return errors.New("number of rows affected by delete: " + fmt.Sprint(num))
		}
		return
	}
	return
}

// RevokeRefreshTokensByUserId revokes all refresh tokens for a specific user
func RevokeCustomerRefreshTokensByUserId(userId int64) (err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(CustomerRefreshTokens))
	_, err = qs.Filter("User__UserId", userId).Filter("Revoked", false).Update(orm.Params{
		"Revoked":      true,
		"DateModified": time.Now(),
	})
	return
}

// RevokeRefreshTokenByToken revokes a specific refresh token
func RevokeCustomerRefreshTokenByToken(token string) (err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(CustomerRefreshTokens))
	_, err = qs.Filter("Token", token).Update(orm.Params{
		"Revoked":      true,
		"DateModified": time.Now(),
	})
	return
}

// UpdateRefreshTokenLastUsed updates the last used time for a refresh token
func UpdateCustomerRefreshTokenLastUsed(tokenId int64) (err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(CustomerRefreshTokens))
	_, err = qs.Filter("RefreshTokenId", tokenId).Update(orm.Params{
		"LastUsedAt":   time.Now(),
		"DateModified": time.Now(),
	})
	return
}

// ValidateRefreshToken checks if a refresh token is valid (not expired, not revoked)
func ValidateCustomerRefreshToken(token string) (v *CustomerRefreshTokens, err error) {
	o := orm.NewOrm()
	v = &CustomerRefreshTokens{Token: token}
	if err = o.QueryTable(new(CustomerRefreshTokens)).Filter("Token", token).Filter("Revoked", false).RelatedSel().One(v); err == nil {
		if v.ExpiresAt.After(time.Now()) {
			// Update last used
			v.LastUsedAt = time.Now()
			o.Update(v, "LastUsedAt", "DateModified")
			return v, nil
		}
		return nil, errors.New("refresh token expired")
	}
	return nil, err
}

// GetRefreshTokenDetails returns detailed info about a refresh token
func GetCustomerRefreshTokenDetails(token string) (v *CustomerRefreshTokens, err error) {
	o := orm.NewOrm()
	v = &CustomerRefreshTokens{}
	if err = o.QueryTable(new(CustomerRefreshTokens)).Filter("Token", token).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

func (t *CustomerRefreshTokens) TableName() string {
	return "customer_refresh_tokens"
}

// String() function is the same as the String function in AccessTokens
func (t *CustomerRefreshTokens) String() string {
	return reflect.Indirect(reflect.ValueOf(t)).Type().Name()
}
