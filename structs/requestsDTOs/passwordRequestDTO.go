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
	Role    string
	Message string
	Subject string
}
