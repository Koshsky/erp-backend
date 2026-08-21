// Package user wires the user module's providers and routes.
package user

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/Koshsky/erp-backend/internal/middleware/ratelimit"
	"github.com/Koshsky/erp-backend/internal/user/delivery"
	"github.com/Koshsky/erp-backend/internal/user/repository"
	"github.com/Koshsky/erp-backend/internal/user/service"
	userctx "github.com/Koshsky/erp-backend/internal/userctx"
)

// change-password limit (per user) — защита от брутфорса старого пароля.
// Задержка uniform, чтобы не раскрывать результат по времени ответа.
const (
	changePasswordRatePerSecond  = 0.5
	changePasswordBurst          = 3
	changePasswordCleanupEvery   = time.Minute
	changePasswordLimitExpiresIn = 10 * time.Minute
	changePasswordResponseDelay  = 300 * time.Millisecond
)

// ProviderSet aggregates the user module's dependencies.
var ProviderSet = wire.NewSet(
	repository.NewUserRepository,
	service.NewUserService,
	delivery.NewUserHandler,
	ProvideModule,
)

// Module registers the user module's routes (all protected).
type Module struct {
	handler *delivery.UserHandler
	logger  *slog.Logger
}

// ProvideModule builds the user module.
func ProvideModule(handler *delivery.UserHandler, logger *slog.Logger) Module {
	return Module{handler: handler, logger: logger}
}

// changePasswordGuard возвращает per-user (по JWT id) лимитер для
// /user/change-password; uniform задержка применяется до лимитера.
func (m Module) changePasswordGuard() gin.HandlerFunc {
	limiter := ratelimit.New(ratelimit.Config{
		RequestsPerSecond: changePasswordRatePerSecond,
		Burst:             changePasswordBurst,
		CleanupInterval:   changePasswordCleanupEvery,
		Expiration:        changePasswordLimitExpiresIn,
		Key: func(c *gin.Context) string {
			id, err := userctx.GetUserID(c)
			if err != nil {
				// На защищённом пути userctx обязан быть; на любой промах —
				// общий bucket, чтобы не плодить неограниченные лимитеры.
				return "unknown"
			}
			return strconv.FormatInt(id, 10)
		},
	}, m.logger)

	return func(c *gin.Context) {
		time.Sleep(changePasswordResponseDelay)
		limiter(c)
	}
}

// RegisterPublicRoutes is a no-op: the user module has no public routes.
func (m Module) RegisterPublicRoutes(r *gin.RouterGroup) {
}

// RegisterProtectedRoutes registers the user routes behind authentication.
func (m Module) RegisterProtectedRoutes(r *gin.RouterGroup) {
	m.handler.RegisterRoutes(r, m.changePasswordGuard())
}
