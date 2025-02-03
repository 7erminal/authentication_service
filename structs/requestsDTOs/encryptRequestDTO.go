package requestsDTOs

type EncryptInviteRequestDTO struct {
	Email string
	Role  string
}

type DecryptRequestDTO struct {
	Token string
	Nonce string
	Email string
}

type TokenDTO struct {
	Token string
}
