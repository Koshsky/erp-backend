package policies_test

import (
	"testing"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/policies"
)

// Кодеки: roundtrip для каждого ресурса/действия/зоны + невалидные значения.
func TestResourceCodecs(t *testing.T) {
	t.Parallel()
	for i := rbac.ResourceProject; i <= rbac.ResourceRBACConfig; i++ {
		name := policies.ResourceName(i)
		if name == "" {
			t.Fatalf("ResourceName(%v) пуст", i)
		}
		back, ok := policies.ParseResource(name)
		if !ok || back != i {
			t.Errorf("ParseResource(%q) = %v, %v; want %v", name, back, ok, i)
		}
	}
	for _, bogus := range []string{"bogus", "review", "Organization"} {
		if _, ok := policies.ParseResource(bogus); ok {
			t.Errorf("ParseResource(%q) должен быть false", bogus)
		}
	}
}

func TestActionCodecs(t *testing.T) {
	t.Parallel()
	for i := policies.ActionView; i <= policies.ActionDelete; i++ {
		name := policies.ActionName(i)
		back, ok := policies.ParseAction(name)
		if !ok || back != i {
			t.Errorf("ParseAction(%q) = %v, %v; want %v", name, back, ok, i)
		}
	}
	if _, ok := policies.ParseAction("edit"); ok {
		t.Errorf(`ParseAction("edit") должен быть false`)
	}
}

func TestScopeCodecs(t *testing.T) {
	t.Parallel()
	for _, scope := range []policies.Scope{policies.ScopeAll, policies.ScopeOwn, policies.ScopeParent, policies.ScopeAncestor} {
		name := policies.ScopeName(scope)
		back, ok := policies.ParseScope(name)
		if !ok || back != scope {
			t.Errorf("ParseScope(%q) = %v, %v; want %v", name, back, ok, scope)
		}
	}
	if scope, ok := policies.ParseScope(""); ok && scope != policies.ScopeNone {
		t.Errorf("ParseScope(\"\") = %v, %v; want none", scope, ok)
	}
	if _, ok := policies.ParseScope("bogus"); ok {
		t.Errorf(`ParseScope("bogus") должен быть false`)
	}
}

// Точная таблица применимости зон по ресурсу (зеркалит бэкенд ScopeApplicable).
func TestScopeApplicableTable(t *testing.T) {
	t.Parallel()
	ownRes := map[rbac.Resource]bool{
		rbac.ResourceProject:     true,
		rbac.ResourceProcess:     true,
		rbac.ResourceTask:        true,
		rbac.ResourceMilestone:   false,
		rbac.ResourceAssignment:  false,
		rbac.ResourceState:       false,
		rbac.ResourceResource:    true,
		rbac.ResourceWorker:      true,
		rbac.ResourceComment:     false,
		rbac.ResourceUserCatalog: false,
		rbac.ResourceRBACConfig:  false,
	}
	parentRes := map[rbac.Resource]bool{
		rbac.ResourceProject:     false,
		rbac.ResourceProcess:     true,
		rbac.ResourceTask:        true,
		rbac.ResourceMilestone:   true,
		rbac.ResourceAssignment:  true,
		rbac.ResourceState:       false,
		rbac.ResourceResource:    false,
		rbac.ResourceWorker:      false,
		rbac.ResourceComment:     false,
		rbac.ResourceUserCatalog: false,
		rbac.ResourceRBACConfig:  false,
	}
	ancestorRes := map[rbac.Resource]bool{
		rbac.ResourceProject:     false,
		rbac.ResourceProcess:     true,
		rbac.ResourceTask:        true,
		rbac.ResourceMilestone:   true,
		rbac.ResourceAssignment:  true,
		rbac.ResourceState:       false,
		rbac.ResourceResource:    false,
		rbac.ResourceWorker:      false,
		rbac.ResourceComment:     false,
		rbac.ResourceUserCatalog: false,
		rbac.ResourceRBACConfig:  false,
	}
	for i := rbac.ResourceProject; i <= rbac.ResourceRBACConfig; i++ {
		if !policies.ScopeApplicable(i, policies.ScopeAll) {
			t.Errorf("ScopeApplicable(%v, all) = false; want true", i)
		}
		if policies.ScopeApplicable(i, policies.ScopeNone) {
			t.Errorf("ScopeApplicable(%v, none) = true; want false", i)
		}
		wantOwn := ownRes[i]
		if got := policies.ScopeApplicable(i, policies.ScopeOwn); got != wantOwn {
			t.Errorf("ScopeApplicable(%v, own) = %v; want %v", i, got, wantOwn)
		}
		wantParent := parentRes[i]
		if got := policies.ScopeApplicable(i, policies.ScopeParent); got != wantParent {
			t.Errorf("ScopeApplicable(%v, parent) = %v; want %v", i, got, wantParent)
		}
		wantAncestor := ancestorRes[i]
		if got := policies.ScopeApplicable(i, policies.ScopeAncestor); got != wantAncestor {
			t.Errorf("ScopeApplicable(%v, ancestor) = %v; want %v", i, got, wantAncestor)
		}
	}
}

// NewMatrix: последняя пара (role, resource, action) побеждает.
func TestNewMatrixLastWins(t *testing.T) {
	t.Parallel()
	m := policies.NewMatrix([]policies.MatrixRule{
		{Res: rbac.ResourceTask, Act: policies.ActionView, Role: "rp", Scope: policies.ScopeParent},
		{Res: rbac.ResourceTask, Act: policies.ActionView, Role: "rp", Scope: policies.ScopeAncestor},
	})
	if got := m.ScopeFor("rp", rbac.ResourceTask, policies.ActionView); got != policies.ScopeAncestor {
		t.Errorf("ScopeFor(rp, task, view) = %v; want ancestor (last wins)", got)
	}
	if got := m.ScopeFor("vp", rbac.ResourceTask, policies.ActionView); got != policies.ScopeNone {
		t.Errorf("ScopeFor(vp, task, view) = %v; want none", got)
	}
}

// DefaultMatrixRules ↔ NewMatrix равны встроенной матрице (reset-путь не теряет
// ни одного правила).
func TestDefaultMatrixRoundTrip(t *testing.T) {
	t.Parallel()
	rebuilt := policies.NewMatrix(policies.DefaultMatrixRules())
	roles := []string{"admin", "dp", "rp", "vp", "worker", "ghost"}
	for i := rbac.ResourceProject; i <= rbac.ResourceRBACConfig; i++ {
		for act := policies.ActionView; act <= policies.ActionDelete; act++ {
			for _, role := range roles {
				want := policies.DefaultMatrix().ScopeFor(role, i, act)
				if got := rebuilt.ScopeFor(role, i, act); got != want {
					t.Errorf("roundtrip mismatch role=%s res=%v act=%v: got %v want %v", role, i, act, got, want)
				}
			}
		}
	}
}

// Валидация kind'ов: валидные собираются, невалидные — ошибка.
func TestValidateSpec(t *testing.T) {
	t.Parallel()
	valid := []policies.RouteSpec{
		{Name: "b.list", Kind: "list", Params: map[string]any{"resource": "project", "query_key": "owner_id"}},
		{
			Name:   "b.view",
			Kind:   "entity",
			Params: map[string]any{"resource": "project", "action": "view", "owner": "id"},
		},
		{
			Name:   "b.view.none",
			Kind:   "entity",
			Params: map[string]any{"resource": "state", "action": "view", "owner": "none"},
		},
		{
			Name:   "b.create",
			Kind:   "create",
			Params: map[string]any{"resource": "task", "parent_resource": "process", "parent_from": "process_id"},
		},
		{
			Name:   "b.create.self",
			Kind:   "create",
			Params: map[string]any{"resource": "worker", "owner_key": "manager_id", "default_self": true},
		},
		{Name: "b.owner.match", Kind: "owner_match", Params: map[string]any{
			"resource": "assignment", "primary_resource": "task", "primary_from": "task_id",
			"compare_resource": "resource", "compare_from": "resource_id", "exempt_roles": []any{"admin"},
		}},
		{Name: "b.author.or", Kind: "author_or", Params: map[string]any{
			"author_resource": "comment", "author_id_param": "comment_id",
			"right_resource": "task", "right_action": "update",
		}},
	}
	for _, spec := range valid {
		if err := policies.ValidateSpec(spec); err != nil {
			t.Errorf("ValidateSpec(%s) неожиданная ошибка: %v", spec.Name, err)
		}
	}

	invalid := []struct {
		name string
		spec policies.RouteSpec
	}{
		{"unknown kind", policies.RouteSpec{Name: "x", Kind: "wat", Params: map[string]any{}}},
		{
			"unknown resource",
			policies.RouteSpec{Name: "x", Kind: "entity", Params: map[string]any{"resource": "nope", "action": "view"}},
		},
		{
			"unknown action",
			policies.RouteSpec{
				Name:   "x",
				Kind:   "entity",
				Params: map[string]any{"resource": "project", "action": "edit"},
			},
		},
		{
			"list без query_key",
			policies.RouteSpec{Name: "x", Kind: "list", Params: map[string]any{"resource": "project"}},
		},
		{
			"create: parent_from без parent_resource",
			policies.RouteSpec{
				Name:   "x",
				Kind:   "create",
				Params: map[string]any{"resource": "task", "parent_from": "process_id"},
			},
		},
		{
			"create: ни parent_from ни owner_key",
			policies.RouteSpec{Name: "x", Kind: "create", Params: map[string]any{"resource": "task"}},
		},
		{"owner_match: неизвестный compare", policies.RouteSpec{Name: "x", Kind: "owner_match", Params: map[string]any{
			"resource": "assignment", "primary_resource": "task", "primary_from": "task_id",
			"compare_resource": "nope", "compare_from": "resource_id",
		}}},
		{"author_or: неизвестное право", policies.RouteSpec{Name: "x", Kind: "author_or", Params: map[string]any{
			"author_resource": "comment", "author_id_param": "comment_id",
			"right_resource": "task", "right_action": "edit",
		}}},
		{
			"имя с пробелами",
			policies.RouteSpec{
				Name:   " x ",
				Kind:   "entity",
				Params: map[string]any{"resource": "project", "action": "view"},
			},
		},
		{
			"пустое имя",
			policies.RouteSpec{
				Name:   "",
				Kind:   "entity",
				Params: map[string]any{"resource": "project", "action": "view"},
			},
		},
		{
			"неверный тип параметра",
			policies.RouteSpec{
				Name:   "x",
				Kind:   "create",
				Params: map[string]any{"resource": 42, "owner_key": "owner_id"},
			},
		},
	}
	for _, tc := range invalid {
		if err := policies.ValidateSpec(tc.spec); err == nil {
			t.Errorf("ValidateSpec(%s) должен дать ошибку", tc.name)
		}
	}

	// BuildPolicies: невалидная спецификация валит всю сборку (fail-closed).
	built, err := policies.BuildPolicies([]policies.RouteSpec{valid[0], invalid[1].spec})
	if err == nil || len(built) != 0 {
		t.Errorf("BuildPolicies с битой спецификацией: err=%v built=%d; want error и пусто", err, len(built))
	}
}

// Kinds-справочник содержит все kind'ы и схемы параметров.
func TestKinds(t *testing.T) {
	t.Parallel()
	infos := policies.Kinds()
	if len(infos) != 5 {
		t.Fatalf("Kinds() = %d; want 5", len(infos))
	}
	seen := map[string]bool{}
	for _, k := range infos {
		seen[k.Name] = true
		if len(k.Params) == 0 {
			t.Errorf("kind %s без схемы параметров", k.Name)
		}
	}
	for _, name := range []string{"list", "entity", "create", "owner_match", "author_or"} {
		if !seen[name] {
			t.Errorf("в Kinds() нет kind'а %s", name)
		}
	}
}
