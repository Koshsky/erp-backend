package service

import (
	"context"
	"log/slog"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/policies"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy/domain"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy/dto"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy/repository"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy/repository/sqlc"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
	userctx "github.com/Koshsky/erp-backend/internal/userctx"
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

// ListPresets returns the preset catalog.
func (s *Service) ListPresets(ctx context.Context) ([]sqlc.ListActivePresetsRow, error) {
	return s.repo.ListActivePresets(ctx)
}

// ListRules returns the active matrix rows.
func (s *Service) ListRules(ctx context.Context) ([]sqlc.ListActivePresetRulesRow, error) {
	return s.repo.ListActiveRules(ctx)
}

// UpsertRule validates and writes a matrix row, then applies the
// changes immediately.
func (s *Service) UpsertRule(ctx context.Context, in dto.PresetRuleInput, updatedBy int64) error {
	if in.Preset == "" {
		return errors.BadRequest("пресет обязателен")
	}
	presets, err := s.repo.ListActivePresets(ctx)
	if err != nil {
		return err
	}
	if !presetExists(presets, in.Preset) {
		return errors.BadRequest("неизвестный пресет " + in.Preset)
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

	if _, upsertErr := s.repo.UpsertRule(
		ctx,
		in.Preset,
		in.Resource,
		in.Action,
		in.Scope,
		&updatedBy,
	); upsertErr != nil {
		return upsertErr
	}
	return s.apply(ctx)
}

// DeleteRule removes a matrix row (archived by the DB trigger).
func (s *Service) DeleteRule(ctx context.Context, id int64) error {
	if err := s.repo.DeleteRule(ctx, id); err != nil {
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

// DeleteRoutePolicy removes a route policy by name (archived).
func (s *Service) DeleteRoutePolicy(ctx context.Context, name string) error {
	if err := s.repo.DeleteRoutePolicy(ctx, name); err != nil {
		return err
	}
	return s.apply(ctx)
}

// Reset restores rules and route policies to the built-in defaults
// (an escape hatch after erroneous edits).
func (s *Service) Reset(ctx context.Context, updatedBy int64) error {
	if err := s.repo.DeleteAllRules(ctx); err != nil {
		return err
	}
	if err := s.repo.DeleteAllRoutePolicies(ctx); err != nil {
		return err
	}
	for _, r := range policies.DefaultMatrixRules() {
		scope := policies.ScopeName(r.Scope)
		if _, err := s.repo.UpsertRule(
			ctx,
			r.Role,
			policies.ResourceName(r.Res),
			policies.ActionName(r.Act),
			scope,
			&updatedBy,
		); err != nil {
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
	_ = ctx // matrix and presets are read from memory/DB without extra queries
	presets, err := s.repo.ListActivePresets(ctx)
	if err != nil {
		return nil, err
	}
	matrix := policies.CurrentMatrix()
	names := make([]string, 0, len(presets)+1)
	names = append(names, userdomain.PresetAdmin)
	for _, p := range presets {
		names = append(names, p.Name)
	}
	var cells []dto.MatrixCell
	for _, preset := range names {
		for res := rbac.ResourceProject; res <= rbac.ResourceAudit; res++ {
			for act := policies.ActionView; act <= policies.ActionDelete; act++ {
				if scope := matrix.ScopeFor(preset, res, act); scope != policies.ScopeNone {
					cells = append(cells, dto.MatrixCell{
						Preset: preset, Resource: policies.ResourceName(res), Action: policies.ActionName(act),
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
	scope := policies.CurrentMatrix().ScopeFor(in.Preset, res, act)
	allowed := policies.Authorize(in.Preset, res, act, rbac.Owners{
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

func presetExists(presets []sqlc.ListActivePresetsRow, name string) bool {
	for _, p := range presets {
		if p.Name == name {
			return true
		}
	}
	return false
}

// MyPermissions returns the caller's principal permissions (all allowed
// actions with an effective scope != none; admin — everything). Used by the
// frontend to display capabilities by permissions rather than by presets.
func (s *Service) MyPermissions(_ context.Context, user userctx.UserContext) []dto.Permission {
	out := []dto.Permission{}
	for res := rbac.ResourceProject; res <= rbac.ResourceAudit; res++ {
		for act := policies.ActionView; act <= policies.ActionDelete; act++ {
			scope := policies.CurrentMatrix().ScopeForUser(user, res, act)
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

// ListUserPermissions returns the editor view of a user's permissions: the
// assigned preset, the current overrides, the preset baseline (without
// overrides) and the resulting effective matrix.
func (s *Service) ListUserPermissions(
	ctx context.Context,
	userID int64,
) (dto.UserPermissionsView, error) {
	presetCode, presetSet, err := s.repo.FindUserPreset(ctx, userID)
	if err != nil {
		return dto.UserPermissionsView{}, err
	}
	rows, err := s.repo.ListUserPermissions(ctx, userID)
	if err != nil {
		return dto.UserPermissionsView{}, err
	}
	var preset *string
	if presetSet {
		preset = &presetCode
	}
	overrides := make([]dto.PermissionOverride, 0, len(rows))
	for _, p := range rows {
		overrides = append(overrides, dto.PermissionOverride{
			Resource: p.Resource,
			Action:   p.Action,
			Scope:    p.Scope,
			Granted:  p.Granted,
		})
	}
	admin := presetSet && presetCode == userdomain.PresetAdmin
	// Preset baseline: the assigned preset without individual overrides.
	base := userctx.UserContext{
		ID:     userID,
		Admin:  admin,
		Preset: presetNameRef(preset),
	}
	// Effective matrix: the caller's principal applied to every resource/action.
	principal := base
	principal.Rules = toRules(overrides)
	return dto.UserPermissionsView{
		UserID:      userID,
		Preset:      preset,
		Admin:       admin,
		Overrides:   overrides,
		PresetScope: s.MyPermissions(ctx, base),
		Effective:   s.MyPermissions(ctx, principal),
	}, nil
}

// ReplaceUserPermissions validates and replaces the user's override set
// (full replacement), then applies the change immediately.
func (s *Service) ReplaceUserPermissions(
	ctx context.Context,
	userID int64,
	in dto.UserPermissionsInput,
	updatedBy int64,
) error {
	rows := make([]sqlc.InsertUserPermissionParams, 0, len(in.Overrides))
	seen := map[string]bool{}
	for _, o := range in.Overrides {
		res, ok := policies.ParseResource(o.Resource)
		if !ok {
			return errors.BadRequest("неизвестный ресурс " + o.Resource)
		}
		if _, okAction := policies.ParseAction(o.Action); !okAction {
			return errors.BadRequest("неизвестное действие " + o.Action)
		}
		key := o.Resource + "/" + o.Action
		if seen[key] {
			return errors.BadRequest("дублируется право " + key)
		}
		seen[key] = true
		scope := "all"
		if o.Granted {
			parsed, okScope := policies.ParseScope(o.Scope)
			if !okScope || parsed == policies.ScopeNone {
				return errors.BadRequest("недопустимая зона " + o.Scope + " (all|own|parent|ancestor)")
			}
			if !policies.ScopeApplicable(res, parsed) {
				return errors.BadRequest("зона " + o.Scope + " неприменима к ресурсу " + o.Resource)
			}
			scope = o.Scope
		}
		rows = append(rows, sqlc.InsertUserPermissionParams{
			UserID:    userID,
			Resource:  o.Resource,
			Action:    o.Action,
			Scope:     scope,
			Granted:   o.Granted,
			UpdatedBy: pgtype.Int8{Int64: updatedBy, Valid: true},
		})
	}
	if err := s.repo.ReplaceUserPermissions(ctx, userID, rows); err != nil {
		return err
	}
	return s.apply(ctx)
}

// toRules converts overrides into the engine's user-rule model.
func toRules(overrides []dto.PermissionOverride) []userctx.PermissionRule {
	out := make([]userctx.PermissionRule, 0, len(overrides))
	for _, o := range overrides {
		out = append(out, userctx.PermissionRule{
			Resource: o.Resource,
			Action:   o.Action,
			Scope:    o.Scope,
			Granted:  o.Granted,
		})
	}
	return out
}

// presetNameRef unwraps a preset pointer ("" — none).
func presetNameRef(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// maxPresetNameLen — maximum preset name length.
const maxPresetNameLen = 32

// CreatePreset creates a preset (or updates an existing one) and applies it.
func (s *Service) CreatePreset(ctx context.Context, in dto.PresetUpsertInput) (sqlc.UpsertPresetRow, error) {
	if err := validatePresetName(in.Name); err != nil {
		return sqlc.UpsertPresetRow{}, err
	}
	preset, err := s.repo.UpsertPreset(ctx, in.Name, in.Description)
	if err != nil {
		return sqlc.UpsertPresetRow{}, err
	}
	_ = s.apply(ctx)
	return preset, nil
}

// UpdatePreset updates the preset description.
func (s *Service) UpdatePreset(
	ctx context.Context,
	name string,
	in dto.PresetUpdateInput,
) (sqlc.UpdatePresetDescriptionRow, error) {
	preset, err := s.repo.UpdatePresetDescription(ctx, name, in.Description)
	if err != nil {
		return sqlc.UpdatePresetDescriptionRow{}, err
	}
	return preset, nil
}

// DeletePreset deletes a preset, its rules and clears the preset ref on
// assigned users (they keep the account but lose the preset's base rights;
// individual overrides survive).
func (s *Service) DeletePreset(ctx context.Context, name string) error {
	if err := s.repo.DeletePreset(ctx, name); err != nil {
		return err
	}
	return s.apply(ctx)
}

// validPresetName — allowed characters of a preset name (system access code):
// lowercase latin letters, digits, "-" and "_".
func validPresetName(name string) bool {
	return regexp.MustCompile(`^[a-z0-9_-]+$`).MatchString(name)
}

// validatePresetName validates a preset name: non-empty, no longer than 32, a code.
func validatePresetName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.BadRequest("имя пресета не может быть пустым")
	}
	if len(name) > maxPresetNameLen {
		return errors.BadRequest("имя пресета не длиннее 32 символов")
	}
	if !validPresetName(name) {
		return errors.BadRequest("имя пресета: только латиница в нижнем регистре, цифры, «-» и «_»")
	}
	return nil
}
