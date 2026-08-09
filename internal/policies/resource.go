package policies

import "github.com/Koshsky/erp-backend/internal/middleware/rbac"

//nolint:gochecknoglobals // rule registry
var resourcePolicies = []rbac.Policy{
	{Name: "resource.view", Check: EntityCheck(rbac.ResourceResource, ActionView)},
	{Name: "resource.update", Check: EntityCheck(rbac.ResourceResource, ActionUpdate)},
	{Name: "resource.delete", Check: EntityCheck(rbac.ResourceResource, ActionDelete)},
	{Name: "resource.create", Check: createResource},
}

// createResource lets vp create a resource into their own ownership, admin — any owner.
func createResource(rc *rbac.CheckCtx) error {
	ownerID, _ := rc.BodyID("owner_id")
	if ownerID == 0 {
		ownerID = rc.User.ID
	}
	owners := rbac.Owners{Owner: ownerID}
	if !Authorize(rc.User.Role, rbac.ResourceResource, ActionCreate, owners, rc.User.ID) {
		return rbac.ErrForbidden
	}
	return nil
}
