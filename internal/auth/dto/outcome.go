package dto

import "github.com/Koshsky/erp-backend/internal/security/jwt"

type RefreshResponse struct {
	Tokens  *jwt.TokenPair `json:"tokens"`
	Message string         `json:"message" example:"Token refreshed successfully"`
}

// AuthResponse is the login/register payload: token fields are flattened to
// the top level (access_token, token_type, expires_in, refresh_token) plus user.
type AuthResponse struct {
	AccessToken  string   `json:"access_token"  example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType    string   `json:"token_type"    example:"Bearer"`
	ExpiresIn    int      `json:"expires_in"    example:"3600"`
	RefreshToken string   `json:"refresh_token" example:"abcdef123456..."`
	User         UserInfo `json:"user"`
}

type UserInfo struct {
	ID       int64  `json:"id"       example:"1"`
	Name     string `json:"name"     example:"Ivan Ivanov"`
	Username string `json:"username" example:"ivanov"`
	Role     string `json:"role"     example:"rp"`
}
