package auth

import (
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/security/jwt"
)

// ProvideAuthMiddleware builds the JWT auth middleware.
func ProvideAuthMiddleware(logger *slog.Logger, jwtService *jwt.Service) *Middleware {
	return &Middleware{
		logger:     logger,
		jwtManager: jwtService,
	}
}
