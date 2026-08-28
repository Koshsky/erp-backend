// Package policies — the single place holding access rules. The engine (rbac)
// stays a pure mechanism; rules are split into a matrix (role × resource ×
// action → ownership scope) and route policies (kind + parameters).
package policies

import (
	"sync/atomic"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
)

// Action — an operation on an entity.
type Action int

const (
	ActionView Action = iota
	ActionCreate
	ActionUpdate
	ActionDelete
)

// Scope — ownership zone required for an action. The mechanism is single: "owner of
// the parent element". Absence of a rule = ScopeNone (no access).
type Scope int

const (
	ScopeNone Scope = iota
	ScopeAll
	ScopeOwn      // owner of the row itself (project for rp; resource/worker for vp)
	ScopeParent   // owner of the immediate parent (managing one level down)
	ScopeAncestor // owner of the ancestor one level above the parent (viewing tasks through the project)
)

// String resource codes (mirror V15 and the kind schemas).
const (
	resProject      = "project"
	resProcess      = "process"
	resTask         = "task"
	resMilestone    = "milestone"
	resAssignment   = "assignment"
	resState        = "state"
	resResource     = "resource"
	resWorker       = "worker"
	resComment      = "comment"
	resUserCatalog  = "user_catalog"
	resRBACConfig   = "rbac_config"
	resUserAdmin    = "user_admin"
	resStateAdmin   = "state_admin"
	resOrgStructure = "org_structure"
)

// String action codes.
const (
	actView   = "view"
	actCreate = "create"
	actUpdate = "update"
	actDelete = "delete"
)

// String ownership zone codes.
const (
	scopeAll      = "all"
	scopeOwn      = "own"
	scopeParent   = "parent"
	scopeAncestor = "ancestor"
)

// Rule binds a role and the required access zone.
type Rule struct {
	Role  string
	Scope Scope
}

// MatrixRule — a matrix row (for building from the DB and defaults).
type MatrixRule struct {
	Res   rbac.Resource
	Act   Action
	Role  string
	Scope Scope
}

// CurrentMatrix returns the active matrix (for the matrix/explain API).
func CurrentMatrix() Matrix { return snapshot() }

// DefaultMatrixRules returns the built-in matrix rules (for reset).
func DefaultMatrixRules() []MatrixRule {
	var out []MatrixRule
	for res, byAction := range defaultMatrix.rules {
		for act, rules := range byAction {
			for _, r := range rules {
				out = append(out, MatrixRule{Res: res, Act: act, Role: r.Role, Scope: r.Scope})
			}
		}
	}
	return out
}

// Matrix — a snapshot of the permission matrix (role × resource × action → rules).
type Matrix struct {
	rules map[rbac.Resource]map[Action][]Rule
}

// NewMatrix assembles a matrix from rules. (role, resource, action) pairs
// are unique: the last rule wins.
func NewMatrix(rules []MatrixRule) Matrix {
	m := Matrix{rules: make(map[rbac.Resource]map[Action][]Rule)}
	for _, r := range rules {
		byAction, ok := m.rules[r.Res]
		if !ok {
			byAction = make(map[Action][]Rule)
			m.rules[r.Res] = byAction
		}
		replaced := false
		for i, existing := range byAction[r.Act] {
			if existing.Role == r.Role {
				byAction[r.Act][i] = Rule{Role: r.Role, Scope: r.Scope}
				replaced = true
				break
			}
		}
		if !replaced {
			byAction[r.Act] = append(byAction[r.Act], Rule{Role: r.Role, Scope: r.Scope})
		}
	}
	return m
}

//nolint:gochecknoglobals // live snapshot; updated by PolicyStore from the DB, fallback — DefaultMatrix
var currentRules atomic.Pointer[Matrix]

// SetMatrix atomically replaces the active permission matrix.
func SetMatrix(m Matrix) {
	currentRules.Store(&m)
}

// snapshot returns the active matrix or the built-in defaults (before the first
// DB load and in tests).
func snapshot() Matrix {
	if m := currentRules.Load(); m != nil {
		return *m
	}
	return DefaultMatrix()
}

// DefaultMatrix — the built-in default matrix: serialization of the seed
// V10__rbac_policies.sql. Used as a fallback, a reset source, and the golden
// equivalence test. admin and worker are not listed explicitly: admin — ScopeAll
// (invariant), worker — ScopeNone (no rows).
//
//nolint:gochecknoglobals // rule registry
var defaultMatrix = Matrix{rules: map[rbac.Resource]map[Action][]Rule{
	rbac.ResourceProject: {
		ActionView: {
			{userdomain.ProjectDirector, ScopeAll},
			{userdomain.ProjectManager, ScopeOwn},
		},
		ActionCreate: {
			// rp creates a project into own ownership (the owner defaults to themselves).
			{userdomain.ProjectManager, ScopeOwn},
		},
		ActionUpdate: {
			// dp and admin edit any project, rp — their own.
			{userdomain.ProjectDirector, ScopeAll},
			{userdomain.ProjectManager, ScopeOwn},
		},
		ActionDelete: {
			{userdomain.ProjectManager, ScopeOwn},
		},
	},
	rbac.ResourceProcess: {
		ActionView: {
			{userdomain.ProjectDirector, ScopeAll},
			// rp — processes of own projects (parent), vp — all for reference.
			{userdomain.ProjectManager, ScopeParent},
			{userdomain.ProcessOwner, ScopeAll},
		},
		ActionCreate: {
			{userdomain.ProjectManager, ScopeParent},
		},
		ActionUpdate: {
			{userdomain.ProjectManager, ScopeParent},
		},
		ActionDelete: {
			{userdomain.ProjectManager, ScopeParent},
		},
	},
	rbac.ResourceTask: {
		ActionView: {
			{userdomain.ProjectDirector, ScopeAll},
			{userdomain.ProjectManager, ScopeAncestor},
			{userdomain.ProcessOwner, ScopeParent},
		},
		ActionCreate: {
			{userdomain.ProcessOwner, ScopeParent},
		},
		ActionUpdate: {
			{userdomain.ProcessOwner, ScopeParent},
		},
		ActionDelete: {
			{userdomain.ProcessOwner, ScopeParent},
		},
	},
	rbac.ResourceMilestone: {
		ActionView: {
			{userdomain.ProjectDirector, ScopeAll},
			{userdomain.ProjectManager, ScopeAncestor},
			{userdomain.ProcessOwner, ScopeParent},
		},
		ActionCreate: {
			{userdomain.ProcessOwner, ScopeParent},
		},
		ActionUpdate: {
			{userdomain.ProcessOwner, ScopeParent},
		},
		ActionDelete: {
			{userdomain.ProcessOwner, ScopeParent},
		},
	},
	rbac.ResourceAssignment: {
		ActionView: {
			{userdomain.ProjectDirector, ScopeAll},
			{userdomain.ProjectManager, ScopeAncestor},
			{userdomain.ProcessOwner, ScopeParent},
		},
		ActionCreate: {
			{userdomain.ProcessOwner, ScopeParent},
		},
		ActionUpdate: {
			{userdomain.ProcessOwner, ScopeParent},
		},
		ActionDelete: {
			{userdomain.ProcessOwner, ScopeParent},
		},
	},
	// === Timesheet ===
	// States: an ownerless reference; vp sees them (for the timesheet), only admin manages them.
	rbac.ResourceState: {
		ActionView: {
			{userdomain.ProcessOwner, ScopeAll},
		},
		ActionCreate: {},
		ActionUpdate: {},
		ActionDelete: {},
	},
	// Resource categories: admin — all, vp — own (own); vp creates into own ownership.
	rbac.ResourceResource: {
		ActionView: {
			{userdomain.ProcessOwner, ScopeOwn},
		},
		ActionCreate: {
			{userdomain.ProcessOwner, ScopeOwn},
		},
		ActionUpdate: {
			{userdomain.ProcessOwner, ScopeOwn},
		},
		ActionDelete: {
			{userdomain.ProcessOwner, ScopeOwn},
		},
	},
	// Workers: creating employees — admin only (bypass); vp — own
	// subordinates (manager_id): view and edit.
	rbac.ResourceWorker: {
		ActionView: {
			{userdomain.ProcessOwner, ScopeOwn},
		},
		ActionUpdate: {
			{userdomain.ProcessOwner, ScopeOwn},
		},
		ActionDelete: {
			{userdomain.ProcessOwner, ScopeOwn},
		},
	},
	// Comments: no access in the common matrix — rights are derived from the parent
	// task (see task.comment.* route checks: list/create by task.view,
	// delete — the author or the task update right).
	rbac.ResourceComment: {},
	// Virtual resources: user_catalog — the user catalog for pickers
	// (dp/rp/vp + admin); rbac_config — auto-create/RBAC administration (admin
	// only, bypass) — no rows in the matrix.
	rbac.ResourceUserCatalog: {
		ActionView: {
			{userdomain.ProjectDirector, ScopeAll},
			{userdomain.ProjectManager, ScopeAll},
			{userdomain.ProcessOwner, ScopeAll},
		},
	},
	rbac.ResourceRBACConfig:   {},
	rbac.ResourceUserAdmin:    {},
	rbac.ResourceStateAdmin:   {},
	rbac.ResourceOrgStructure: {},
}}

// DefaultMatrix returns the built-in default matrix (a copy).
func DefaultMatrix() Matrix {
	return defaultMatrix
}

// ScopeFor returns the required access zone for (role, resource, action).
// admin gets ScopeAll (a protective invariant, not stored in the DB).
func (m Matrix) ScopeFor(role string, res rbac.Resource, act Action) Scope {
	if role == userdomain.Admin {
		return ScopeAll
	}
	rules, ok := m.rules[res][act]
	if !ok {
		return ScopeNone
	}
	for _, r := range rules {
		if r.Role == role {
			return r.Scope
		}
	}
	return ScopeNone
}

// ownField returns the owner of the row itself (chain L0) for a resource
// (0 — the entity has no own owner: own is not applicable).
func ownField(res rbac.Resource, owners rbac.Owners) int64 {
	switch res {
	case rbac.ResourceProject:
		return owners.ProjectOwner
	case rbac.ResourceProcess:
		return owners.ProcessOwner
	case rbac.ResourceTask, rbac.ResourceResource, rbac.ResourceWorker:
		return owners.Owner
	case rbac.ResourceMilestone, rbac.ResourceAssignment, rbac.ResourceState,
		rbac.ResourceComment, rbac.ResourceUserCatalog, rbac.ResourceRBACConfig,
		rbac.ResourceUserAdmin, rbac.ResourceStateAdmin, rbac.ResourceOrgStructure:
		return 0
	}
	return 0
}

// parentField returns the immediate parent owner for a resource
// (0 — the resource has no parent in the project → process → task/… hierarchy).
func parentField(res rbac.Resource, owners rbac.Owners) int64 {
	switch res {
	case rbac.ResourceProcess:
		return owners.ProjectOwner
	case rbac.ResourceTask, rbac.ResourceMilestone, rbac.ResourceAssignment:
		return owners.ProcessOwner
	case rbac.ResourceProject, rbac.ResourceState, rbac.ResourceResource,
		rbac.ResourceWorker, rbac.ResourceComment,
		rbac.ResourceUserCatalog, rbac.ResourceRBACConfig,
		rbac.ResourceUserAdmin, rbac.ResourceStateAdmin, rbac.ResourceOrgStructure:
		return 0
	}
	return 0
}

// ancestorMatch reports whether the user matches any owner of the entity's
// ownership chain (the L0 row owner or any higher one).
// For process/milestone the self-owner is absent (Owners.Owner = 0) — then
// the process and project owners are considered.
func ancestorMatch(res rbac.Resource, owners rbac.Owners, userID int64) bool {
	if userID == 0 {
		return false
	}
	switch res {
	case rbac.ResourceTask, rbac.ResourceMilestone, rbac.ResourceAssignment,
		rbac.ResourceProcess:
		return owners.Owner == userID || owners.ProcessOwner == userID || owners.ProjectOwner == userID
	case rbac.ResourceProject, rbac.ResourceState, rbac.ResourceResource,
		rbac.ResourceWorker, rbac.ResourceComment,
		rbac.ResourceUserCatalog, rbac.ResourceRBACConfig,
		rbac.ResourceUserAdmin, rbac.ResourceStateAdmin, rbac.ResourceOrgStructure:
		return false
	}
	return false
}

// Authorize reports whether the user may perform an action on an entity
// with its owners.
func Authorize(role string, res rbac.Resource, act Action, owners rbac.Owners, userID int64) bool {
	switch snapshot().ScopeFor(role, res, act) {
	case ScopeNone:
		return false
	case ScopeAll:
		return true
	case ScopeOwn:
		owner := ownField(res, owners)
		return userID != 0 && owner != 0 && owner == userID
	case ScopeParent:
		parent := parentField(res, owners)
		return userID != 0 && parent != 0 && parent == userID
	case ScopeAncestor:
		return ancestorMatch(res, owners, userID)
	default:
		return false
	}
}

// Can reports whether a role can perform an action at all
// (a coarse check before loading lists).
func Can(role string, res rbac.Resource, act Action) bool {
	return scopeFor(role, res, act) != ScopeNone
}

// ViewScopeCode returns the string code of the view zone for listing requests
// (all|own|parent|ancestor). SQL applies exactly this zone to the owner chain.
func ViewScopeCode(role string, res rbac.Resource) string {
	return ScopeName(scopeFor(role, res, ActionView))
}

// scopeFor — internal wrapper for the check builders.
func scopeFor(role string, res rbac.Resource, act Action) Scope {
	return snapshot().ScopeFor(role, res, act)
}

//nolint:gochecknoglobals // resource codex (stable dictionary, mirrors V15)
var resourceNames = map[rbac.Resource]string{
	rbac.ResourceProject:      resProject,
	rbac.ResourceProcess:      resProcess,
	rbac.ResourceTask:         resTask,
	rbac.ResourceMilestone:    resMilestone,
	rbac.ResourceAssignment:   resAssignment,
	rbac.ResourceState:        resState,
	rbac.ResourceResource:     resResource,
	rbac.ResourceWorker:       resWorker,
	rbac.ResourceComment:      resComment,
	rbac.ResourceUserCatalog:  resUserCatalog,
	rbac.ResourceRBACConfig:   resRBACConfig,
	rbac.ResourceUserAdmin:    resUserAdmin,
	rbac.ResourceStateAdmin:   resStateAdmin,
	rbac.ResourceOrgStructure: resOrgStructure,
}

//nolint:gochecknoglobals // action codex
var actionNames = map[Action]string{
	ActionView:   actView,
	ActionCreate: actCreate,
	ActionUpdate: actUpdate,
	ActionDelete: actDelete,
}

//nolint:gochecknoglobals // scope codex ("none" is not stored: absence of a row = no access)
var scopeNames = map[Scope]string{
	ScopeNone:     "",
	ScopeAll:      scopeAll,
	ScopeOwn:      scopeOwn,
	ScopeParent:   scopeParent,
	ScopeAncestor: scopeAncestor,
}

// ResourceName returns the string resource code ("" — unknown).
func ResourceName(res rbac.Resource) string { return resourceNames[res] }

// ParseResource parses a string resource code.
func ParseResource(s string) (rbac.Resource, bool) {
	for res, name := range resourceNames {
		if name == s {
			return res, true
		}
	}
	return 0, false
}

// ActionName returns the string action code ("" — unknown).
func ActionName(act Action) string { return actionNames[act] }

// ParseAction parses a string action code.
func ParseAction(s string) (Action, bool) {
	for act, name := range actionNames {
		if name == s {
			return act, true
		}
	}
	return 0, false
}

// ScopeName returns the string zone code ("" — no access, not stored).
func ScopeName(scope Scope) string { return scopeNames[scope] }

// ParseScope parses a string zone code.
func ParseScope(s string) (Scope, bool) {
	for scope, name := range scopeNames {
		if name == s {
			return scope, true
		}
	}
	return 0, false
}

//nolint:gochecknoglobals // scope applicability maps (complete: every resource listed)
var ownApplicable = map[rbac.Resource]bool{
	rbac.ResourceProject:      true,
	rbac.ResourceProcess:      true,
	rbac.ResourceTask:         true,
	rbac.ResourceMilestone:    false,
	rbac.ResourceAssignment:   false,
	rbac.ResourceState:        false,
	rbac.ResourceResource:     true,
	rbac.ResourceWorker:       true,
	rbac.ResourceComment:      false,
	rbac.ResourceUserCatalog:  false,
	rbac.ResourceRBACConfig:   false,
	rbac.ResourceUserAdmin:    false,
	rbac.ResourceStateAdmin:   false,
	rbac.ResourceOrgStructure: false,
}

//nolint:gochecknoglobals // scope applicability maps (complete: every resource listed)
var parentApplicable = map[rbac.Resource]bool{
	rbac.ResourceProject:      false,
	rbac.ResourceProcess:      true,
	rbac.ResourceTask:         true,
	rbac.ResourceMilestone:    true,
	rbac.ResourceAssignment:   true,
	rbac.ResourceState:        false,
	rbac.ResourceResource:     false,
	rbac.ResourceWorker:       false,
	rbac.ResourceComment:      false,
	rbac.ResourceUserCatalog:  false,
	rbac.ResourceRBACConfig:   false,
	rbac.ResourceUserAdmin:    false,
	rbac.ResourceStateAdmin:   false,
	rbac.ResourceOrgStructure: false,
}

//nolint:gochecknoglobals // scope applicability maps (complete: every resource listed)
var ancestorApplicable = map[rbac.Resource]bool{
	rbac.ResourceProject:      false,
	rbac.ResourceProcess:      true,
	rbac.ResourceTask:         true,
	rbac.ResourceMilestone:    true,
	rbac.ResourceAssignment:   true,
	rbac.ResourceState:        false,
	rbac.ResourceResource:     false,
	rbac.ResourceWorker:       false,
	rbac.ResourceComment:      false,
	rbac.ResourceUserCatalog:  false,
	rbac.ResourceRBACConfig:   false,
	rbac.ResourceUserAdmin:    false,
	rbac.ResourceStateAdmin:   false,
	rbac.ResourceOrgStructure: false,
}

// ScopeApplicable reports whether a zone is applicable to a resource (for rule validation).
func ScopeApplicable(res rbac.Resource, scope Scope) bool {
	switch scope {
	case ScopeAll:
		return true
	case ScopeOwn:
		return ownApplicable[res]
	case ScopeParent:
		return parentApplicable[res]
	case ScopeAncestor:
		return ancestorApplicable[res]
	case ScopeNone:
		return false
	}
	return false
}
