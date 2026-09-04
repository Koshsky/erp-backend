package rbac

import (
	"log/slog"

	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"
)

// ProvideMiddleware builds the engine with a startup set of policies (defaults).
// At runtime the set is updated by PolicyStore via Middleware.Refresh.
func ProvideMiddleware(logger *slog.Logger, tracer *tracingpkg.Tracer, data Data, policies []Policy) *Middleware {
	m := &Middleware{logger: logger, tracer: tracer, data: data}
	m.Refresh(policies)
	return m
}
