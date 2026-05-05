package requestsDTOs

type CustomerCredentialRequestDTO struct {
	CustomerId int64
	Username   string
	Password   string
	Pin        string
}

type CustomerCredentialUpdateRequestDTO struct {
	Username string
	Password string
	Pin      string
}

type ValidatePin struct {
	Username string
	Password string
}

type ResetPin struct {
	Number   string
	NewPin   string
	Password string
}
