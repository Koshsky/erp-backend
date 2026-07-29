package service

import (
	"context"

	userDTO "github.com/Koshsky/erp-backend/internal/user/dto"
)

// UserService is the interface auth depends on instead of a direct repository.
type UserService interface {
	FindUserByUsername(ctx context.Context, username string) (*userDTO.UserResponse, error)
	FindUserByID(ctx context.Context, userID int64) (*userDTO.UserResponse, error)
	CreateUser(ctx context.Context, req userDTO.CreateUserRequest) (*userDTO.UserResponse, error)
	UpdatePassword(ctx context.Context, userID int64, newHash string) error
}

// PasswordHasher is used to hash and compare passwords.
type PasswordHasher interface {
	Hash(raw string) (string, error)
	Compare(hashed, raw string) error
}
