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

type Customer_credentials struct {
	Id           int64      `orm:"auto;column(customer_credential_id)"`
	Customer     *Customers `orm:"rel(fk)"`
	Username     string     `orm:"size(255)"`
	Password     string     `orm:"size(255)"`
	Pin          string     `orm:"size(10)"`
	DateCreated  time.Time  `orm:"type(datetime)"`
	DateModified time.Time  `orm:"type(datetime)"`
	CreatedBy    int
	ModifiedBy   int
	Active       int
}

func (t *Customer_credentials) TableName() string {
	return "customer_credentials"
}

func init() {
	orm.RegisterModel(new(Customer_credentials))
}

// AddCustomer_credentials insert a new Customer_credentials into database and returns
// last inserted Id on success.
func AddCustomer_credentials(m *Customer_credentials) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetCustomer_credentialsById retrieves Customer_credentials by Id. Returns error if
// Id doesn't exist
func GetCustomer_credentialsById(id int64) (v *Customer_credentials, err error) {
	o := orm.NewOrm()
	v = &Customer_credentials{Id: id}
	if err = o.QueryTable(new(Customer_credentials)).Filter("customer_credential_id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetCustomer_credentialsByCustomerId retrieves Customer_credentials by Id. Returns error if
// Id doesn't exist
func GetCustomer_credentialsByCustomerId(id Customers) (v *Customer_credentials, err error) {
	o := orm.NewOrm()
	v = &Customer_credentials{Customer: &id}
	if err = o.QueryTable(new(Customer_credentials)).Filter("Customer", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetCustomer_credentialsByCustomerId retrieves Customer_credentials by Id. Returns error if
// Id doesn't exist
func GetCustomer_credentialsByCustomerUsername(username string) (v *Customer_credentials, err error) {
	logs.Info("Fetching customer credentials by username:", username)
	o := orm.NewOrm()
	v = &Customer_credentials{}
	if err = o.QueryTable(new(Customer_credentials)).Filter("Username", username).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCustomer_credentials retrieves all Customer_credentials matches certain condition. Returns empty list if
// no records exist
func GetAllCustomer_credentials(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Customer_credentials))
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

	var l []Customer_credentials
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

// UpdateCustomer_credentials updates Customer_credentials by Id and returns error if
// the record to be updated doesn't exist
func UpdateCustomer_credentialsById(m *Customer_credentials) (err error) {
	o := orm.NewOrm()
	v := Customer_credentials{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteCustomer_credentials deletes Customer_credentials by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCustomer_credentials(id int64) (err error) {
	o := orm.NewOrm()
	v := Customer_credentials{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Customer_credentials{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
