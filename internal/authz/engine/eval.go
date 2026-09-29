package engine

// Scope-expression evaluation for a single entity: the caller may access the
// entity when the effective expression is satisfied — the caller owns at least
// one node reachable from the entity by the expression's moves. Owner-chain
// moves (self/upN/up) use the Owners chain; data-dependent moves (sib/down)
// use the optional probes (nil — evaluate to false).

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
)

// OwnerProbe — the data-dependent questions of scope evaluation (need the DB):
//   - descendant — does uid own a node below (res, id) in the tree?
//   - sibling — does uid own a direct sibling node of (res, id) (same parent)?
type OwnerProbe struct {
	Descendant func(ctx context.Context, res rbac.Resource, id, userID int64) (bool, error)
	Sibling    func(ctx context.Context, res rbac.Resource, id, userID int64) (bool, error)
}

// evalOwnerChain — chain-based evaluation of one move (self/up1/up) against
// the entity's owner chain.
func evalChainMove(m moveToken, res rbac.Resource, owners rbac.Owners, uid int64) bool {
	if uid == 0 {
		return false
	}
	switch m.kind {
	case moveSelf:
		return ownField(res, owners) == uid
	case moveUp:
		// up1 — the immediate parent owner; up (any depth) / upN — any chain
		// owner (the project→process→task chain is ≤3 levels, so up2 == up).
		if m.depth == 1 {
			return parentField(res, owners) == uid
		}
		return ancestorMatch(res, owners, uid)
	}
	return false
}

// EvalScope — evaluation of a scope expression against an entity: chain moves
// from Owners, sib/down via the probes (nil probe → the move grants nothing).
func EvalScope(
	ctx context.Context,
	expr string,
	res rbac.Resource,
	owners rbac.Owners,
	uid int64,
	id int64,
	probe *OwnerProbe,
) bool {
	p, ok := parseScope(expr)
	if !ok || p.none {
		return false
	}
	if p.all {
		return true
	}
	for _, m := range p.moves {
		if evalMove(ctx, m, res, owners, uid, id, probe) {
			return true
		}
	}
	return false
}

// evalMove evaluates one move of an expression.
func evalMove(
	ctx context.Context,
	m moveToken,
	res rbac.Resource,
	owners rbac.Owners,
	uid int64,
	id int64,
	probe *OwnerProbe,
) bool {
	switch m.kind {
	case moveSelf, moveUp:
		return evalChainMove(m, res, owners, uid)
	case moveDown:
		if probe != nil && probe.Descendant != nil && id != 0 {
			ok, err := probe.Descendant(ctx, res, id, uid)
			return err == nil && ok
		}
	case moveSib:
		if probe != nil && probe.Sibling != nil && id != 0 {
			ok, err := probe.Sibling(ctx, res, id, uid)
			return err == nil && ok
		}
	}
	return false
}
