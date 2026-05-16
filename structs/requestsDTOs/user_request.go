package requestsDTOs

type GetUserRequest struct {
	UserId string
}

type GetUserWithUsernameRequest struct {
	Username string
}

type GetCustomerRequest struct {
	CustomerId string
}

type GetCustomerWithUsernameRequest struct {
	Username string
}

type UpdateUserPasswordRequest struct {
	UserId      int64  `validate:"required"`
	OldPassword string `validate:"required"`
	NewPassword string `validate:"required"`
}
