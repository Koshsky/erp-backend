package policies

import (
	"strconv"

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

// ListCheck is the standard rule for a listing route: denies roles without
// view access (ScopeNone) and validates an optional scope query param
// (owner_id/manager_id) — any id (or 0 = all) for ScopeAll roles, only the
// caller's own id otherwise. The SQL then filters by the caller's JWT scope
// and applies the validated param as an extra owner filter.
func ListCheck(rsrc rbac.Resource, key string) func(*rbac.CheckCtx) error {
	return func(rc *rbac.CheckCtx) error {
		scope := scopeFor(rc.User.Role, rsrc, ActionView)
		if scope == ScopeNone {
			return errors.ErrForbidden
		}
		raw := rc.C.Query(key)
		if raw == "" {
			return nil
		}
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || id < 0 {
			return errors.BadRequest("invalid " + key)
		}
		if scope != ScopeAll && id != rc.User.ID {
			return errors.ErrForbidden
		}
		return nil
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
