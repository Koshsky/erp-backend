package service

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/policies"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy/domain"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy/dto"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy/repository"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

// Service — администрирование RBAC-политик (валидация + запись в БД +
// немедленное применение через PolicyStore).
type Service struct {
	logger *slog.Logger
	repo   *repository.RuleRepository
	store  *PolicyStore
}

// NewRBACService builds the RBAC administration service.
func NewRBACService(logger *slog.Logger, repo *repository.RuleRepository, store *PolicyStore) *Service {
	return &Service{logger: logger, repo: repo, store: store}
}

// ListRoles returns the role catalog.
func (s *Service) ListRoles(ctx context.Context) ([]domain.Role, error) {
	return s.repo.ListActiveRoles(ctx)
}

// ListRules returns the active matrix rows.
func (s *Service) ListRules(ctx context.Context) ([]domain.Rule, error) {
	return s.repo.ListActiveRules(ctx)
}

// UpsertRule валидирует и записывает строку матрицы, затем применяет
// изменения немедленно.
func (s *Service) UpsertRule(ctx context.Context, in dto.RuleInput, updatedBy int64) error {
	if in.Role == "" {
		return errors.BadRequest("role обязательна")
	}
	roles, err := s.repo.ListActiveRoles(ctx)
	if err != nil {
		return err
	}
	if !roleExists(roles, in.Role) {
		return errors.BadRequest("неизвестная роль " + in.Role)
	}
	res, ok := policies.ParseResource(in.Resource)
	if !ok {
		return errors.BadRequest("неизвестный ресурс " + in.Resource)
	}
	if _, okAction := policies.ParseAction(in.Action); !okAction {
		return errors.BadRequest("неизвестное действие " + in.Action)
	}
	scope, ok := policies.ParseScope(in.Scope)
	if !ok || scope == policies.ScopeNone {
		return errors.BadRequest("недопустимая зона " + in.Scope + " (all|own|parent|ancestor)")
	}
	if !policies.ScopeApplicable(res, scope) {
		return errors.BadRequest("зона " + in.Scope + " неприменима к ресурсу " + in.Resource)
	}

	if _, upsertErr := s.repo.UpsertRule(ctx, domain.Rule{
		Role: in.Role, Resource: in.Resource, Action: in.Action, Scope: in.Scope,
		UpdatedBy: &updatedBy,
	}); upsertErr != nil {
		return upsertErr
	}
	return s.apply(ctx)
}

// DeleteRule мягко удаляет строку матрицы.
func (s *Service) DeleteRule(ctx context.Context, id int64) error {
	if err := s.repo.SoftDeleteRule(ctx, id); err != nil {
		return err
	}
	return s.apply(ctx)
}

// ListRoutePolicies returns the active route policy definitions.
func (s *Service) ListRoutePolicies(ctx context.Context) ([]domain.RoutePolicy, error) {
	return s.repo.ListActiveRoutePolicies(ctx)
}

// UpsertRoutePolicy валидирует (kind + параметры по схеме) и записывает
// маршрутную проверку.
func (s *Service) UpsertRoutePolicy(ctx context.Context, in dto.RoutePolicyInput, updatedBy int64) error {
	if err := policies.ValidateSpec(policies.RouteSpec{Name: in.Name, Kind: in.Kind, Params: in.Params}); err != nil {
		return errors.BadRequest(err.Error())
	}
	active := true
	if in.Active != nil {
		active = *in.Active
	}
	if _, err := s.repo.UpsertRoutePolicy(ctx, domain.RoutePolicy{
		Name: in.Name, Kind: in.Kind, Params: in.Params, Active: active,
		UpdatedBy: &updatedBy,
	}); err != nil {
		return err
	}
	return s.apply(ctx)
}

// DeleteRoutePolicy мягко удаляет маршрутную проверку по имени.
func (s *Service) DeleteRoutePolicy(ctx context.Context, name string) error {
	if err := s.repo.SoftDeleteRoutePolicy(ctx, name); err != nil {
		return err
	}
	return s.apply(ctx)
}

// Reset возвращает правила и маршрутные проверки к встроенным дефолтам
// (запасной люк после ошибочных правок).
func (s *Service) Reset(ctx context.Context, updatedBy int64) error {
	if err := s.repo.SoftDeleteAllRules(ctx); err != nil {
		return err
	}
	if err := s.repo.SoftDeleteAllRoutePolicies(ctx); err != nil {
		return err
	}
	for _, r := range policies.DefaultMatrixRules() {
		scope := policies.ScopeName(r.Scope)
		if _, err := s.repo.UpsertRule(ctx, domain.Rule{
			Role: r.Role, Resource: policies.ResourceName(r.Res), Action: policies.ActionName(r.Act),
			Scope: scope, UpdatedBy: &updatedBy,
		}); err != nil {
			return err
		}
	}
	for _, spec := range policies.DefaultRouteSpecs() {
		if _, err := s.repo.UpsertRoutePolicy(ctx, domain.RoutePolicy{
			Name: spec.Name, Kind: spec.Kind, Params: spec.Params, Active: true,
			UpdatedBy: &updatedBy,
		}); err != nil {
			return err
		}
	}
	return s.apply(ctx)
}

// EffectiveMatrix возвращает эффективную матрицу (с admin-байпасом) для API.
func (s *Service) EffectiveMatrix(ctx context.Context) ([]dto.MatrixCell, error) {
	_ = ctx // матрица и роли читаются из памяти/БД без доп. запросов
	roles, err := s.repo.ListActiveRoles(ctx)
	if err != nil {
		return nil, err
	}
	matrix := policies.CurrentMatrix()
	names := make([]string, 0, len(roles)+1)
	names = append(names, "admin")
	for _, r := range roles {
		names = append(names, r.Name)
	}
	var cells []dto.MatrixCell
	for _, role := range names {
		for res := rbac.ResourceProject; res <= rbac.ResourceRBACConfig; res++ {
			for act := policies.ActionView; act <= policies.ActionDelete; act++ {
				if scope := matrix.ScopeFor(role, res, act); scope != policies.ScopeNone {
					cells = append(cells, dto.MatrixCell{
						Role: role, Resource: policies.ResourceName(res), Action: policies.ActionName(act),
						Scope: policies.ScopeName(scope),
					})
				}
			}
		}
	}
	return cells, nil
}

// Kinds возвращает справочник kind'ов маршрутных проверок.
func (s *Service) Kinds() []policies.KindInfo {
	return policies.Kinds()
}

// Explain отвечает «почему allow/deny» для отладки правил из БД.
func (s *Service) Explain(_ context.Context, in dto.ExplainInput) (dto.ExplainResult, error) {
	res, ok := policies.ParseResource(in.Resource)
	if !ok {
		return dto.ExplainResult{}, errors.BadRequest("неизвестный ресурс " + in.Resource)
	}
	act, ok := policies.ParseAction(in.Action)
	if !ok {
		return dto.ExplainResult{}, errors.BadRequest("неизвестное действие " + in.Action)
	}
	scope := policies.CurrentMatrix().ScopeFor(in.Role, res, act)
	allowed := policies.Authorize(in.Role, res, act, rbac.Owners{
		ProjectOwner: in.ProjectOwner,
		ProcessOwner: in.ProcessOwner,
		Owner:        in.Owner,
	}, in.UserID)
	return dto.ExplainResult{Scope: policies.ScopeName(scope), Allowed: allowed}, nil
}

// apply применяет текущее состояние БД в движок.
func (s *Service) apply(ctx context.Context) error {
	if err := s.store.Reload(ctx); err != nil {
		s.logger.ErrorContext(
			ctx,
			"rbac: применение изменений не удалось (продолжаю на последнем снапшоте)",
			"error",
			err,
		)
	}
	return nil
}

func roleExists(roles []domain.Role, name string) bool {
	for _, r := range roles {
		if r.Name == name {
			return true
		}
	}
	return false
}
