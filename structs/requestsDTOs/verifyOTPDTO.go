package requestsDTOs

type VerifyOtpDTO struct {
	Username string
	Password string
}

type SendActivationCode struct {
	MobileNumber string
}

type VerifyActivationCodeDTO struct {
	MobileNumber string
	Password     string
}
