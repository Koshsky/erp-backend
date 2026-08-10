package rbac

import "log/slog"

func ProvideMiddleware(logger *slog.Logger, data Data, policies []Policy) *Middleware {
	byName := make(map[string]Policy, len(policies))
	for _, p := range policies {
		byName[p.Name] = p
	}
	return &Middleware{logger: logger, data: data, byName: byName}
}
