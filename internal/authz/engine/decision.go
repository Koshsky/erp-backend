// Runtime access decisions over the Casbin snapshot: scope resolution for the
// caller (preset + per-user overrides + admin bypass), entity authorization
// (the ABAC matcher) and the route-check entry points (EnforceTarget).
package engine

import (
	"strconv"
	"strings"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
	userctx "github.com/Koshsky/erp-backend/internal/userctx"
)

// resActKey returns the canonical "<resource>/<action>" policy object.
func resActKey(res rbac.Resource, act Action) string {
	return ResourceName(res) + "/" + ActionName(act)
}

// splitKey parses a canonical "<resource>/<action>" object.
func splitKey(key string) (rbac.Resource, Action, error) {
	resName, actName, found := strings.Cut(key, "/")
	if !found {
		return 0, 0, errInvalidKey(key)
	}
	res, okRes := ParseResource(resName)
	if !okRes {
		return 0, 0, errInvalidKey(key)
	}
	act, okAct := ParseAction(actName)
	if !okAct {
		return 0, 0, errInvalidKey(key)
	}
	return res, act, nil
}

func errInvalidKey(key string) error {
	return &invalidKeyError{key: key}
}

// invalidKeyError — a malformed or unknown "<resource>/<action>" object.
type invalidKeyError struct {
	key string
}

func (e *invalidKeyError) Error() string {
	return "неизвестный объект доступа " + strconv.Quote(e.key)
}

// numericSub reports whether a policy subject is a user id (the ACL rows) as
// opposed to a preset name (the matrix rows).
func numericSub(sub string) bool {
	if sub == "" {
		return false
	}
	_, err := strconv.ParseInt(sub, 10, 64)
	return err == nil
}

// enforce runs the ABAC matcher for a subject and object; the admin bypass is
// a code-level invariant (deny rows must not affect admins).
func (s *Snapshot) enforce(sub Subject, obj Object) bool {
	if sub.Admin {
		return true
	}
	allowed, err := s.e.Enforce(sub, obj)
	if err != nil {
		return false
	}
	return allowed
}

// CurrentMatrix returns the effective preset matrix assembled from the Casbin
// allow rows (the matrix/explain API source; admin bypass is not stored).
func CurrentMatrix() Matrix {
	rules := make([]MatrixRule, 0, 64)
	rows, err := current().e.GetPolicy()
	if err != nil {
		return NewMatrix(rules)
	}
	for _, row := range rows {
		if len(row) < 4 || row[3] != "allow" || numericSub(row[0]) {
			continue
		}
		res, act, keyErr := splitKey(row[1])
		if keyErr != nil {
			continue
		}
		scope, ok := ParseScope(row[2])
		if !ok || scope == ScopeNone {
			continue
		}
		rules = append(rules, MatrixRule{Res: res, Act: act, Role: row[0], Scope: scope})
	}
	return NewMatrix(rules)
}

// scopeFor resolves the preset zone (admin gets ScopeAll — invariant).
func scopeFor(role string, res rbac.Resource, act Action) Scope {
	if role == userdomain.PresetAdmin {
		return ScopeAll
	}
	rows, err := current().e.GetFilteredPolicy(0, role, resActKey(res, act))
	if err != nil || len(rows) == 0 {
		return ScopeNone
	}
	scope, ok := ParseScope(rows[0][2])
	if !ok {
		return ScopeNone
	}
	return scope
}

// scopeForUser resolves the caller's effective zone: admin — ScopeAll; a
// user-level deny (p2) — ScopeNone; a user-level grant — its scope; otherwise
// the preset zone through the role hierarchy (g).
func scopeForUser(u userctx.UserContext, res rbac.Resource, act Action) Scope {
	if u.Admin {
		return ScopeAll
	}
	s := current()
	key := resActKey(res, act)
	uid := strconv.FormatInt(u.ID, 10)

	// A user-level revoke (eft deny on the same (sub, obj)) shadows both the
	// user grants and the preset rows.
	if rows, err := s.e.GetFilteredPolicy(0, uid, key); err == nil {
		for _, row := range rows {
			if len(row) >= 4 && row[3] == "deny" {
				return ScopeNone
			}
		}
	}
	perms, err := s.e.GetImplicitPermissionsForUser(uid)
	if err != nil {
		return ScopeNone
	}
	fallback := ScopeNone
	for _, p := range perms {
		if len(p) < 4 || p[1] != key || p[3] != "allow" {
			continue
		}
		scope, ok := ParseScope(p[2])
		if !ok {
			continue
		}
		if p[0] == uid {
			return scope // a user-level grant wins over the preset
		}
		fallback = scope
	}
	return fallback
}

// ScopeForUser returns the caller's effective zone (admin bypass + per-user
// overrides applied; Casbin snapshot).
func ScopeForUser(u userctx.UserContext, res rbac.Resource, act Action) Scope {
	return scopeForUser(u, res, act)
}

// Can reports whether a preset can perform an action at all.
func Can(role string, res rbac.Resource, act Action) bool {
	return scopeFor(role, res, act) != ScopeNone
}

// CanUser reports whether the caller can perform an action at all.
func CanUser(u userctx.UserContext, res rbac.Resource, act Action) bool {
	return scopeForUser(u, res, act) != ScopeNone
}

// ViewScopeCode returns the string code of the preset view zone for listing
// requests (all|own|parent|ancestor).
func ViewScopeCode(role string, res rbac.Resource) string {
	return ScopeName(scopeFor(role, res, ActionView))
}

// ViewScopeCodeUser returns the string code of the caller's view zone for
// listing requests (admin bypass + per-user overrides applied).
func ViewScopeCodeUser(u userctx.UserContext, res rbac.Resource) string {
	return ScopeName(scopeForUser(u, res, ActionView))
}

// Authorize reports whether a preset may perform an action on an entity with
// its owners (the preset introspection path: matrix scope + owner comparison;
// per-user overrides do not apply — mirror the matcher's scope leg).
func Authorize(role string, res rbac.Resource, act Action, owners rbac.Owners, userID int64) bool {
	if role == userdomain.PresetAdmin {
		return true
	}
	return authorizeScope(scopeFor(role, res, act), res, owners, userID)
}

// AuthorizeUser reports whether the caller may perform an action on an entity
// with its owners (admin bypass + per-user overrides applied; the ABAC matcher
// evaluates the ownership zone).
func AuthorizeUser(u userctx.UserContext, res rbac.Resource, act Action, owners rbac.Owners, userID int64) bool {
	return current().enforce(
		Subject{ID: itoa(userID), Role: u.Preset, Admin: u.Admin},
		buildObject(res, act, owners, "", 0, rbac.Owners{}),
	)
}

// authorizeScope — the owner-chain mechanism for a resolved zone (the preset
// introspection mirror of the matcher's scope leg).
func authorizeScope(scope Scope, res rbac.Resource, owners rbac.Owners, userID int64) bool {
	switch scope {
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
		rbac.ResourceUserAdmin, rbac.ResourceStateAdmin, rbac.ResourceOrgStructure,
		rbac.ResourceAudit:
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
		rbac.ResourceUserAdmin, rbac.ResourceStateAdmin, rbac.ResourceOrgStructure,
		rbac.ResourceAudit:
		return 0
	}
	return 0
}

// ancestorMatch reports whether the user matches any owner of the entity's
// ownership chain (the L0 row owner or any higher one). For process/milestone
// the self-owner is absent (Owners.Owner = 0) — then the process and project
// owners are considered.
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
		rbac.ResourceUserAdmin, rbac.ResourceStateAdmin, rbac.ResourceOrgStructure,
		rbac.ResourceAudit:
		return false
	}
	return false
}

// Target — a route-check target assembled by the kind builders: the implied
// (resource, action), the resolved owner chain and the special route modes.
type Target struct {
	Res    rbac.Resource
	Act    Action
	Owners rbac.Owners
	Mode   string
	// AuthorID — the author leg of ModeAuthorOr.
	AuthorID int64
	// OwnerB — the compared owner chain of ModeOwnerMatch.
	OwnerB rbac.Owners
}

// EnforceTarget runs the ABAC matcher for a route-check target.
func EnforceTarget(u userctx.UserContext, t Target) bool {
	obj := buildObject(t.Res, t.Act, t.Owners, t.Mode, t.AuthorID, t.OwnerB)
	return current().enforce(Subject{ID: itoa(u.ID), Role: u.Preset, Admin: u.Admin}, obj)
}

// buildObject maps an owner chain into the object's ABAC attributes by the
// resource's zone semantics (mirror ownField/parentField; the ancestor zone
// uses the full chain).
func buildObject(res rbac.Resource, act Action, owners rbac.Owners, mode string, authorID int64, ownerB rbac.Owners) Object {
	return Object{
		Key:            resActKey(res, act),
		Mode:           mode,
		OwnerID:        itoa(ownField(res, owners)),
		ParentOwnerID:  itoa(parentField(res, owners)),
		ProcessOwnerID: itoa(owners.ProcessOwner),
		ProjectOwnerID: itoa(owners.ProjectOwner),
		AuthorID:       itoa(authorID),
		OwnerA:         owners,
		OwnerB:         ownerB,
	}
}

// itoa returns the decimal string of a user/owner id ("" — 0: not set).
func itoa(id int64) string {
	if id == 0 {
		return ""
	}
	return strconv.FormatInt(id, 10)
}
