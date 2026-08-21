package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/auth/dto"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
}
