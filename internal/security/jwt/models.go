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
