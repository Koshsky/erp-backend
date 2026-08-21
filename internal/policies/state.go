package policies

import (
	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

//nolint:gochecknoglobals // rule registry
var statePolicies = []rbac.Policy{
	// Список справочника состояний = та же матрица view (admin/vp; по умолчанию 404 для остальных).
	{Name: "state.list", Check: stateCheck(ActionView)},
	{Name: "state.view", Check: stateCheck(ActionView)},
	{Name: "state.create", Check: stateCheck(ActionCreate)},
	{Name: "state.update", Check: stateCheck(ActionUpdate)},
	{Name: "state.delete", Check: stateCheck(ActionDelete)},
}

// stateCheck applies when states have no owner: access is decided by the matrix.
func stateCheck(act Action) func(*rbac.CheckCtx) error {
	return func(rc *rbac.CheckCtx) error {
		if !Authorize(rc.User.Role, rbac.ResourceState, act, rbac.Owners{}, rc.User.ID) {
			if act == ActionView {
				return errors.ErrNotFound
			}
			return errors.ErrForbidden
		}
		return nil
	}
}
