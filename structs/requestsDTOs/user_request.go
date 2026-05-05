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
