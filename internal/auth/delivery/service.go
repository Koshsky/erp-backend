package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/auth/dto"
)

type AuthService interface {
	Register(ctx context.Context, name, username, password string) (*dto.AuthResponse, error)
	Login(ctx context.Context, username, password string) (*dto.AuthResponse, error)
	ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error
	RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
}
