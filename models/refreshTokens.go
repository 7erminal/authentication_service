package models

import (
	"authentication_service/internal/trace"
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type RefreshTokens struct {
	RefreshTokenId int64         `orm:"auto;column(refresh_token_id)"`
	Token          string        `orm:"size(255);unique"`
	User           int64         `orm:"column(user_id)"`
	AccessToken    *AccessTokens `orm:"rel(fk);null"`
	DateCreated    time.Time     `orm:"type(datetime)"`
	DateModified   time.Time     `orm:"type(datetime)"`
	ExpiresAt      time.Time     `orm:"type(datetime)"`
	Revoked        bool          `orm:"default(false)"`
	IPAddress      string        `orm:"size(45);null"`
	UserAgent      string        `orm:"size(255);null"`
	LastUsedAt     time.Time     `orm:"type(datetime);null"`
}

func init() {
	orm.RegisterModel(new(RefreshTokens))
}

// AddRefreshTokens insert a new RefreshTokens into database and returns
// last inserted Id on success.
func AddRefreshTokensContext(ctx context.Context, m *RefreshTokens) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.InsertWithCtx(ctx, m)
	if err == nil {
		trace.Logger(ctx).Info("record created", "record_type", "refresh_tokens", "record_id", id)
	}
	return
}

// GetRefreshTokensById retrieves RefreshTokens by Id. Returns error if
// Id doesn't exist
func GetRefreshTokensByIdContext(ctx context.Context, id int64) (v *RefreshTokens, err error) {
	o := orm.NewOrm()
	v = &RefreshTokens{RefreshTokenId: id}
	if err = o.QueryTable(new(RefreshTokens)).Filter("RefreshTokenId", id).RelatedSel().OneWithCtx(ctx, v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetRefreshTokensByToken retrieves RefreshTokens by Token. Returns error if
// Token doesn't exist
func GetRefreshTokensByTokenContext(ctx context.Context, token string) (v *RefreshTokens, err error) {
	o := orm.NewOrm()
	v = &RefreshTokens{Token: token}
	if err = o.QueryTable(new(RefreshTokens)).Filter("Token", token).RelatedSel().OneWithCtx(ctx, v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetRefreshTokensByUser retrieves all RefreshTokens for a specific user
func GetRefreshTokensByUserContext(ctx context.Context, userId int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RefreshTokens)).Filter("User__UserId", userId).Filter("Revoked", false).RelatedSel()
	_, err = qs.AllWithCtx(ctx, &ml)
	if err == nil {
		return ml, nil
	}
	return nil, err
}

// GetAllRefreshTokens retrieves all RefreshTokens matches certain condition. Returns empty list if
// no records exist
func GetAllRefreshTokensContext(ctx context.Context, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RefreshTokens))
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
	_, err = qs.Limit(limit, offset).AllWithCtx(ctx, &ml, fields...)
	if err == nil {
		return ml, nil
	}
	return nil, err
}

// UpdateRefreshTokensById updates RefreshTokens by Id and returns error if
// the record to be updated doesn't exist
func UpdateRefreshTokensByIdContext(ctx context.Context, m *RefreshTokens) (err error) {
	o := orm.NewOrm()
	v := RefreshTokens{RefreshTokenId: m.RefreshTokenId}
	// ascertain id exists in the database
	if err = o.ReadWithCtx(ctx, &v); err == nil {
		var num int64
		if num, err = o.UpdateWithCtx(ctx, m); !errors.Is(err, nil) {
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
func DeleteRefreshTokensContext(ctx context.Context, id int64) (err error) {
	o := orm.NewOrm()
	v := RefreshTokens{RefreshTokenId: id}
	// ascertain id exists in the database
	if err = o.ReadWithCtx(ctx, &v); err == nil {
		var num int64
		if num, err = o.DeleteWithCtx(ctx, &v); !errors.Is(err, nil) {
			return
		} else if num != 1 {
			return errors.New("number of rows affected by delete: " + fmt.Sprint(num))
		}
		return
	}
	return
}

// RevokeRefreshTokensByUserId revokes all refresh tokens for a specific user
func RevokeRefreshTokensByUserIdContext(ctx context.Context, userId int64) (err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RefreshTokens))
	_, err = qs.Filter("User__UserId", userId).Filter("Revoked", false).UpdateWithCtx(ctx, orm.Params{
		"Revoked":      true,
		"DateModified": time.Now(),
	})
	return
}

// RevokeRefreshTokenByToken revokes a specific refresh token
func RevokeRefreshTokenByTokenContext(ctx context.Context, token string) (err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RefreshTokens))
	_, err = qs.Filter("Token", token).UpdateWithCtx(ctx, orm.Params{
		"Revoked":      true,
		"DateModified": time.Now(),
	})
	return
}

// UpdateRefreshTokenLastUsed updates the last used time for a refresh token
func UpdateRefreshTokenLastUsedContext(ctx context.Context, tokenId int64) (err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RefreshTokens))
	_, err = qs.Filter("RefreshTokenId", tokenId).UpdateWithCtx(ctx, orm.Params{
		"LastUsedAt":   time.Now(),
		"DateModified": time.Now(),
	})
	return
}

// ValidateRefreshToken checks if a refresh token is valid (not expired, not revoked)
func ValidateRefreshTokenContext(ctx context.Context, token string) (v *RefreshTokens, err error) {
	o := orm.NewOrm()
	v = &RefreshTokens{Token: token}
	if err = o.QueryTable(new(RefreshTokens)).Filter("Token", token).Filter("Revoked", false).RelatedSel().OneWithCtx(ctx, v); err == nil {
		if v.ExpiresAt.After(time.Now()) {
			// Update last used
			v.LastUsedAt = time.Now()
			o.UpdateWithCtx(ctx, v, "LastUsedAt", "DateModified")
			return v, nil
		}
		return nil, errors.New("refresh token expired")
	}
	return nil, err
}

// GetRefreshTokenDetails returns detailed info about a refresh token
func GetRefreshTokenDetailsContext(ctx context.Context, token string) (v *RefreshTokens, err error) {
	o := orm.NewOrm()
	v = &RefreshTokens{}
	if err = o.QueryTable(new(RefreshTokens)).Filter("Token", token).RelatedSel().OneWithCtx(ctx, v); err == nil {
		return v, nil
	}
	return nil, err
}

func (t *RefreshTokens) TableName() string {
	return "refresh_tokens"
}

// String() function is the same as the String function in AccessTokens
func (t *RefreshTokens) String() string {
	return reflect.Indirect(reflect.ValueOf(t)).Type().Name()
}

// AddRefreshTokens is a compatibility wrapper. Request handlers use AddRefreshTokensContext.
func AddRefreshTokens(m *RefreshTokens) (id int64, err error) {
	return AddRefreshTokensContext(context.Background(), m)
}

// GetRefreshTokensById is a compatibility wrapper. Request handlers use GetRefreshTokensByIdContext.
func GetRefreshTokensById(id int64) (v *RefreshTokens, err error) {
	return GetRefreshTokensByIdContext(context.Background(), id)
}

// GetRefreshTokensByToken is a compatibility wrapper. Request handlers use GetRefreshTokensByTokenContext.
func GetRefreshTokensByToken(token string) (v *RefreshTokens, err error) {
	return GetRefreshTokensByTokenContext(context.Background(), token)
}

// GetRefreshTokensByUser is a compatibility wrapper. Request handlers use GetRefreshTokensByUserContext.
func GetRefreshTokensByUser(userId int64) (ml []interface{}, err error) {
	return GetRefreshTokensByUserContext(context.Background(), userId)
}

// GetAllRefreshTokens is a compatibility wrapper. Request handlers use GetAllRefreshTokensContext.
func GetAllRefreshTokens(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	return GetAllRefreshTokensContext(context.Background(), query, fields, sortby, order, offset, limit)
}

// UpdateRefreshTokensById is a compatibility wrapper. Request handlers use UpdateRefreshTokensByIdContext.
func UpdateRefreshTokensById(m *RefreshTokens) (err error) {
	return UpdateRefreshTokensByIdContext(context.Background(), m)
}

// DeleteRefreshTokens is a compatibility wrapper. Request handlers use DeleteRefreshTokensContext.
func DeleteRefreshTokens(id int64) (err error) {
	return DeleteRefreshTokensContext(context.Background(), id)
}

// RevokeRefreshTokensByUserId is a compatibility wrapper. Request handlers use RevokeRefreshTokensByUserIdContext.
func RevokeRefreshTokensByUserId(userId int64) (err error) {
	return RevokeRefreshTokensByUserIdContext(context.Background(), userId)
}

// RevokeRefreshTokenByToken is a compatibility wrapper. Request handlers use RevokeRefreshTokenByTokenContext.
func RevokeRefreshTokenByToken(token string) (err error) {
	return RevokeRefreshTokenByTokenContext(context.Background(), token)
}

// UpdateRefreshTokenLastUsed is a compatibility wrapper. Request handlers use UpdateRefreshTokenLastUsedContext.
func UpdateRefreshTokenLastUsed(tokenId int64) (err error) {
	return UpdateRefreshTokenLastUsedContext(context.Background(), tokenId)
}

// ValidateRefreshToken is a compatibility wrapper. Request handlers use ValidateRefreshTokenContext.
func ValidateRefreshToken(token string) (v *RefreshTokens, err error) {
	return ValidateRefreshTokenContext(context.Background(), token)
}

// GetRefreshTokenDetails is a compatibility wrapper. Request handlers use GetRefreshTokenDetailsContext.
func GetRefreshTokenDetails(token string) (v *RefreshTokens, err error) {
	return GetRefreshTokenDetailsContext(context.Background(), token)
}
