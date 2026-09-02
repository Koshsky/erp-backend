package service

import (
	"context"
	"log/slog"
	"regexp"
	"strings"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/policies"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy/domain"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy/dto"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy/repository"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

// Service — RBAC policy administration (validation + DB writes +
// immediate application through the PolicyStore).
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

// UpsertRule validates and writes a matrix row, then applies the
// changes immediately.
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

// DeleteRule softly deletes a matrix row.
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

// UpsertRoutePolicy validates (kind + parameters against the schema) and writes
// the route policy.
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

// DeleteRoutePolicy softly deletes a route policy by name.
func (s *Service) DeleteRoutePolicy(ctx context.Context, name string) error {
	if err := s.repo.SoftDeleteRoutePolicy(ctx, name); err != nil {
		return err
	}
	return s.apply(ctx)
}

// Reset restores rules and route policies to the built-in defaults
// (an escape hatch after erroneous edits).
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

// EffectiveMatrix returns the effective matrix (with the admin bypass) for the API.
func (s *Service) EffectiveMatrix(ctx context.Context) ([]dto.MatrixCell, error) {
	_ = ctx // matrix and roles are read from memory/DB without extra queries
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
		for res := rbac.ResourceProject; res <= rbac.ResourceOrgStructure; res++ {
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

// Kinds returns the catalog of route policy kinds.
func (s *Service) Kinds() []policies.KindInfo {
	return policies.Kinds()
}

// Explain answers "why allow/deny" for debugging DB-backed rules.
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

// apply applies the current DB state to the engine.
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

// MyPermissions returns the role's principal permissions (all allowed
// actions with ScopeFor != none; admin — everything). Used by the frontend to
// display capabilities by permissions rather than by roles.
func (s *Service) MyPermissions(_ context.Context, role string) []dto.Permission {
	out := []dto.Permission{}
	for res := rbac.ResourceProject; res <= rbac.ResourceAudit; res++ {
		for act := policies.ActionView; act <= policies.ActionDelete; act++ {
			scope := policies.CurrentMatrix().ScopeFor(role, res, act)
			if scope == policies.ScopeNone {
				continue
			}
			out = append(out, dto.Permission{
				Resource: policies.ResourceName(res),
				Action:   policies.ActionName(act),
				Scope:    policies.ScopeName(scope),
			})
		}
	}
	return out
}

// maxRoleNameLen — maximum role name length.
//

const maxRoleNameLen = 32

// CreateRole creates a role (or revives a deleted one) and returns it.
func (s *Service) CreateRole(ctx context.Context, in dto.RoleUpsertInput) (domain.Role, error) {
	if err := validateRoleName(in.Name); err != nil {
		return domain.Role{}, err
	}
	role, err := s.repo.UpsertRole(ctx, in.Name, in.Description)
	if err != nil {
		return domain.Role{}, err
	}
	return role, nil
}

// UpdateRole updates the role description.
func (s *Service) UpdateRole(ctx context.Context, name string, in dto.RoleUpdateInput) (domain.Role, error) {
	role, err := s.repo.UpdateRoleDescription(ctx, name, in.Description)
	if err != nil {
		return domain.Role{}, err
	}
	return role, nil
}

// DeleteRole softly deletes a role and its rules; assigned users
// keep existing but lose permissions (the role disappears from the matrix).
func (s *Service) DeleteRole(ctx context.Context, name string) error {
	if err := s.repo.SoftDeleteRole(ctx, name); err != nil {
		return err
	}
	return s.apply(ctx)
}

// validRoleName — allowed characters of a role name (system access code):
// lowercase latin letters, digits, "-" and "_".
func validRoleName(name string) bool {
	return regexp.MustCompile(`^[a-z0-9_-]+$`).MatchString(name)
}

// validateRoleName validates a role name: non-empty, no longer than 32, a code.
func validateRoleName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.BadRequest("имя роли не может быть пустым")
	}
	if len(name) > maxRoleNameLen {
		return errors.BadRequest("имя роли не длиннее 32 символов")
	}
	if !validRoleName(name) {
		return errors.BadRequest("имя роли: только латиница в нижнем регистре, цифры, «-» и «_»")
	}
	return nil
}
