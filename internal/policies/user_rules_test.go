package policies_test

import (
	"testing"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/policies"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
	userctx "github.com/Koshsky/erp-backend/internal/userctx"
)

// TestScopeForUserOverride verifies the per-user resolution rule:
// admin → all; a per-user override wins over the preset; a revoke beats a
// preset grant; no preset → the override decides; nothing → none.
func TestScopeForUserOverride(t *testing.T) {
	t.Parallel()
	m := policies.DefaultMatrix()

	// Preset-only user: rp views own projects.
	rp := userctx.UserContext{Preset: userdomain.PresetProjectManager}
	if got := m.ScopeForUser(rp, rbac.ResourceProject, policies.ActionView); got != policies.ScopeOwn {
		t.Errorf("rp without overrides: got %v; want own", got)
	}

	// Admin bypass wins regardless of rules/overrides.
	adminUser := userctx.UserContext{
		Preset: userdomain.PresetAdmin,
		Admin:  true,
		Rules:  []userctx.PermissionRule{{Resource: "project", Action: "view", Granted: false}},
	}
	if got := m.ScopeForUser(adminUser, rbac.ResourceProject, policies.ActionView); got != policies.ScopeAll {
		t.Errorf("admin bypass: got %v; want all", got)
	}

	// Worker (no preset rights) with an explicit grant override.
	workerGranted := userctx.UserContext{
		Preset: userdomain.PresetWorker,
		Rules: []userctx.PermissionRule{
			{Resource: "project", Action: "view", Scope: "all", Granted: true},
		},
	}
	if got := m.ScopeForUser(workerGranted, rbac.ResourceProject, policies.ActionView); got != policies.ScopeAll {
		t.Errorf("grant override: got %v; want all", got)
	}

	// Revoke beats the preset grant (rp loses project.view).
	rpRevoked := userctx.UserContext{
		Preset: userdomain.PresetProjectManager,
		Rules: []userctx.PermissionRule{
			{Resource: "project", Action: "view", Granted: false},
		},
	}
	if got := m.ScopeForUser(rpRevoked, rbac.ResourceProject, policies.ActionView); got != policies.ScopeNone {
		t.Errorf("revoke override: got %v; want none", got)
	}

	// An unrelated override does not affect the preset rule.
	rpOther := userctx.UserContext{
		Preset: userdomain.PresetProjectManager,
		Rules: []userctx.PermissionRule{
			{Resource: "task", Action: "create", Scope: "parent", Granted: true},
		},
	}
	if got := m.ScopeForUser(rpOther, rbac.ResourceProject, policies.ActionView); got != policies.ScopeOwn {
		t.Errorf("unrelated override: got %v; want own (preset)", got)
	}
}

// TestAuthorizeUserWithOverrides verifies the owner-chain check on top of the
// per-user resolution: an own-scope grant matches the row owner, a revoke
// denies even when the preset would allow it.
func TestAuthorizeUserWithOverrides(t *testing.T) {
	t.Parallel()

	// worker granted project.update (own) — allowed only on own rows.
	workerGranted := userctx.UserContext{
		Preset: userdomain.PresetWorker,
		Rules: []userctx.PermissionRule{
			{Resource: "project", Action: "update", Scope: "own", Granted: true},
		},
	}
	if !policies.AuthorizeUser(
		workerGranted,
		rbac.ResourceProject,
		policies.ActionUpdate,
		rbac.Owners{ProjectOwner: 42},
		42,
	) {
		t.Errorf("own-scope grant: want allowed for the owner")
	}
	if policies.AuthorizeUser(
		workerGranted,
		rbac.ResourceProject,
		policies.ActionUpdate,
		rbac.Owners{ProjectOwner: 42},
		7,
	) {
		t.Errorf("own-scope grant: want denied for a foreign user")
	}

	// rp preset allows own project.delete; the revoke must deny entirely.
	rpRevoked := userctx.UserContext{
		Preset: userdomain.PresetProjectManager,
		Rules: []userctx.PermissionRule{
			{Resource: "project", Action: "delete", Granted: false},
		},
	}
	if policies.AuthorizeUser(
		rpRevoked,
		rbac.ResourceProject,
		policies.ActionDelete,
		rbac.Owners{ProjectOwner: 42},
		42,
	) {
		t.Errorf("revoke override: want denied")
	}
}

// TestViewScopeCodeUser mirrors the preset-based listing scope with overrides.
func TestViewScopeCodeUser(t *testing.T) {
	t.Parallel()

	rp := userctx.UserContext{Preset: userdomain.PresetProjectManager}
	if got := policies.ViewScopeCodeUser(rp, rbac.ResourceProject); got != "own" {
		t.Errorf("rp project view scope: got %q; want own", got)
	}

	rpRevoked := userctx.UserContext{
		Preset: userdomain.PresetProjectManager,
		Rules:  []userctx.PermissionRule{{Resource: "project", Action: "view", Granted: false}},
	}
	if got := policies.ViewScopeCodeUser(
		rpRevoked,
		rbac.ResourceProject,
	); got != "" {
		t.Errorf("revoked view scope: got %q; want empty", got)
	}

	adminUser := userctx.UserContext{Preset: userdomain.PresetAdmin, Admin: true}
	if got := policies.ViewScopeCodeUser(adminUser, rbac.ResourceProject); got != "all" {
		t.Errorf("admin view scope: got %q; want all", got)
	}
}
