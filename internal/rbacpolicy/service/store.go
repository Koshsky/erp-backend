package service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Koshsky/erp-backend/internal/config"
	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/policies"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy/domain"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy/repository"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
	userctx "github.com/Koshsky/erp-backend/internal/userctx"
)

// reloadTimeout — budget of a single rule reload from the DB.
const reloadTimeout = 10 * time.Second

// PolicyStore keeps DB rules in memory and publishes them to the engine:
// the matrix via policies.SetMatrix, route policies via
// rbac.Middleware.Refresh, and the per-user principal snapshot (admin bypass,
// assigned preset, individual overrides) served to the auth middleware via
// EffectiveUser. When the DB is unavailable it runs on the built-in defaults
// and "heals" itself by TTL. No load error brings the service
// down: only a valid, consistent snapshot is applied.
type PolicyStore struct {
	logger   *slog.Logger
	repo     *repository.RuleRepository
	mw       *rbac.Middleware
	interval time.Duration
	stop     chan struct{}

	mu      sync.RWMutex
	users   map[int64]userctx.UserContext
	started bool
}

// NewPolicyStore builds the policy store.
func NewPolicyStore(
	logger *slog.Logger,
	repo *repository.RuleRepository,
	mw *rbac.Middleware,
	interval config.Duration,
) *PolicyStore {
	return &PolicyStore{
		logger:   logger,
		repo:     repo,
		mw:       mw,
		interval: time.Duration(interval),
		stop:     make(chan struct{}),
		users:    map[int64]userctx.UserContext{},
	}
}

// Start performs the initial load (best-effort: when the DB is unavailable —
// a WARN and defaults) and starts the background TTL refresh.
func (s *PolicyStore) Start() {
	s.mu.Lock()
	s.started = true
	s.mu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), reloadTimeout)
	if err := s.Reload(ctx); err != nil {
		s.logger.Warn("rbac: стартовая загрузка правил не удалась, работаю на дефолтах", "error", err)
	}
	cancel()
	go s.refreshLoop()
}

// refreshLoop periodically reloads rules from the DB (eventual
// consistency across instances; local mutations are applied immediately
// via Reload from the service).
func (s *PolicyStore) refreshLoop() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-s.stop:
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), reloadTimeout)
			if err := s.Reload(ctx); err != nil {
				s.logger.Warn(
					"rbac: фоновая перезагрузка правил не удалась, остаюсь на последнем снапшоте",
					"error",
					err,
				)
			}
			cancel()
		}
	}
}

// Stop stops the background refresh (for graceful shutdown).
func (s *PolicyStore) Stop() {
	close(s.stop)
}

// Reload re-reads rules and user principals from the DB and publishes them to
// the engine. On any error (including an empty route policy set) the snapshot
// stays unchanged.
func (s *PolicyStore) Reload(ctx context.Context) error {
	rules, err := s.repo.ListActiveRules(ctx)
	if err != nil {
		return fmt.Errorf("загрузка правил: %w", err)
	}
	matrix, err := rulesToMatrix(rules)
	if err != nil {
		return err
	}

	routePolicies, err := s.repo.ListActiveRoutePolicies(ctx)
	if err != nil {
		return fmt.Errorf("загрузка маршрутных проверок: %w", err)
	}
	if len(routePolicies) == 0 {
		return fmt.Errorf(
			"в БД нет ни одной активной маршрутной проверки — применение отменено (защита от полной блокировки)",
		)
	}
	built, err := policies.BuildPolicies(routePoliciesToSpecs(routePolicies))
	if err != nil {
		return fmt.Errorf("сборка маршрутных проверок: %w", err)
	}

	principals, err := s.loadPrincipals(ctx)
	if err != nil {
		return fmt.Errorf("загрузка прав пользователей: %w", err)
	}

	policies.SetMatrix(matrix)
	s.mw.Refresh(built)
	s.mu.Lock()
	s.users = principals
	s.mu.Unlock()
	return nil
}

// EffectiveUser returns the in-memory principal of a user (admin bypass,
// assigned preset, individual overrides). It never touches the DB per request:
// the snapshot is refreshed by Start/Reload. A missing entry (unknown or
// not-yet-loaded user) returns the default deny principal.
func (s *PolicyStore) EffectiveUser(_ context.Context, userID int64) (userctx.UserContext, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[userID]
	if !ok {
		return userctx.UserContext{ID: userID}, nil
	}
	return u, nil
}

// loadPrincipals builds the user → principal snapshot: every active user with
// the assigned preset (admin preset ⇒ the admin bypass) and the individual
// overrides grouped by user.
func (s *PolicyStore) loadPrincipals(ctx context.Context) (map[int64]userctx.UserContext, error) {
	rows, err := s.repo.ListUserPrincipals(ctx)
	if err != nil {
		return nil, err
	}
	perms, err := s.repo.ListAllUserPermissions(ctx)
	if err != nil {
		return nil, err
	}
	overrides := make(map[int64][]userctx.PermissionRule, len(perms))
	for _, p := range perms {
		overrides[p.UserID] = append(overrides[p.UserID], userctx.PermissionRule{
			Resource: p.Resource,
			Action:   p.Action,
			Scope:    p.Scope,
			Granted:  p.Granted,
		})
	}
	users := make(map[int64]userctx.UserContext, len(rows))
	for _, row := range rows {
		preset := ""
		if row.Preset.Valid {
			preset = row.Preset.String
		}
		users[row.UserID] = userctx.UserContext{
			ID:     row.UserID,
			Admin:  preset == userdomain.PresetAdmin,
			Preset: preset,
			Rules:  overrides[row.UserID],
		}
	}
	return users, nil
}

// IsReady reports whether the store has completed at least one reload attempt
// (used by the auth middleware to fall back to the JWT preset before then).
func (s *PolicyStore) IsReady() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.started
}

// rulesToMatrix converts DB rows into a matrix, validating the codecs.
func rulesToMatrix(rules []domain.PresetRule) (policies.Matrix, error) {
	rows := make([]policies.MatrixRule, 0, len(rules))
	for _, r := range rules {
		res, ok := policies.ParseResource(r.Resource)
		if !ok {
			return policies.Matrix{}, fmt.Errorf(
				"неизвестный ресурс %q в правиле (preset=%s, action=%s)",
				r.Resource,
				r.Preset,
				r.Action,
			)
		}
		act, ok := policies.ParseAction(r.Action)
		if !ok {
			return policies.Matrix{}, fmt.Errorf(
				"неизвестное действие %q в правиле (preset=%s, resource=%s)",
				r.Action,
				r.Preset,
				r.Resource,
			)
		}
		scope, ok := policies.ParseScope(r.Scope)
		if !ok || scope == policies.ScopeNone {
			return policies.Matrix{}, fmt.Errorf(
				"недопустимая зона %q в правиле (preset=%s, resource=%s, action=%s)",
				r.Scope,
				r.Preset,
				r.Resource,
				r.Action,
			)
		}
		if !policies.ScopeApplicable(res, scope) {
			return policies.Matrix{}, fmt.Errorf(
				"зона %q неприменима к ресурсу %q (preset=%s)",
				r.Scope,
				r.Resource,
				r.Preset,
			)
		}
		rows = append(rows, policies.MatrixRule{Res: res, Act: act, Role: r.Preset, Scope: scope})
	}
	return policies.NewMatrix(rows), nil
}

// routePoliciesToSpecs converts route policy definitions into engine
// specifications (kind and parameter validation happens in BuildPolicies).
func routePoliciesToSpecs(routes []domain.RoutePolicy) []policies.RouteSpec {
	specs := make([]policies.RouteSpec, 0, len(routes))
	for _, p := range routes {
		specs = append(specs, policies.RouteSpec{Name: p.Name, Kind: p.Kind, Params: p.Params})
	}
	return specs
}
