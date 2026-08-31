package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Koshsky/erp-backend/internal/config"
	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/policies"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy/domain"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy/repository"
)

// reloadTimeout — бюджет одной перезагрузки правил из БД.
const reloadTimeout = 10 * time.Second

// PolicyStore держит правила из БД в памяти и публикует их в движок:
// матрицу — через policies.SetMatrix, маршрутные проверки — через
// rbac.Middleware.Refresh. При недоступной БД работает на встроенных
// дефолтах и «долечивается» по TTL. Никакая ошибка загрузки не роняет
// сервис: применяется только валидный целостный снапшот.
type PolicyStore struct {
	logger   *slog.Logger
	repo     *repository.RuleRepository
	mw       *rbac.Middleware
	interval time.Duration
	stop     chan struct{}
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
	}
}

// Start выполняет стартовую загрузку (best-effort: при недоступной БД —
// WARN и дефолты) и запускает фоновое TTL-обновление.
func (s *PolicyStore) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), reloadTimeout)
	if err := s.Reload(ctx); err != nil {
		s.logger.Warn("rbac: стартовая загрузка правил не удалась, работаю на дефолтах", "error", err)
	}
	cancel()
	go s.refreshLoop()
}

// refreshLoop периодически перезагружает правила из БД (эвентуальная
// консистентность между инстансами; локальные мутации применяются сразу
// через Reload из сервиса).
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

// Stop останавливает фоновое обновление (для graceful shutdown).
func (s *PolicyStore) Stop() {
	close(s.stop)
}

// Reload перечитывает правила из БД и публикует их в движок. При любой
// ошибке (включая пустой набор маршрутных проверок) снапшот не меняется.
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

	policies.SetMatrix(matrix)
	s.mw.Refresh(built)
	return nil
}

// rulesToMatrix преобразует строки БД в матрицу, валидируя кодеки.
func rulesToMatrix(rules []domain.Rule) (policies.Matrix, error) {
	rows := make([]policies.MatrixRule, 0, len(rules))
	for _, r := range rules {
		res, ok := policies.ParseResource(r.Resource)
		if !ok {
			return policies.Matrix{}, fmt.Errorf(
				"неизвестный ресурс %q в правиле (role=%s, action=%s)",
				r.Resource,
				r.Role,
				r.Action,
			)
		}
		act, ok := policies.ParseAction(r.Action)
		if !ok {
			return policies.Matrix{}, fmt.Errorf(
				"неизвестное действие %q в правиле (role=%s, resource=%s)",
				r.Action,
				r.Role,
				r.Resource,
			)
		}
		scope, ok := policies.ParseScope(r.Scope)
		if !ok || scope == policies.ScopeNone {
			return policies.Matrix{}, fmt.Errorf(
				"недопустимая зона %q в правиле (role=%s, resource=%s, action=%s)",
				r.Scope,
				r.Role,
				r.Resource,
				r.Action,
			)
		}
		if !policies.ScopeApplicable(res, scope) {
			return policies.Matrix{}, fmt.Errorf(
				"зона %q неприменима к ресурсу %q (role=%s)",
				r.Scope,
				r.Resource,
				r.Role,
			)
		}
		rows = append(rows, policies.MatrixRule{Res: res, Act: act, Role: r.Role, Scope: scope})
	}
	return policies.NewMatrix(rows), nil
}

// routePoliciesToSpecs преобразует определения маршрутных проверок в
// спецификации движка (валидация kind и параметров — в BuildPolicies).
func routePoliciesToSpecs(routes []domain.RoutePolicy) []policies.RouteSpec {
	specs := make([]policies.RouteSpec, 0, len(routes))
	for _, p := range routes {
		specs = append(specs, policies.RouteSpec{Name: p.Name, Kind: p.Kind, Params: p.Params})
	}
	return specs
}
