package dto

// AuthResponse is the login/refresh payload. The refresh token itself is NOT
// returned here: it lives in an HttpOnly cookie (AD-05).
type AuthResponse struct {
	AccessToken string   `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string   `json:"token_type"   example:"Bearer"`
	ExpiresIn   int      `json:"expires_in"   example:"900"`
	User        UserInfo `json:"user"`
}

// SessionResult is what the auth service returns to the delivery layer:
// the response body plus the opaque refresh token for the HttpOnly cookie.
type SessionResult struct {
	Auth         *AuthResponse
	RefreshToken string
}

type UserInfo struct {
	ID         int64   `json:"id"          example:"1"`
	Name       string  `json:"name"        example:"Иванов Иван Иванович"`
	LastName   string  `json:"last_name"   example:"Иванов"`
	FirstName  string  `json:"first_name"  example:"Иван"`
	MiddleName *string `json:"middle_name" example:"Иванович"`
	Username   string  `json:"username"    example:"ivanov"`
	Role       string  `json:"role"        example:"rp"`
}
