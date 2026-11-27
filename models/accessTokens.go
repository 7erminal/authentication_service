package models

import (
	"authentication_service/internal/trace"
	"context"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type AccessTokens struct {
	AccessTokenId int64     `orm:"auto;column(access_token_id)"`
	Token         string    `orm:"size(255)"`
	User          int64     `orm:"column(user_id)"`
	DateCreated   time.Time `orm:"type(datetime)"`
	DateModified  time.Time `orm:"type(datetime)"`
	ExpiresAt     time.Time `orm:"type(datetime)"`
	Revoked       bool
	IPAddress     string    `orm:"size(80);column(ip_address)"`
	LastUsedAt    time.Time `orm:"type(datetime);null"`
}

func init() {
	orm.RegisterModel(new(AccessTokens))
}

// AddAccessTokens insert a new AccessTokens into database and returns
// last inserted Id on success.
func AddAccessTokensContext(ctx context.Context, m *AccessTokens) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.InsertWithCtx(ctx, m)
	if err == nil {
		trace.Logger(ctx).Info("record created", "record_type", "access_tokens", "record_id", id)
	}
	return
}

// GetAccessTokensById retrieves AccessTokens by Id. Returns error if
// Id doesn't exist
func GetAccessTokensByIdContext(ctx context.Context, id int64) (v *AccessTokens, err error) {
	o := orm.NewOrm()
	v = &AccessTokens{AccessTokenId: id}
	if err = o.QueryTable(new(AccessTokens)).Filter("AccessTokenId", id).RelatedSel().OneWithCtx(ctx, v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAccessTokensById retrieves AccessTokens by Id. Returns error if
// Id doesn't exist
func GetAccessTokensByTokenContext(ctx context.Context, token string) (v *AccessTokens, err error) {
	o := orm.NewOrm()
	v = &AccessTokens{Token: token}
	if err = o.QueryTable(new(AccessTokens)).Filter("Token", token).RelatedSel().OneWithCtx(ctx, v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAccessTokens retrieves all AccessTokens matches certain condition. Returns empty list if
// no records exist
func GetAllAccessTokensContext(ctx context.Context, query map[string]string, fields []string, sortby []string, order []string,
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
	if _, err = qs.Limit(limit, offset).AllWithCtx(ctx, &l, fields...); err == nil {
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
func UpdateAccessTokensByIdContext(ctx context.Context, m *AccessTokens) (err error) {
	o := orm.NewOrm()
	v := AccessTokens{AccessTokenId: m.AccessTokenId}
	// ascertain id exists in the database
	if err = o.ReadWithCtx(ctx, &v); err == nil {
		var num int64
		if num, err = o.UpdateWithCtx(ctx, m); err == nil {
			trace.Logger(ctx).Debug("database rows changed", "rows", num)
		}
	}
	return
}

// UpdateAccessTokens updates AccessTokens by User Id and returns error if
// the record to be updated doesn't exist
func UpdateAccessTokensByUserIdContext(ctx context.Context, m *AccessTokens) (err error) {
	o := orm.NewOrm()
	v := AccessTokens{User: m.User}
	// ascertain id exists in the database
	if err = o.ReadWithCtx(ctx, &v, "User"); err == nil {
		var num int64
		if num, err = o.UpdateWithCtx(ctx, m, "Revoked"); err == nil {
			trace.Logger(ctx).Debug("database rows changed", "rows", num)
		}
	}
	return
}

// DeleteAccessTokens deletes AccessTokens by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAccessTokensContext(ctx context.Context, id int64) (err error) {
	o := orm.NewOrm()
	v := AccessTokens{AccessTokenId: id}
	// ascertain id exists in the database
	if err = o.ReadWithCtx(ctx, &v); err == nil {
		var num int64
		if num, err = o.DeleteWithCtx(ctx, &AccessTokens{AccessTokenId: id}); err == nil {
			trace.Logger(ctx).Debug("database rows changed", "rows", num)
		}
	}
	return
}

// AddAccessTokens is a compatibility wrapper. Request handlers use AddAccessTokensContext.
func AddAccessTokens(m *AccessTokens) (id int64, err error) {
	return AddAccessTokensContext(context.Background(), m)
}

// GetAccessTokensById is a compatibility wrapper. Request handlers use GetAccessTokensByIdContext.
func GetAccessTokensById(id int64) (v *AccessTokens, err error) {
	return GetAccessTokensByIdContext(context.Background(), id)
}

// GetAccessTokensByToken is a compatibility wrapper. Request handlers use GetAccessTokensByTokenContext.
func GetAccessTokensByToken(token string) (v *AccessTokens, err error) {
	return GetAccessTokensByTokenContext(context.Background(), token)
}

// GetAllAccessTokens is a compatibility wrapper. Request handlers use GetAllAccessTokensContext.
func GetAllAccessTokens(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	return GetAllAccessTokensContext(context.Background(), query, fields, sortby, order, offset, limit)
}

// UpdateAccessTokensById is a compatibility wrapper. Request handlers use UpdateAccessTokensByIdContext.
func UpdateAccessTokensById(m *AccessTokens) (err error) {
	return UpdateAccessTokensByIdContext(context.Background(), m)
}

// UpdateAccessTokensByUserId is a compatibility wrapper. Request handlers use UpdateAccessTokensByUserIdContext.
func UpdateAccessTokensByUserId(m *AccessTokens) (err error) {
	return UpdateAccessTokensByUserIdContext(context.Background(), m)
}

// DeleteAccessTokens is a compatibility wrapper. Request handlers use DeleteAccessTokensContext.
func DeleteAccessTokens(id int64) (err error) {
	return DeleteAccessTokensContext(context.Background(), id)
}
