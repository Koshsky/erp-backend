package policies

import (
	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

//nolint:gochecknoglobals // rule registry
var assignmentPolicies = []rbac.Policy{
	{Name: "assignment.list", Check: ListCheck(rbac.ResourceAssignment, "owner_id")},
	{Name: "assignment.view", Check: EntityCheck(rbac.ResourceAssignment, ActionView)},
	{Name: "assignment.update", Check: EntityCheck(rbac.ResourceAssignment, ActionUpdate)},
	{Name: "assignment.delete", Check: EntityCheck(rbac.ResourceAssignment, ActionDelete)},
	{Name: "assignment.create", Check: createAssignment},
}

// createAssignment checks the matrix permission by the task (admin/vp in their
// own process) plus the business rule: a resource can only be assigned to a
// task of its own owner. Admin is exempt from the business rule: the matrix
// grants full access (ScopeAll), and the owner-matching rule exists to keep
// vp inside their own owners, not to restrict admin.
func createAssignment(rc *rbac.CheckCtx) error {
	taskID, err := rc.BodyID("task_id")
	if err != nil {
		return err
	}
	resourceID, err := rc.BodyID("resource_id")
	if err != nil {
		return err
	}

	taskOwners, err := rc.Owners(rbac.ResourceTask, taskID)
	if err != nil {
		return err
	}
	if !Authorize(rc.User.Role, rbac.ResourceAssignment, ActionCreate, taskOwners, rc.User.ID) {
		return errors.Forbidden(
			"недостаточно прав: назначать ресурс на задачу может только " +
				"владелец её процесса (или администратор)",
		)
	}

	if rc.User.Role != userdomain.Admin {
		resourceOwners, ownerErr := rc.Owners(rbac.ResourceResource, resourceID)
		if ownerErr != nil {
			return ownerErr
		}
		if !taskOwners.SharesOwner(resourceOwners) {
			return errors.Forbidden("ресурс не принадлежит владельцу задачи")
		}
	}
	return nil
}
