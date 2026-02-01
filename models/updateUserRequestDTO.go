package models

type UpdateUserRequestDTO struct {
	FullName      string `orm:"size(255)"`
	Username      string `orm:"size(255)"`
	PhoneNumber   string `orm:"size(255)"`
	Gender        string `orm:"size(10)"`
	Dob           string `orm:"size(50)"`
	Address       string `orm:"size(255)"`
	MaritalStatus string `orm:"size(255)"`
}
