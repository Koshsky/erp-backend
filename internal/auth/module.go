// Package auth wires the auth module's providers and routes.
package auth

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/Koshsky/erp-backend/internal/auth/delivery"
	"github.com/Koshsky/erp-backend/internal/auth/repository"
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
	// loginResponseDelay uniformly slows every /auth/login response so that
	// timing does not reveal whether a username exists.
	loginResponseDelay = 500 * time.Millisecond
)

// refresh rate limit (per IP) — стpоже общего API-лимита: refresh-токены
// 256-бит (подбор нереален), лимит гасит злоупотребление/DoS по /auth/refresh.
const (
	refreshRatePerSecond  = 2.0
	refreshBurst          = 10
	refreshCleanupEvery   = time.Minute
	refreshLimitExpiresIn = 10 * time.Minute
)

// ProviderSet aggregates the auth module's dependencies.
var ProviderSet = wire.NewSet(
	service.NewAuthService,
	repository.NewAuthRepository,
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
// The uniform loginResponseDelay is applied before the limiter so that timing
// does not reveal whether a username exists.
func (m Module) loginGuard() gin.HandlerFunc {
	limiter := ratelimit.New(ratelimit.Config{
		RequestsPerSecond: loginRatePerSecond,
		Burst:             loginBurst,
		CleanupInterval:   loginCleanupEvery,
		Expiration:        loginLimitExpiresIn,
	}, m.logger)

	return func(c *gin.Context) {
		time.Sleep(loginResponseDelay)
		limiter(c)
	}
}

// refreshGuard returns a per-IP rate limiter applied to the endpoint /auth/refresh.
func (m Module) refreshGuard() gin.HandlerFunc {
	return ratelimit.New(ratelimit.Config{
		RequestsPerSecond: refreshRatePerSecond,
		Burst:             refreshBurst,
		CleanupInterval:   refreshCleanupEvery,
		Expiration:        refreshLimitExpiresIn,
	}, m.logger)
}

// RegisterPublicRoutes registers the auth routes without authentication.
func (m Module) RegisterPublicRoutes(r *gin.RouterGroup) {
	m.handler.RegisterRoutes(r, m.loginGuard(), m.refreshGuard())
}

// RegisterProtectedRoutes is a no-op: auth has no protected routes.
func (m Module) RegisterProtectedRoutes(r *gin.RouterGroup) {
}
