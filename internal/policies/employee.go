package policies

import (
	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

//nolint:gochecknoglobals // rule registry
var employeePolicies = []rbac.Policy{
	{Name: "employee.list", Check: ListCheck(rbac.ResourceEmployee, "manager_id")},
	{Name: "employee.view", Check: EntityCheck(rbac.ResourceEmployee, ActionView)},
	{Name: "employee.update", Check: EntityCheck(rbac.ResourceEmployee, ActionUpdate)},
	{Name: "employee.delete", Check: EntityCheck(rbac.ResourceEmployee, ActionDelete)},
	{Name: "employee.create", Check: createEmployee},
}

// createEmployee lets vp create an employee into their own team, admin — anyone.
func createEmployee(rc *rbac.CheckCtx) error {
	ownerID, _ := rc.BodyID("manager_id")
	if ownerID == 0 {
		ownerID = rc.User.ID
	}
	owners := rbac.Owners{Owner: ownerID}
	if !Authorize(rc.User.Role, rbac.ResourceEmployee, ActionCreate, owners, rc.User.ID) {
		return errors.ErrForbidden
	}
	return nil
}
