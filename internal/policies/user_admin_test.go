package policies_test

import (
	"testing"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/policies"
)

// The user-edit right (user_admin.*): an employee IS a system user, so profile
// mutations are gated by the user_admin virtual resource. The default matrix
// grants it to nobody but admin (the bypass); it stays grantable via the
// matrix (scope "all" is the only applicable one).
func TestUserAdminRights(t *testing.T) {
	t.Parallel()
	for _, act := range []policies.Action{
		policies.ActionCreate, policies.ActionUpdate, policies.ActionDelete,
	} {
		if got := policies.DefaultMatrix().ScopeFor("vp", rbac.ResourceUserAdmin, act); got != policies.ScopeNone {
			t.Errorf("ScopeFor(vp, user_admin, %v) = %v; want none", policies.ActionName(act), got)
		}
		if got := policies.DefaultMatrix().ScopeFor("admin", rbac.ResourceUserAdmin, act); got != policies.ScopeAll {
			t.Errorf("ScopeFor(admin, user_admin, %v) = %v; want all", policies.ActionName(act), got)
		}
		if policies.Authorize("vp", rbac.ResourceUserAdmin, act, rbac.Owners{}, 42) {
			t.Errorf("Authorize(vp, user_admin, %v) = true; want false", policies.ActionName(act))
		}
		if !policies.Authorize("admin", rbac.ResourceUserAdmin, act, rbac.Owners{}, 42) {
			t.Errorf("Authorize(admin, user_admin, %v) = false; want true", policies.ActionName(act))
		}
	}
}

// DefaultRouteSpecs carry the user_admin.* route policies (mirror seed V12)
// and they are valid (nullable owner — the virtual resource has no owner).
func TestUserAdminRouteSpecs(t *testing.T) {
	t.Parallel()
	byName := make(map[string]policies.RouteSpec, len(policies.DefaultRouteSpecs()))
	for _, spec := range policies.DefaultRouteSpecs() {
		byName[spec.Name] = spec
	}
	for _, name := range []string{"user_admin.create", "user_admin.update", "user_admin.delete"} {
		spec, ok := byName[name]
		if !ok {
			t.Fatalf("DefaultRouteSpecs не содержит %s", name)
		}
		if err := policies.ValidateSpec(spec); err != nil {
			t.Errorf("ValidateSpec(%s) неожиданная ошибка: %v", name, err)
		}
	}
}
