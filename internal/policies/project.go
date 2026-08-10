package policies

import (
	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

//nolint:gochecknoglobals // rule registry
var projectPolicies = []rbac.Policy{
	{Name: "project.list", Check: ListCheck(rbac.ResourceProject, "owner_id")},
	{Name: "project.view", Check: EntityCheck(rbac.ResourceProject, ActionView)},
	{Name: "project.update", Check: EntityCheck(rbac.ResourceProject, ActionUpdate)},
	{Name: "project.delete", Check: EntityCheck(rbac.ResourceProject, ActionDelete)},
	{Name: "project.create", Check: createProject},
}

// createProject lets rp create a project into their own ownership (owner
// defaults to self), admin — any owner.
func createProject(rc *rbac.CheckCtx) error {
	ownerID, _ := rc.BodyID("owner_id")
	if ownerID == 0 {
		ownerID = rc.User.ID
	}
	owners := rbac.Owners{ProjectOwner: ownerID}
	if !Authorize(rc.User.Role, rbac.ResourceProject, ActionCreate, owners, rc.User.ID) {
		return errors.ErrForbidden
	}
	return nil
}
