package rbac

import (
	"log/slog"

	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"
)

// ProvideMiddleware строит движок со стартовым набором политик (дефолты).
// На рантайме набор обновляет PolicyStore через Middleware.Refresh.
func ProvideMiddleware(logger *slog.Logger, tracer *tracingpkg.Tracer, data Data, policies []Policy) *Middleware {
	m := &Middleware{logger: logger, tracer: tracer, data: data}
	m.Refresh(policies)
	return m
}
