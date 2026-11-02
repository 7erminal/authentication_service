package requestsDTOs

type UsernameDTO struct {
	Username string
}

type PinAuthenticationDTO struct {
	Username string `orm:"size(255)"`
	Pin      string `orm:"size(255)"`
}
