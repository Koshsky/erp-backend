package rbac

import (
	"log/slog"

	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"
)

func ProvideMiddleware(logger *slog.Logger, tracer *tracingpkg.Tracer, data Data, policies []Policy) *Middleware {
	byName := make(map[string]Policy, len(policies))
	for _, p := range policies {
		byName[p.Name] = p
	}
	return &Middleware{logger: logger, tracer: tracer, data: data, byName: byName}
}
