package requestsDTOs

type ChangePassword struct {
	OldPassword string
	NewPassword string
}

type ResetPassword struct {
	NewPassword string
}

type ResetPasswordLink struct {
	Email   string
	Message string
	Subject string
	Links   []*string
}
