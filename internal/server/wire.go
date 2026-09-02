//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"

	"github.com/Koshsky/erp-backend/internal/audit"
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
	projectmgmt "github.com/Koshsky/erp-backend/internal/project_mgmt"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy"
	"github.com/Koshsky/erp-backend/internal/security/jwt"
	"github.com/Koshsky/erp-backend/internal/server/profiler"
	"github.com/Koshsky/erp-backend/internal/timesheet"
	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"
	"github.com/Koshsky/erp-backend/internal/user"
	userservice "github.com/Koshsky/erp-backend/internal/user/service"
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
		config.ProvideRBACRefreshInterval,
		config.ProvideAuditConfig,
		idempotency.ProvideIdempotencyRepository,
		idempotency.ProvideIdempotencyMiddleware,
		jwt.ProvideJWTService,
		authMw.ProvideAuthMiddleware,
		rbacMW.ProvideMiddleware,
		policies.ProvideAll,
		profiler.ProvideProfiler,
		ProvideRBACData,
		// The audit client consumes the user service to resolve the `user`
		// filter (login/full name) and to enrich actor display names.
		wire.Bind(new(audit.UserLookup), new(*userservice.UserService)),

		user.ProviderSet,
		auth.ProviderSet,
		planning.ProviderSet,
		projectmgmt.ProviderSet,
		timesheet.ProviderSet,
		autocreate.ProviderSet,
		rbacpolicy.ProviderSet,
		audit.ProviderSet,

		ProvideModules,
		New,
	)
	return nil, nil
}
