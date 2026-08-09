package policies

import "github.com/Koshsky/erp-backend/internal/middleware/rbac"

//nolint:gochecknoglobals // rule registry
var milestonePolicies = []rbac.Policy{
	{Name: "milestone.view", Check: EntityCheck(rbac.ResourceMilestone, ActionView)},
	{Name: "milestone.update", Check: EntityCheck(rbac.ResourceMilestone, ActionUpdate)},
	{Name: "milestone.delete", Check: EntityCheck(rbac.ResourceMilestone, ActionDelete)},
	{
		Name:  "milestone.create",
		Check: CreateCheck(rbac.ResourceMilestone, parentByID(rbac.ResourceProcess, "process_id")),
	},
}
