package engine

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

// PublishInput is the DB-derived policy set fed to Publish: the preset matrix
// rows (p allow), the per-user grants (p allow) and revokes (p2 deny) and the
// user -> preset assignments (g).
type PublishInput struct {
	Grants []RuleGrant
	Denies []RuleDeny
	Roles  []RoleAssign
}

// RuleGrant — an allow policy row: sub (preset or user id), canonical
// "<resource>/<action>" object and the ownership zone.
type RuleGrant struct {
	Sub   string
	Key   string
	Scope string
}

// RuleDeny — a deny policy row: sub (user id) and the canonical object.
type RuleDeny struct {
	Sub string
	Key string
}

// RoleAssign — a user -> preset grouping row.
type RoleAssign struct {
	User   string
	Preset string
}

// Snapshot is an immutable, ready-to-evaluate Casbin enforcer. A fresh
// snapshot is built on every rule reload and published atomically; the old one
// stays valid until the swap, so concurrent readers never observe a
// half-built enforcer.
type Snapshot struct {
	e *casbin.Enforcer
}

//nolint:gochecknoglobals // live snapshot; published by the PolicyStore (fallback — built-in defaults)
var published atomic.Pointer[Snapshot]

//nolint:gochecknoglobals // lazy built-in default snapshot (no init())
var (
	defaultOnce sync.Once
	defaultSnap *Snapshot
)

// current returns the published snapshot or the built-in defaults (before the
// first DB load and in tests).
func current() *Snapshot {
	if s := published.Load(); s != nil {
		return s
	}
	defaultOnce.Do(func() {
		snap, err := build(DefaultPublishInput())
		if err != nil {
			panic("engine: invalid built-in default snapshot: " + err.Error())
		}
		defaultSnap = snap
	})
	return defaultSnap
}

// DefaultPublishInput returns the built-in default policy set — the initial
// snapshot and fallback when the DB is unavailable/empty, and the restore
// source for tests.
func DefaultPublishInput() PublishInput {
	var input PublishInput
	for _, r := range DefaultMatrixRules() {
		input.Grants = append(input.Grants, RuleGrant{
			Sub:   r.Role,
			Key:   resActKey(r.Res, r.Act),
			Scope: ScopeName(r.Scope),
		})
	}
	return input
}

// Publish validates the input, builds a fresh snapshot and atomically replaces
// the active one. Any invalid rule cancels the whole build — the caller stays
// on the last valid snapshot (fail-closed).
func Publish(input PublishInput) error {
	snap, err := build(input)
	if err != nil {
		return err
	}
	published.Store(snap)
	return nil
}

// SetMatrix publishes a preset-only snapshot from a matrix (tests and the
// built-in reset source).
func SetMatrix(m Matrix) {
	var input PublishInput
	for _, r := range m.Rules() {
		input.Grants = append(input.Grants, RuleGrant{
			Sub:   r.Role,
			Key:   resActKey(r.Res, r.Act),
			Scope: ScopeName(r.Scope),
		})
	}
	if err := Publish(input); err != nil {
		panic("engine: invalid matrix snapshot: " + err.Error())
	}
}

// build validates and assembles an enforcer from the input.
func build(input PublishInput) (*Snapshot, error) {
	m, err := model.NewModelFromString(modelText())
	if err != nil {
		return nil, fmt.Errorf("casbin model: %w", err)
	}
	e, err := casbin.NewEnforcer(m)
	if err != nil {
		return nil, fmt.Errorf("casbin enforcer: %w", err)
	}
	registerFunctions(e)
	for _, g := range input.Roles {
		if g.User == "" || g.Preset == "" {
			return nil, errors.New("grouping rule: пустое значение user/preset")
		}
		if _, err := e.AddNamedGroupingPolicy("g", g.User, g.Preset); err != nil {
			return nil, fmt.Errorf("grouping rule %s -> %s: %w", g.User, g.Preset, err)
		}
	}
	for _, gr := range input.Grants {
		if err := validateGrant(gr); err != nil {
			return nil, err
		}
		if _, err := e.AddNamedPolicy("p", gr.Sub, gr.Key, gr.Scope, "allow"); err != nil {
			return nil, fmt.Errorf("allow rule (%s, %s, %s): %w", gr.Sub, gr.Key, gr.Scope, err)
		}
	}
	for _, d := range input.Denies {
		if _, _, err := splitKey(d.Key); err != nil {
			return nil, err
		}
		if d.Sub == "" {
			return nil, errors.New("deny rule: пустой subject")
		}
		if _, err := e.AddNamedPolicy("p", d.Sub, d.Key, "", "deny"); err != nil {
			return nil, fmt.Errorf("deny rule (%s, %s): %w", d.Sub, d.Key, err)
		}
	}
	if err := e.BuildRoleLinks(); err != nil {
		return nil, fmt.Errorf("casbin role links: %w", err)
	}
	return &Snapshot{e: e}, nil
}

// validateGrant checks the canonical object and the ownership zone of an allow
// rule (mirrors the matrix rule validation).
func validateGrant(g RuleGrant) error {
	res, act, err := splitKey(g.Key)
	if err != nil {
		return err
	}
	scope, ok := ParseScope(g.Scope)
	if !ok || scope == ScopeNone {
		return fmt.Errorf("недопустимая зона %q (preset=%s, resource=%s, action=%s)", g.Scope, g.Sub, ResourceName(res), ActionName(act))
	}
	if !ScopeApplicable(res, scope) {
		return fmt.Errorf("зона %q неприменима к ресурсу %q (preset=%s)", g.Scope, ResourceName(res), g.Sub)
	}
	if g.Sub == "" {
		return fmt.Errorf("пустой subject (resource=%s, action=%s)", ResourceName(res), ActionName(act))
	}
	return nil
}
