package models

type UpdateUserRequestDTO struct {
	FullName      string `orm:"size(255)"`
	Username      string `orm:"size(255)" orm: "omitempty"`
	PhoneNumber   string `orm:"size(255)" orm: "omitempty"`
	Gender        string `orm:"size(10)" orm: "omitempty"`
	Dob           string `orm:"size(50)" orm: "omitempty"`
	Address       string `orm:"size(255)" orm: "omitempty"`
	MaritalStatus string `orm:"size(255)" orm: "omitempty"`
}
