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

type Customer_access_tokens struct {
	Id           int64     `orm:"auto;column(customer_access_token_id)"`
	Token        string    `orm:"size(255)"`
	Customer     int64     `orm:"column(customer_id)"`
	DateCreated  time.Time `orm:"type(datetime)"`
	DateModified time.Time `orm:"type(datetime)"`
	ExpiresAt    time.Time `orm:"type(datetime)"`
	Revoked      bool
	IpAddress    string    `orm:"size(80);column(ip_address)"`
	LastUsedAt   time.Time `orm:"type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Customer_access_tokens))
}

// AddCustomer_access_tokens insert a new Customer_access_tokens into database and returns
// last inserted Id on success.
func AddCustomer_access_tokensContext(ctx context.Context, m *Customer_access_tokens) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.InsertWithCtx(ctx, m)
	if err == nil {
		trace.Logger(ctx).Info("record created", "record_type", "customer_access_tokens", "record_id", id)
	}
	return
}

// GetCustomer_access_tokensById retrieves Customer_access_tokens by Id. Returns error if
// Id doesn't exist
func GetCustomer_access_tokensByIdContext(ctx context.Context, id int64) (v *Customer_access_tokens, err error) {
	o := orm.NewOrm()
	v = &Customer_access_tokens{Id: id}
	if err = o.QueryTable(new(Customer_access_tokens)).Filter("Id", id).RelatedSel().OneWithCtx(ctx, v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetCustomer_access_tokensById retrieves Customer_access_tokens by Id. Returns error if
// Id doesn't exist
func GetCustomer_access_tokensByTokenContext(ctx context.Context, token string) (v *Customer_access_tokens, err error) {
	o := orm.NewOrm()
	v = &Customer_access_tokens{Token: token}
	if err = o.QueryTable(new(Customer_access_tokens)).Filter("Token", token).RelatedSel().OneWithCtx(ctx, v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCustomer_access_tokens retrieves all Customer_access_tokens matches certain condition. Returns empty list if
// no records exist
func GetAllCustomer_access_tokensContext(ctx context.Context, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Customer_access_tokens))
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

	var l []Customer_access_tokens
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

// UpdateCustomer_access_tokens updates Customer_access_tokens by Id and returns error if
// the record to be updated doesn't exist
func UpdateCustomer_access_tokensByIdContext(ctx context.Context, m *Customer_access_tokens) (err error) {
	o := orm.NewOrm()
	v := Customer_access_tokens{Id: m.Id}
	// ascertain id exists in the database
	if err = o.ReadWithCtx(ctx, &v); err == nil {
		var num int64
		if num, err = o.UpdateWithCtx(ctx, m); err == nil {
			trace.Logger(ctx).Debug("database rows changed", "rows", num)
		}
	}
	return
}

// UpdateCustomer_access_tokens updates Customer_access_tokens by Id and returns error if
// the record to be updated doesn't exist
func UpdateCustomer_access_tokensByCustomerContext(ctx context.Context, m *Customer_access_tokens) (err error) {
	o := orm.NewOrm()
	v := Customer_access_tokens{Customer: m.Customer}
	// ascertain id exists in the database
	if err = o.ReadWithCtx(ctx, &v, "Customer"); err == nil {
		var num int64
		if num, err = o.UpdateWithCtx(ctx, m); err == nil {
			trace.Logger(ctx).Debug("database rows changed", "rows", num)
		}
	}
	return
}

// DeleteCustomer_access_tokens deletes Customer_access_tokens by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCustomer_access_tokensContext(ctx context.Context, id int64) (err error) {
	o := orm.NewOrm()
	v := Customer_access_tokens{Id: id}
	// ascertain id exists in the database
	if err = o.ReadWithCtx(ctx, &v); err == nil {
		var num int64
		if num, err = o.DeleteWithCtx(ctx, &Customer_access_tokens{Id: id}); err == nil {
			trace.Logger(ctx).Debug("database rows changed", "rows", num)
		}
	}
	return
}

// AddCustomer_access_tokens is a compatibility wrapper. Request handlers use AddCustomer_access_tokensContext.
func AddCustomer_access_tokens(m *Customer_access_tokens) (id int64, err error) {
	return AddCustomer_access_tokensContext(context.Background(), m)
}

// GetCustomer_access_tokensById is a compatibility wrapper. Request handlers use GetCustomer_access_tokensByIdContext.
func GetCustomer_access_tokensById(id int64) (v *Customer_access_tokens, err error) {
	return GetCustomer_access_tokensByIdContext(context.Background(), id)
}

// GetCustomer_access_tokensByToken is a compatibility wrapper. Request handlers use GetCustomer_access_tokensByTokenContext.
func GetCustomer_access_tokensByToken(token string) (v *Customer_access_tokens, err error) {
	return GetCustomer_access_tokensByTokenContext(context.Background(), token)
}

// GetAllCustomer_access_tokens is a compatibility wrapper. Request handlers use GetAllCustomer_access_tokensContext.
func GetAllCustomer_access_tokens(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	return GetAllCustomer_access_tokensContext(context.Background(), query, fields, sortby, order, offset, limit)
}

// UpdateCustomer_access_tokensById is a compatibility wrapper. Request handlers use UpdateCustomer_access_tokensByIdContext.
func UpdateCustomer_access_tokensById(m *Customer_access_tokens) (err error) {
	return UpdateCustomer_access_tokensByIdContext(context.Background(), m)
}

// UpdateCustomer_access_tokensByCustomer is a compatibility wrapper. Request handlers use UpdateCustomer_access_tokensByCustomerContext.
func UpdateCustomer_access_tokensByCustomer(m *Customer_access_tokens) (err error) {
	return UpdateCustomer_access_tokensByCustomerContext(context.Background(), m)
}

// DeleteCustomer_access_tokens is a compatibility wrapper. Request handlers use DeleteCustomer_access_tokensContext.
func DeleteCustomer_access_tokens(id int64) (err error) {
	return DeleteCustomer_access_tokensContext(context.Background(), id)
}
