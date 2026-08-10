// Package user wires the user module's providers and routes.
package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/Koshsky/erp-backend/internal/user/delivery"
	"github.com/Koshsky/erp-backend/internal/user/repository"
	"github.com/Koshsky/erp-backend/internal/user/service"
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
}

// ProvideModule builds the user module.
func ProvideModule(handler *delivery.UserHandler) Module {
	return Module{handler: handler}
}

// RegisterPublicRoutes is a no-op: the user module has no public routes.
func (m Module) RegisterPublicRoutes(r *gin.RouterGroup) {
}

// RegisterProtectedRoutes registers the user routes behind authentication.
func (m Module) RegisterProtectedRoutes(r *gin.RouterGroup) {
	m.handler.RegisterRoutes(r)
}
