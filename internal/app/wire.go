//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"

	"github.com/Koshsky/erp-backend/internal/auth"
	"github.com/Koshsky/erp-backend/internal/config"
	"github.com/Koshsky/erp-backend/internal/database"
	"github.com/Koshsky/erp-backend/internal/logger"
	authMw "github.com/Koshsky/erp-backend/internal/middleware/auth"
	rbacMW "github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/planning"
	"github.com/Koshsky/erp-backend/internal/policies"
	"github.com/Koshsky/erp-backend/internal/project_mgmt"
	"github.com/Koshsky/erp-backend/internal/security/jwt"
	"github.com/Koshsky/erp-backend/internal/timesheet"
	"github.com/Koshsky/erp-backend/internal/user"
)

// InitializeApp builds the whole application dependency graph.
func InitializeApp() (*App, error) {
	wire.Build(
		config.ProvideConfig,
		logger.ProvideLogger,
		database.ProvidePostgresDB,
		config.ProvidePostgresConfig,
		config.ProvideJWTConfig,
		jwt.ProvideJWTService,
		authMw.ProvideAuthMiddleware,
		rbacMW.ProvideMiddleware,
		policies.ProvideAll,
		ProvideRBACData,

		user.ProviderSet,
		auth.ProviderSet,
		planning.ProviderSet,
		project_mgmt.ProviderSet,
		timesheet.ProviderSet,

		ProvideModules,
		New,
	)
	return nil, nil
}
