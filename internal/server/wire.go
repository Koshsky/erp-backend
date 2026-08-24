//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"

	"github.com/Koshsky/erp-backend/internal/auth"
	autocreate "github.com/Koshsky/erp-backend/internal/auto_create"
	"github.com/Koshsky/erp-backend/internal/config"
	"github.com/Koshsky/erp-backend/internal/database"
	"github.com/Koshsky/erp-backend/internal/idempotency"
	"github.com/Koshsky/erp-backend/internal/logger"
	authMw "github.com/Koshsky/erp-backend/internal/middleware/auth"
	rbacMW "github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/planning"
	"github.com/Koshsky/erp-backend/internal/policies"
	"github.com/Koshsky/erp-backend/internal/project_mgmt"
	"github.com/Koshsky/erp-backend/internal/security/jwt"
	"github.com/Koshsky/erp-backend/internal/server/profiler"
	"github.com/Koshsky/erp-backend/internal/timesheet"
	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"
	"github.com/Koshsky/erp-backend/internal/user"
)

// InitializeApp builds the whole application dependency graph.
func InitializeApp() (*App, error) {
	wire.Build(
		config.ProvideConfig,
		logger.ProvideLogger,
		tracingpkg.ProvideTracer,
		config.ProvideTracingConfig,
		database.ProvidePostgresDB,
		config.ProvidePostgresConfig,
		config.ProvideJWTConfig,
		config.ProvideProfilingConfig,
		idempotency.ProvideIdempotencyRepository,
		idempotency.ProvideIdempotencyMiddleware,
		jwt.ProvideJWTService,
		authMw.ProvideAuthMiddleware,
		rbacMW.ProvideMiddleware,
		policies.ProvideAll,
		profiler.ProvideProfiler,
		ProvideRBACData,

		user.ProviderSet,
		auth.ProviderSet,
		planning.ProviderSet,
		project_mgmt.ProviderSet,
		timesheet.ProviderSet,
		autocreate.ProviderSet,

		ProvideModules,
		New,
	)
	return nil, nil
}
