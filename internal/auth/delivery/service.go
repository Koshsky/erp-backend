package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/auth/dto"
)

type AuthService interface {
	Register(
		ctx context.Context,
		lastName, firstName, middleName, username, password string,
	) (*dto.AuthResponse, error)
	Login(ctx context.Context, username, password string) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
}
