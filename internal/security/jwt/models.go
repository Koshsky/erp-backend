package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims

	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"            example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token,omitempty" example:"abcdef123456..."`
	TokenType    string `json:"token_type"              example:"Bearer"`
	ExpiresIn    int    `json:"expires_in"              example:"3600"` // seconds
}
