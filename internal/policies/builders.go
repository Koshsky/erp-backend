package policies

import (
	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

// parentByID resolves the parent owner from the body by key (for CreateCheck).
func parentByID(rsrc rbac.Resource, key string) func(*rbac.CheckCtx) (rbac.Owners, error) {
	return func(rc *rbac.CheckCtx) (rbac.Owners, error) {
		id, err := rc.BodyID(key)
		if err != nil {
			return rbac.Owners{}, err
		}
		return rc.Owners(rsrc, id)
	}
}

// EntityCheck is the standard rule for an id from the URL: checks the matrix
// against the owners. A denied view does not reveal the record existence (404).
func EntityCheck(rsrc rbac.Resource, act Action) func(*rbac.CheckCtx) error {
	return func(rc *rbac.CheckCtx) error {
		id, err := rc.ParamID()
		if err != nil {
			return err
		}
		owners, err := rc.Owners(rsrc, id)
		if err != nil {
			return err
		}
		if !Authorize(rc.User.Role, rsrc, act, owners, rc.User.ID) {
			if act == ActionView {
				return errors.ErrNotFound
			}
			return errors.ErrForbidden
		}
		return nil
	}
}

// CreateCheck is the standard create rule: by the parent owner from the body.
func CreateCheck(rsrc rbac.Resource, parent func(*rbac.CheckCtx) (rbac.Owners, error)) func(*rbac.CheckCtx) error {
	return func(rc *rbac.CheckCtx) error {
		owners, err := parent(rc)
		if err != nil {
			return err
		}
		if !Authorize(rc.User.Role, rsrc, ActionCreate, owners, rc.User.ID) {
			return errors.ErrForbidden
		}
		return nil
	}
}
