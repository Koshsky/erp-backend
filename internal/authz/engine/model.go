// Package engine is the access-control core of the ERP. It merges the former
// internal/policies rules (permission matrix, codecs, route-policy kinds) with
// the Casbin policy engine: RBAC role inheritance (user -> preset), ACL
// allow/deny policy rows (the matrix rules + per-user overrides) and an ABAC
// matcher that compares the caller's attributes (Subject) with the target's
// attributes (Object). The rules come from the Postgres policy tables via the
// internal/authz/service PolicyStore; the defaults here are the
// fallback/reset source.
package engine

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/govaluate"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
)

// modelText is the Casbin model.
//
//   - r.sub / r.obj are ABAC attribute structs (Subject/Object); p.sub is
//     either a preset (inherited through g) or a concrete user id (ACL).
//   - p carries the ownership zone of a (resource, action) grant (scope);
//     p2 carries per-user revokes. The effect allows only when an allow rule
//     matched and no deny rule matched for the same (sub, obj).
//   - The matcher compares the owner attributes of the object with the caller
//     (own/parent/ancestor zones); the special route modes (author_or,
//     owner_match) are evaluated as additional ABAC legs.
//
// modelTmpl is the Casbin model skeleton; %s is substituted with the shared
// ownership-zone sub-expression (scopeLegExpr) for the standard and the
// owner_match legs. Allow rows carry the zone in p.scope, user-level revokes
// are the same policy type with eft "deny" (effect: allow only when an allow
// rule matched and no deny rule matched for the same (sub, obj)).
const modelTmpl = `
[request_definition]
r = sub, obj

[policy_definition]
p = sub, obj, scope, eft

[role_definition]
g = _, _

[policy_effect]
e = some(where (p_eft == allow)) && !some(where (p_eft == deny))

[matchers]
m = (p.eft == "allow" && r.obj.Mode == "author_or" && r.obj.AuthorID != "" && r.obj.AuthorID == r.sub.ID) || (r.obj.Key == p.obj && (p.sub == r.sub.ID || p.sub == r.sub.Role || g(r.sub.ID, p.sub) || g(r.sub.Role, p.sub)) && (p.eft == "deny" || (p.eft == "allow" && ((r.obj.Mode != "owner_match" && %s) || (r.obj.Mode == "owner_match" && eval(sharesOwner(r.obj.OwnerA, r.obj.OwnerB)) && %s)))))
`

// modelText composes the final Casbin model.
func modelText() string {
	return fmt.Sprintf(modelTmpl, scopeLegExpr, scopeLegExpr)
}

// scopeLegExpr is the shared ownership-zone sub-expression of the matcher:
// p.scope published as an attribute of the allow rule, compared against the
// object's owner attributes (own/parent/ancestor).
const scopeLegExpr = `(p.scope == "all" || (p.scope == "own" && r.obj.OwnerID != "" && r.obj.OwnerID == r.sub.ID) || (p.scope == "parent" && r.obj.ParentOwnerID != "" && r.obj.ParentOwnerID == r.sub.ID) || (p.scope == "ancestor" && (r.obj.OwnerID == r.sub.ID || r.obj.ProcessOwnerID == r.sub.ID || r.obj.ProjectOwnerID == r.sub.ID)))`

// Route check modes of the object (special ABAC legs).
const (
	// ModeAuthorOr — the target is allowed when the caller authored it
	// (comment delete: the author always, others by the parent's right).
	ModeAuthorOr = "author_or"
	// ModeOwnerMatch — cross-entity rule: the target's owner chain must match
	// the compared entity (assignment create: the resource must belong to the
	// task's owner).
	ModeOwnerMatch = "owner_match"
)

// Subject is the caller's attribute bundle evaluated by the model (ABAC).
type Subject struct {
	// ID — the caller's user id ("" for pure-preset checks).
	ID string
	// Role — the caller's assigned preset ("" for pure-user checks).
	Role string
	// Admin — the admin bypass (checked in code before Enforce; kept here for
	// the model's defense-in-depth).
	Admin bool
}

// Object is the target's attribute bundle evaluated by the model (ABAC).
type Object struct {
	// Key — the canonical "<resource>/<action>" policy object.
	Key string
	// Mode — "" for standard checks, otherwise a special route mode.
	Mode string
	// OwnerID — the row owner (own zone).
	OwnerID string
	// ParentOwnerID — the immediate parent owner (parent zone).
	ParentOwnerID string
	// ProcessOwnerID / ProjectOwnerID — the ancestor chain (ancestor zone).
	ProcessOwnerID string
	ProjectOwnerID string
	// AuthorID — the author leg of the author_or mode.
	AuthorID string
	// OwnerA / OwnerB — the compared owner chains of the owner_match mode.
	OwnerA rbac.Owners
	OwnerB rbac.Owners
}

// sharesOwnerFunc is the ABAC function backing eval(sharesOwner(...)): the
// cross-entity business rule ("a resource can only be assigned to a task of
// its own owner").
func sharesOwnerFunc(args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return false, nil
	}
	a, okA := args[0].(rbac.Owners)
	b, okB := args[1].(rbac.Owners)
	if !okA || !okB {
		return false, nil
	}
	return a.SharesOwner(b), nil
}

// registerFunctions attaches the ABAC functions to the enforcer.
func registerFunctions(e *casbin.Enforcer) {
	e.AddFunction("sharesOwner", govaluate.ExpressionFunction(sharesOwnerFunc))
}
