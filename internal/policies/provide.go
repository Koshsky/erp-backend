package policies

import "github.com/Koshsky/erp-backend/internal/middleware/rbac"

// ProvideAll returns all policies.
func ProvideAll() []rbac.Policy {
	return All()
}
