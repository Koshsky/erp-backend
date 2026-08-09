package delivery

import (
	"log/slog"

	authservice "github.com/Koshsky/erp-backend/internal/auth/service"
)

// ProvideAuthHandler builds the auth handler.
func ProvideAuthHandler(logger *slog.Logger, svc *authservice.AuthService) *AuthHandler {
	return &AuthHandler{
		logger:  logger,
		service: svc,
	}
}
