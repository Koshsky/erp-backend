// Package auth wires the auth module's providers and routes.
package auth

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/Koshsky/erp-backend/internal/auth/delivery"
	"github.com/Koshsky/erp-backend/internal/auth/service"
	"github.com/Koshsky/erp-backend/internal/middleware/ratelimit"
)

// login rate limit (per IP) — stricter than the global API limit to slow
// down password brute-forcing on /auth/login.
const (
	loginRatePerSecond  = 1.0
	loginBurst          = 5
	loginCleanupEvery   = time.Minute
	loginLimitExpiresIn = 10 * time.Minute
)

// ProviderSet aggregates the auth module's dependencies.
var ProviderSet = wire.NewSet(
	service.NewAuthService,
	delivery.NewAuthHandler,
	ProvideModule,
)

// Module registers the auth module's routes. All auth routes are public
// (register/login/refresh); there are no protected auth routes.
type Module struct {
	handler *delivery.AuthHandler
	logger  *slog.Logger
}

// ProvideModule builds the auth module.
func ProvideModule(handler *delivery.AuthHandler, logger *slog.Logger) Module {
	return Module{handler: handler, logger: logger}
}

// loginGuard returns a per-IP rate limiter applied to the login endpoint.
func (m Module) loginGuard() gin.HandlerFunc {
	return ratelimit.New(ratelimit.Config{
		RequestsPerSecond: loginRatePerSecond,
		Burst:             loginBurst,
		CleanupInterval:   loginCleanupEvery,
		Expiration:        loginLimitExpiresIn,
	}, m.logger)
}

// RegisterPublicRoutes registers the auth routes without authentication.
func (m Module) RegisterPublicRoutes(r *gin.RouterGroup) {
	m.handler.RegisterRoutes(r, m.loginGuard())
}

// RegisterProtectedRoutes is a no-op: auth has no protected routes.
func (m Module) RegisterProtectedRoutes(r *gin.RouterGroup) {
}
