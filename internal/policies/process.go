package policies

import "github.com/Koshsky/erp-backend/internal/middleware/rbac"

//nolint:gochecknoglobals // rule registry
var processPolicies = []rbac.Policy{
	{Name: "process.list", Check: ListCheck(rbac.ResourceProcess, "owner_id")},
	{Name: "process.view", Check: EntityCheck(rbac.ResourceProcess, ActionView)},
	{Name: "process.update", Check: EntityCheck(rbac.ResourceProcess, ActionUpdate)},
	{Name: "process.delete", Check: EntityCheck(rbac.ResourceProcess, ActionDelete)},
	{Name: "process.create", Check: CreateCheck(rbac.ResourceProcess, parentByID(rbac.ResourceProject, "project_id"))},
}
