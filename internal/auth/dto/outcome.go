package dto

import "github.com/Koshsky/erp-backend/internal/security/jwt"

type RefreshResponse struct {
	Tokens  *jwt.TokenPair `json:"tokens"`
	Message string         `json:"message" example:"Token refreshed successfully"`
}

type AuthResponse struct {
	User   UserInfo       `json:"user"`
	Tokens *jwt.TokenPair `json:"tokens"`
}

type ChangePasswordResponse struct {
	Message string `json:"message" example:"password changed"`
}

type UserInfo struct {
	ID       int64  `json:"id" example:"1"`
	Name     string `json:"name" example:"Ivan Ivanov"`
	Username string `json:"username" example:"ivanov"`
	Role     string `json:"role" example:"РП"`
}
