package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
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
func AddActivationCodes(m *ActivationCodes) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetActivationCodesById retrieves ActivationCodes by Id. Returns error if
// Id doesn't exist
func GetActivationCodesById(id int64) (v *ActivationCodes, err error) {
	o := orm.NewOrm()
	v = &ActivationCodes{ActivationCodeId: id}
	if err = o.QueryTable(new(ActivationCodes)).Filter("ActivationCodeId", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetActivationCodesById retrieves ActivationCodes by Id. Returns error if
// Id doesn't exist
func GetActivationCodeByNumber(number string) (v *ActivationCodes, err error) {
	o := orm.NewOrm()
	v = &ActivationCodes{}
	if err = o.QueryTable(new(ActivationCodes)).Filter("Number", number).Filter("ExpiryDate__gt", time.Now()).RelatedSel().One(v); err == nil {
		return v, nil
	}
	logs.Info("Error fetching activation code ", err.Error())
	return nil, err
}

// GetAllActivationCodes retrieves all ActivationCodes matches certain condition. Returns empty list if
// no records exist
func GetActivationCodesByNumber(number string) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ActivationCodes))
	// query k=v

	var l []ActivationCodes
	qs = qs.Filter("Number", number).Filter("ExpiryDate__gt", time.Now()).RelatedSel()
	if _, err = qs.All(&l); err == nil {
		for _, v := range l {
			ml = append(ml, v)
		}
		return ml, nil
	}
	return nil, err
}

// GetAllActivationCodes retrieves all ActivationCodes matches certain condition. Returns empty list if
// no records exist
func GetAllActivationCodes(query map[string]string, fields []string, sortby []string, order []string,
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

// UpdateActivationCodes updates ActivationCodes by Id and returns error if
// the record to be updated doesn't exist
func UpdateActivationCodesById(m *ActivationCodes) (err error) {
	o := orm.NewOrm()
	v := ActivationCodes{ActivationCodeId: m.ActivationCodeId}
	logs.Debug("Activation ID in model is ", v.ActivationCodeId, " and value is ", v.ExpiryDate)
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteActivationCodes deletes ActivationCodes by Id and returns error if
// the record to be deleted doesn't exist
func DeleteActivationCodes(id int64) (err error) {
	o := orm.NewOrm()
	v := ActivationCodes{ActivationCodeId: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ActivationCodes{ActivationCodeId: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
