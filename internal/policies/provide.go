package policies

import "github.com/Koshsky/erp-backend/internal/middleware/rbac"

// ProvideAll returns the default policies — the initial snapshot and fallback
// when the DB is unavailable/empty. At runtime the real set is provided by
// PolicyStore via rbac.Middleware.Refresh.
func ProvideAll() []rbac.Policy {
	policies, err := BuildPolicies(DefaultRouteSpecs())
	if err != nil {
		panic("policies: invalid default spec: " + err.Error())
	}
	return policies
}
