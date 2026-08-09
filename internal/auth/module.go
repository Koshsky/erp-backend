// Package auth wires the auth module's providers and routes.
package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/Koshsky/erp-backend/internal/auth/delivery"
	"github.com/Koshsky/erp-backend/internal/auth/service"
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
}

// ProvideModule builds the auth module.
func ProvideModule(handler *delivery.AuthHandler) Module {
	return Module{handler: handler}
}

// RegisterPublicRoutes registers the auth routes without authentication.
func (m Module) RegisterPublicRoutes(r *gin.RouterGroup) {
	m.handler.RegisterRoutes(r)
}

// RegisterProtectedRoutes is a no-op: auth has no protected routes.
func (m Module) RegisterProtectedRoutes(r *gin.RouterGroup) {
}
