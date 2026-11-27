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

type ActivationCodes struct {
	ActivationCodeId int64     `orm:"auto;column(activation_code_id)"`
	Code             string    `orm:"size(80)"`
	Number           string    `orm:"size(80)"`
	ExpiryDate       time.Time `orm:"type(datetime)"`
	DateCreated      time.Time `orm:"type(datetime)"`
	DateModified     time.Time `orm:"type(datetime)"`
	CreatedBy        int
	ModifiedBy       int
	Active           int
}

func init() {
	orm.RegisterModel(new(ActivationCodes))
}

func (u *ActivationCodes) TableName() string {
	return "activation_codes"
}

// AddActivationCodes insert a new ActivationCodes into database and returns
// last inserted Id on success.
func AddActivationCodesContext(ctx context.Context, m *ActivationCodes) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.InsertWithCtx(ctx, m)
	if err == nil {
		trace.Logger(ctx).Info("record created", "record_type", "activation_codes", "record_id", id)
	}
	return
}

// GetActivationCodesById retrieves ActivationCodes by Id. Returns error if
// Id doesn't exist
func GetActivationCodesByIdContext(ctx context.Context, id int64) (v *ActivationCodes, err error) {
	o := orm.NewOrm()
	v = &ActivationCodes{ActivationCodeId: id}
	if err = o.QueryTable(new(ActivationCodes)).Filter("ActivationCodeId", id).RelatedSel().OneWithCtx(ctx, v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetActivationCodesById retrieves ActivationCodes by Id. Returns error if
// Id doesn't exist
func GetActivationCodeByNumberContext(ctx context.Context, number string) (v *ActivationCodes, err error) {
	o := orm.NewOrm()
	v = &ActivationCodes{}
	if err = o.QueryTable(new(ActivationCodes)).Filter("Number", number).Filter("ExpiryDate__gt", time.Now()).RelatedSel().OneWithCtx(ctx, v); err == nil {
		return v, nil
	}

	return nil, err
}

// GetAllActivationCodes retrieves all ActivationCodes matches certain condition. Returns empty list if
// no records exist
func GetActivationCodesByNumberContext(ctx context.Context, number string) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ActivationCodes))
	// query k=v

	var l []ActivationCodes
	qs = qs.Filter("Number", number).Filter("ExpiryDate__gt", time.Now()).RelatedSel()
	if _, err = qs.AllWithCtx(ctx, &l); err == nil {
		for _, v := range l {
			ml = append(ml, v)
		}
		return ml, nil
	}
	return nil, err
}

// GetAllActivationCodes retrieves all ActivationCodes matches certain condition. Returns empty list if
// no records exist
func GetAllActivationCodesContext(ctx context.Context, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ActivationCodes))
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

	var l []ActivationCodes
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

// UpdateActivationCodes updates ActivationCodes by Id and returns error if
// the record to be updated doesn't exist
func UpdateActivationCodesByIdContext(ctx context.Context, m *ActivationCodes) (err error) {
	o := orm.NewOrm()
	v := ActivationCodes{ActivationCodeId: m.ActivationCodeId}

	// ascertain id exists in the database
	if err = o.ReadWithCtx(ctx, &v); err == nil {
		var num int64
		if num, err = o.UpdateWithCtx(ctx, m); err == nil {
			trace.Logger(ctx).Debug("database rows changed", "rows", num)
		}
	}
	return
}

// DeleteActivationCodes deletes ActivationCodes by Id and returns error if
// the record to be deleted doesn't exist
func DeleteActivationCodesContext(ctx context.Context, id int64) (err error) {
	o := orm.NewOrm()
	v := ActivationCodes{ActivationCodeId: id}
	// ascertain id exists in the database
	if err = o.ReadWithCtx(ctx, &v); err == nil {
		var num int64
		if num, err = o.DeleteWithCtx(ctx, &ActivationCodes{ActivationCodeId: id}); err == nil {
			trace.Logger(ctx).Debug("database rows changed", "rows", num)
		}
	}
	return
}

// AddActivationCodes is a compatibility wrapper. Request handlers use AddActivationCodesContext.
func AddActivationCodes(m *ActivationCodes) (id int64, err error) {
	return AddActivationCodesContext(context.Background(), m)
}

// GetActivationCodesById is a compatibility wrapper. Request handlers use GetActivationCodesByIdContext.
func GetActivationCodesById(id int64) (v *ActivationCodes, err error) {
	return GetActivationCodesByIdContext(context.Background(), id)
}

// GetActivationCodeByNumber is a compatibility wrapper. Request handlers use GetActivationCodeByNumberContext.
func GetActivationCodeByNumber(number string) (v *ActivationCodes, err error) {
	return GetActivationCodeByNumberContext(context.Background(), number)
}

// GetActivationCodesByNumber is a compatibility wrapper. Request handlers use GetActivationCodesByNumberContext.
func GetActivationCodesByNumber(number string) (ml []interface{}, err error) {
	return GetActivationCodesByNumberContext(context.Background(), number)
}

// GetAllActivationCodes is a compatibility wrapper. Request handlers use GetAllActivationCodesContext.
func GetAllActivationCodes(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	return GetAllActivationCodesContext(context.Background(), query, fields, sortby, order, offset, limit)
}

// UpdateActivationCodesById is a compatibility wrapper. Request handlers use UpdateActivationCodesByIdContext.
func UpdateActivationCodesById(m *ActivationCodes) (err error) {
	return UpdateActivationCodesByIdContext(context.Background(), m)
}

// DeleteActivationCodes is a compatibility wrapper. Request handlers use DeleteActivationCodesContext.
func DeleteActivationCodes(id int64) (err error) {
	return DeleteActivationCodesContext(context.Background(), id)
}
