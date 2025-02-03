package requestsDTOs

type ChangePassword struct {
	OldPassword string
	NewPassword string
}

type ResetPassword struct {
	NewPassword string
}
