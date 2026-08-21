package policies

import "github.com/Koshsky/erp-backend/internal/middleware/rbac"

//nolint:gochecknoglobals // rule registry
var taskPolicies = []rbac.Policy{
	{Name: "task.list", Check: ListCheck(rbac.ResourceTask, "owner_id")},
	{Name: "task.view", Check: EntityCheck(rbac.ResourceTask, ActionView)},
	{Name: "task.update", Check: EntityCheck(rbac.ResourceTask, ActionUpdate)},
	{Name: "task.delete", Check: EntityCheck(rbac.ResourceTask, ActionDelete)},
	{Name: "task.create", Check: CreateCheck(rbac.ResourceTask, parentByID(rbac.ResourceProcess, "process_id"))},
}
