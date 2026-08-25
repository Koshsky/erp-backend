// Package rbacpolicy wires the module that stores and serves runtime
// RBAC policies (matrix + route checks) from Postgres.
package rbacpolicy

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/Koshsky/erp-backend/internal/rbacpolicy/delivery"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy/repository"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy/service"
)

// ProviderSet aggregates the rbacpolicy module's dependencies.
//
//nolint:gochecknoglobals // wire provider set (established module pattern)
var ProviderSet = wire.NewSet(
	repository.NewRuleRepository,
	service.NewPolicyStore,
	service.NewRBACService,
	delivery.NewRBACHandler,
	ProvideModule,
)

// Module registers the RBAC administration routes (all protected, admin-only).
type Module struct {
	handler *delivery.RBACHandler
	logger  *slog.Logger
}

// ProvideModule builds the rbacpolicy module.
func ProvideModule(handler *delivery.RBACHandler, logger *slog.Logger) Module {
	return Module{handler: handler, logger: logger}
}

// RegisterPublicRoutes is a no-op: the module has no public routes.
func (m Module) RegisterPublicRoutes(_ *gin.RouterGroup) {
}

// RegisterProtectedRoutes registers the admin routes behind authentication.
func (m Module) RegisterProtectedRoutes(r *gin.RouterGroup) {
	m.handler.RegisterRoutes(r)
}
