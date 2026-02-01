package responsesDTOs

type TokenResponseDTO struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int64
}

type LoginTokenResponseDTO struct {
	StatusCode int
	Result     *TokenResponseDTO
	StatusDesc string
}
