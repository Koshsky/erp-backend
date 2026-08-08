// Package policies is the single place with all access rules. Files are split
// by entity; the engine (rbac) stays a pure mechanism.
package policies

import "github.com/Koshsky/erp-backend/internal/middleware/rbac"

// All returns every policy in the system.
func All() []rbac.Policy {
	return concat(
		projectPolicies,
		processPolicies,
		taskPolicies,
		milestonePolicies,
		assignmentPolicies,
		resourcePolicies,
		employeePolicies,
		statePolicies,
	)
}

func concat(lists ...[]rbac.Policy) []rbac.Policy {
	var out []rbac.Policy
	for _, l := range lists {
		out = append(out, l...)
	}
	return out
}
