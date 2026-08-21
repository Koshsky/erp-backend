package policies

import (
	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

//nolint:gochecknoglobals // rule registry
var workerPolicies = []rbac.Policy{
	{Name: "worker.list", Check: ListCheck(rbac.ResourceWorker, "manager_id")},
	{Name: "worker.view", Check: EntityCheck(rbac.ResourceWorker, ActionView)},
	{Name: "worker.update", Check: EntityCheck(rbac.ResourceWorker, ActionUpdate)},
	{Name: "worker.delete", Check: EntityCheck(rbac.ResourceWorker, ActionDelete)},
	{Name: "worker.create", Check: createWorker},
	{Name: "user.picker", Check: userPicker},
}

// userPicker допускает к нескоупированному пулу пользователей (GET /user)
// только роли с планировочными страницами: им каталог нужен для имён
// владельцев в тултипах и пикеров. Worker (и любые неизвестные роли)
// перечислять каталог не могут.
func userPicker(rc *rbac.CheckCtx) error {
	switch rc.User.Role {
	case userdomain.Admin,
		userdomain.ProjectDirector,
		userdomain.ProjectManager,
		userdomain.ProcessOwner:
		return nil
	default:
		return errors.ErrForbidden
	}
}

// createWorker lets vp create a worker into their own team, admin — anyone.
func createWorker(rc *rbac.CheckCtx) error {
	ownerID, _ := rc.BodyID("manager_id")
	if ownerID == 0 {
		ownerID = rc.User.ID
	}
	owners := rbac.Owners{Owner: ownerID}
	if !Authorize(rc.User.Role, rbac.ResourceWorker, ActionCreate, owners, rc.User.ID) {
		return errors.ErrForbidden
	}
	return nil
}
