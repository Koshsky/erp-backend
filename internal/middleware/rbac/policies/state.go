package policies

import "github.com/Koshsky/erp-backend/internal/middleware/rbac"

//nolint:gochecknoglobals // rule registry
var statePolicies = []rbac.Policy{
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
				return rbac.ErrNotFound
			}
			return rbac.ErrForbidden
		}
		return nil
	}
}
