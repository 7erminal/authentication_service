package models

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type AccessTokens struct {
	AccessTokenId int64     `orm:"auto;column(access_token_id)"`
	Token         string    `orm:"size(255)"`
	User          *Users    `orm:"rel(fk)"`
	DateCreated   time.Time `orm:"type(datetime)"`
	DateModified  time.Time `orm:"type(datetime)"`
	ExpiresAt     time.Time `orm:"type(datetime)"`
	Revoked       bool
	IPAddress     string    `orm:"size(80);column(ip_address)"`
	LastUsedAt    time.Time `orm:"type(datetime)"`
}

func init() {
	orm.RegisterModel(new(AccessTokens))
}

// AddAccessTokens insert a new AccessTokens into database and returns
// last inserted Id on success.
func AddAccessTokens(m *AccessTokens) (id int64, err error) {
	return AddAccessTokensCtx(context.Background(), m)
}

// AddAccessTokensCtx is AddAccessTokens with request context propagation.
func AddAccessTokensCtx(ctx context.Context, m *AccessTokens) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.InsertWithCtx(ctx, m)
	return
}

// GetAccessTokensById retrieves AccessTokens by Id. Returns error if
// Id doesn't exist
func GetAccessTokensById(id int64) (v *AccessTokens, err error) {
	o := orm.NewOrm()
	v = &AccessTokens{AccessTokenId: id}
	if err = o.QueryTable(new(AccessTokens)).Filter("AccessTokenId", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAccessTokensById retrieves AccessTokens by Id. Returns error if
// Id doesn't exist
func GetAccessTokensByToken(token string) (v *AccessTokens, err error) {
	return GetAccessTokensByTokenCtx(context.Background(), token)
}

// GetAccessTokensByTokenCtx is GetAccessTokensByToken with request context
// propagation.
func GetAccessTokensByTokenCtx(ctx context.Context, token string) (v *AccessTokens, err error) {
	o := orm.NewOrm()
	v = &AccessTokens{Token: token}
	if err = o.QueryTable(new(AccessTokens)).Filter("Token", token).RelatedSel().OneWithCtx(ctx, v); err == nil {
		_, err = o.LoadRelatedWithCtx(ctx, v.User, "UserDetails")
		if err != nil {
			slog.ErrorContext(ctx, "error loading related user details", "access_token_id", v.AccessTokenId, "error", err.Error())
		} else {
			_, err = o.LoadRelatedWithCtx(ctx, v.User.UserDetails, "Branch")
			if err != nil {
				slog.ErrorContext(ctx, "error loading related user details branch", "access_token_id", v.AccessTokenId, "error", err.Error())
			}
		}
		return v, nil
	}
	return nil, err
}

// GetAllAccessTokens retrieves all AccessTokens matches certain condition. Returns empty list if
// no records exist
func GetAllAccessTokens(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(AccessTokens))
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
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []AccessTokens
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateAccessTokens updates AccessTokens by Id and returns error if
// the record to be updated doesn't exist
func UpdateAccessTokensById(m *AccessTokens) (err error) {
	return UpdateAccessTokensByIdCtx(context.Background(), m)
}

// UpdateAccessTokensByIdCtx is UpdateAccessTokensById with request context
// propagation.
func UpdateAccessTokensByIdCtx(ctx context.Context, m *AccessTokens) (err error) {
	o := orm.NewOrm()
	v := AccessTokens{AccessTokenId: m.AccessTokenId}
	// ascertain id exists in the database
	if err = o.ReadWithCtx(ctx, &v); err == nil {
		var num int64
		if num, err = o.UpdateWithCtx(ctx, m); err == nil {
			slog.DebugContext(ctx, "records updated", "table", "access_tokens", "access_token_id", m.AccessTokenId, "rows", num)
		}
	}
	return
}

// UpdateAccessTokens updates AccessTokens by User Id and returns error if
// the record to be updated doesn't exist
func UpdateAccessTokensByUserId(m *AccessTokens) (err error) {
	return UpdateAccessTokensByUserIdCtx(context.Background(), m)
}

// UpdateAccessTokensByUserIdCtx is UpdateAccessTokensByUserId with request
// context propagation.
func UpdateAccessTokensByUserIdCtx(ctx context.Context, m *AccessTokens) (err error) {
	o := orm.NewOrm()
	v := AccessTokens{User: m.User}
	// ascertain id exists in the database
	if err = o.ReadWithCtx(ctx, &v, "User"); err == nil {
		var num int64
		if num, err = o.UpdateWithCtx(ctx, m, "Revoked"); err == nil {
			slog.DebugContext(ctx, "records updated", "table", "access_tokens", "access_token_id", m.AccessTokenId, "rows", num)
		}
	}
	return
}

// DeleteAccessTokens deletes AccessTokens by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAccessTokens(id int64) (err error) {
	o := orm.NewOrm()
	v := AccessTokens{AccessTokenId: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&AccessTokens{AccessTokenId: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
