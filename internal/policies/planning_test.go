package policies_test

import (
	"testing"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/policies"
)

// DefaultRouteSpecs carry the planning.* route policies (mirror seed V13):
// /planning/* is gated by the view matrix of the underlying domain via the
// list kind (row scoping stays in the SQL).
func TestPlanningRouteSpecs(t *testing.T) {
	t.Parallel()
	byName := make(map[string]policies.RouteSpec, len(policies.DefaultRouteSpecs()))
	for _, spec := range policies.DefaultRouteSpecs() {
		byName[spec.Name] = spec
	}
	for _, name := range []string{"planning.projects", "planning.processes", "planning.tasks"} {
		spec, ok := byName[name]
		if !ok {
			t.Fatalf("DefaultRouteSpecs не содержит %s", name)
		}
		if err := policies.ValidateSpec(spec); err != nil {
			t.Errorf("ValidateSpec(%s) неожиданная ошибка: %v", name, err)
		}
	}
}

// The list kind checks the view right: roles without a view scope are denied,
// admin (bypass) and viewing roles pass.
func TestPlanningListCheck(t *testing.T) {
	t.Parallel()
	policies.SetMatrix(policies.DefaultMatrix())
	defer policies.SetMatrix(policies.DefaultMatrix())

	for _, act := range []struct {
		resource rbac.Resource
		viewRole string
		noRole   string
	}{
		{rbac.ResourceProject, "rp", "worker"},
		{rbac.ResourceProcess, "vp", "worker"},
		{rbac.ResourceTask, "vp", "worker"},
	} {
		if got := policies.ViewScopeCode(act.noRole, act.resource); got != "" {
			t.Errorf("ViewScopeCode(%s, %v) = %q; want empty (no access)", act.noRole, act.resource, got)
		}
		if got := policies.ViewScopeCode(act.viewRole, act.resource); got == "" {
			t.Errorf("ViewScopeCode(%s, %v) пуст; want a scope", act.viewRole, act.resource)
		}
	}
	if got := policies.ViewScopeCode("admin", rbac.ResourceProject); got != "all" {
		t.Errorf("ViewScopeCode(admin, project) = %q; want all", got)
	}
}
