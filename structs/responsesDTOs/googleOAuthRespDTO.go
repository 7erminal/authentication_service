package responsesDTOs

type GoogleOAuthRespDTO struct {
	Access_token  string
	Expires_in    int
	Scope         string
	Token_type    string
	Refresh_token string
}
