package dto

import "github.com/Koshsky/erp-backend/internal/security/jwt"

type RefreshResponse struct {
	Tokens  *jwt.TokenPair `json:"tokens"`
	Message string         `json:"message"`
}

type AuthResponse struct {
	User   UserInfo       `json:"user"`
	Tokens *jwt.TokenPair `json:"tokens"`
}

type UserInfo struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
