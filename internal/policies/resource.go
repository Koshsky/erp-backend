package policies

import (
	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

//nolint:gochecknoglobals // rule registry
var resourcePolicies = []rbac.Policy{
	{Name: "resource.list", Check: ListCheck(rbac.ResourceResource, "owner_id")},
	{Name: "resource.view", Check: EntityCheck(rbac.ResourceResource, ActionView)},
	{Name: "resource.update", Check: EntityCheck(rbac.ResourceResource, ActionUpdate)},
	{Name: "resource.delete", Check: EntityCheck(rbac.ResourceResource, ActionDelete)},
	{Name: "resource.create", Check: createResource},
	// Членство ресурса: список — видимость ресурса; добавление/снятие — управление ресурсом.
	{Name: "resource.member-list", Check: EntityCheck(rbac.ResourceResource, ActionView)},
	{Name: "resource.member-add", Check: EntityCheck(rbac.ResourceResource, ActionUpdate)},
	{Name: "resource.member-remove", Check: EntityCheck(rbac.ResourceResource, ActionUpdate)},
}

// createResource lets vp create a resource into their own ownership, admin — any owner.
func createResource(rc *rbac.CheckCtx) error {
	ownerID, _ := rc.BodyID("owner_id")
	if ownerID == 0 {
		ownerID = rc.User.ID
	}
	owners := rbac.Owners{Owner: ownerID}
	if !Authorize(rc.User.Role, rbac.ResourceResource, ActionCreate, owners, rc.User.ID) {
		return errors.ErrForbidden
	}
	return nil
}
