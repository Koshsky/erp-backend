//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/rbacpolicy/domain"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy/repository/sqlc"
)

// RuleRepository — access to configurable RBAC policies in Postgres.
type RuleRepository struct {
	logger *slog.Logger
	pool   *pgxpool.Pool
	db     *sqlc.Queries
}

// NewRuleRepository builds the rule repository.
func NewRuleRepository(logger *slog.Logger, pool *pgxpool.Pool) *RuleRepository {
	return &RuleRepository{
		logger: logger,
		pool:   pool,
		db:     sqlc.New(pool),
	}
}

// ListActivePresets returns the preset catalog.
func (r *RuleRepository) ListActivePresets(ctx context.Context) ([]domain.Preset, error) {
	rows, err := r.db.ListActivePresets(ctx)
	if err != nil {
		return nil, err
	}
	presets := make([]domain.Preset, 0, len(rows))
	for _, row := range rows {
		presets = append(presets, domain.Preset{ID: row.ID, Name: row.Name, Description: row.Description})
	}
	return presets, nil
}

// ListActiveRules returns all active matrix rows.
func (r *RuleRepository) ListActiveRules(ctx context.Context) ([]domain.PresetRule, error) {
	rows, err := r.db.ListActivePresetRules(ctx)
	if err != nil {
		return nil, err
	}
	rules := make([]domain.PresetRule, 0, len(rows))
	for _, row := range rows {
		rules = append(rules, domain.PresetRule{
			ID:        row.ID,
			Preset:    row.Preset,
			Resource:  row.Resource,
			Action:    row.Action,
			Scope:     row.Scope,
			UpdatedBy: fromInt8(row.UpdatedBy),
			UpdatedAt: row.UpdatedAt,
		})
	}
	return rules, nil
}

// UpsertRule writes (or updates) a matrix row by the unique key
// (preset, resource, action) and returns the stored row.
func (r *RuleRepository) UpsertRule(ctx context.Context, rule domain.PresetRule) (domain.PresetRule, error) {
	row, err := r.db.UpsertPresetRule(ctx, sqlc.UpsertPresetRuleParams{
		Preset:    rule.Preset,
		Resource:  rule.Resource,
		Action:    rule.Action,
		Scope:     rule.Scope,
		UpdatedBy: toInt8(rule.UpdatedBy),
	})
	if err != nil {
		return domain.PresetRule{}, err
	}
	return domain.PresetRule{
		ID:        row.ID,
		Preset:    row.Preset,
		Resource:  row.Resource,
		Action:    row.Action,
		Scope:     row.Scope,
		UpdatedBy: fromInt8(row.UpdatedBy),
		UpdatedAt: row.UpdatedAt,
	}, nil
}

// SoftDeleteRule marks a matrix row as deleted (no-op if already deleted).
func (r *RuleRepository) SoftDeleteRule(ctx context.Context, id int64) error {
	return r.db.SoftDeletePresetRule(ctx, id)
}

// SoftDeleteAllRules marks every matrix row as deleted (for reset).
func (r *RuleRepository) SoftDeleteAllRules(ctx context.Context) error {
	return r.db.SoftDeleteAllPresetRules(ctx)
}

// ListActiveRoutePolicies returns all active route policy definitions.
func (r *RuleRepository) ListActiveRoutePolicies(ctx context.Context) ([]domain.RoutePolicy, error) {
	rows, err := r.db.ListActiveRoutePolicies(ctx)
	if err != nil {
		return nil, err
	}
	policies := make([]domain.RoutePolicy, 0, len(rows))
	for _, row := range rows {
		params := map[string]any{}
		if len(row.Params) > 0 {
			if unmarshalErr := json.Unmarshal(row.Params, &params); unmarshalErr != nil {
				return nil, unmarshalErr
			}
		}
		policies = append(policies, domain.RoutePolicy{
			Name:      row.Name,
			Kind:      row.Kind,
			Params:    params,
			Active:    row.Active,
			UpdatedBy: fromInt8(row.UpdatedBy),
			UpdatedAt: row.UpdatedAt,
		})
	}
	return policies, nil
}

// UpsertRoutePolicy writes (or updates) a route policy by name.
func (r *RuleRepository) UpsertRoutePolicy(ctx context.Context, p domain.RoutePolicy) (domain.RoutePolicy, error) {
	raw, err := json.Marshal(p.Params)
	if err != nil {
		return domain.RoutePolicy{}, err
	}
	row, err := r.db.UpsertRoutePolicy(ctx, sqlc.UpsertRoutePolicyParams{
		Name:      p.Name,
		Kind:      p.Kind,
		Params:    raw,
		Active:    p.Active,
		UpdatedBy: toInt8(p.UpdatedBy),
	})
	if err != nil {
		return domain.RoutePolicy{}, err
	}
	params := map[string]any{}
	if len(row.Params) > 0 {
		if unmarshalErr := json.Unmarshal(row.Params, &params); unmarshalErr != nil {
			return domain.RoutePolicy{}, unmarshalErr
		}
	}
	return domain.RoutePolicy{
		Name:      row.Name,
		Kind:      row.Kind,
		Params:    params,
		Active:    row.Active,
		UpdatedBy: fromInt8(row.UpdatedBy),
		UpdatedAt: row.UpdatedAt,
	}, nil
}

// SoftDeleteRoutePolicy marks a route policy as deleted by name.
func (r *RuleRepository) SoftDeleteRoutePolicy(ctx context.Context, name string) error {
	return r.db.SoftDeleteRoutePolicy(ctx, name)
}

// SoftDeleteAllRoutePolicies marks every route policy as deleted (for reset).
func (r *RuleRepository) SoftDeleteAllRoutePolicies(ctx context.Context) error {
	return r.db.SoftDeleteAllRoutePolicies(ctx)
}

// UpsertPreset creates a preset (or revives a soft-deleted one by name).
func (r *RuleRepository) UpsertPreset(ctx context.Context, name, description string) (domain.Preset, error) {
	row, err := r.db.UpsertPreset(ctx, sqlc.UpsertPresetParams{Name: name, Description: description})
	if err != nil {
		return domain.Preset{}, err
	}
	return domain.Preset{ID: row.ID, Name: row.Name, Description: row.Description}, nil
}

// UpdatePresetDescription updates the preset description.
func (r *RuleRepository) UpdatePresetDescription(
	ctx context.Context,
	name, description string,
) (domain.Preset, error) {
	row, err := r.db.UpdatePresetDescription(ctx, sqlc.UpdatePresetDescriptionParams{
		Name: name, Description: description,
	})
	if err != nil {
		return domain.Preset{}, err
	}
	return domain.Preset{ID: row.ID, Name: row.Name, Description: row.Description}, nil
}

// SoftDeletePreset softly deletes a preset together with its rules and clears
// the preset ref on users (base rights vanish; individual overrides survive).
func (r *RuleRepository) SoftDeletePreset(ctx context.Context, name string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := sqlc.New(tx)
	if err = q.SoftDeletePreset(ctx, name); err != nil {
		return err
	}
	if err = q.SoftDeletePresetRulesByPreset(ctx, name); err != nil {
		return err
	}
	if err = q.ClearPresetOnUsers(ctx, name); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// ListUserPermissions returns the active per-user overrides of a user.
func (r *RuleRepository) ListUserPermissions(ctx context.Context, userID int64) ([]domain.UserPermission, error) {
	rows, err := r.db.ListUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.UserPermission, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.UserPermission{
			ID:        row.ID,
			UserID:    row.UserID,
			Resource:  row.Resource,
			Action:    row.Action,
			Scope:     row.Scope,
			Granted:   row.Granted,
			UpdatedBy: fromInt8(row.UpdatedBy),
			UpdatedAt: row.UpdatedAt,
		})
	}
	return out, nil
}

// ReplaceUserPermissions atomically replaces the user's override set:
// soft-deletes the previous rows and inserts the new ones.
func (r *RuleRepository) ReplaceUserPermissions(
	ctx context.Context,
	userID int64,
	rows []domain.UserPermission,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := sqlc.New(tx)
	if err = q.SoftDeleteAllUserPermissions(ctx, userID); err != nil {
		return err
	}
	for _, p := range rows {
		if _, err = q.InsertUserPermission(ctx, sqlc.InsertUserPermissionParams{
			UserID:    p.UserID,
			Resource:  p.Resource,
			Action:    p.Action,
			Scope:     p.Scope,
			Granted:   p.Granted,
			UpdatedBy: toInt8(p.UpdatedBy),
		}); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// FindUserPreset returns the user's assigned preset and whether it is set
// (exists=false — the user has no preset).
func (r *RuleRepository) FindUserPreset(ctx context.Context, userID int64) (string, bool, error) {
	preset, err := r.db.FindUserPreset(ctx, userID)
	if err != nil {
		return "", false, err
	}
	if !preset.Valid {
		return "", false, nil
	}
	return preset.String, true, nil
}

// ListUserPrincipals returns the preset of every active user (user_id, preset)
// — the base of the in-memory principal snapshot.
func (r *RuleRepository) ListUserPrincipals(ctx context.Context) ([]UserPreset, error) {
	rows, err := r.db.ListUserPrincipals(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]UserPreset, 0, len(rows))
	for _, row := range rows {
		out = append(out, UserPreset{
			UserID: row.UserID,
			Preset: row.Preset,
		})
	}
	return out, nil
}

// ListAllUserPermissions returns every active per-user override (for the
// in-memory principal snapshot; grouped by the caller).
func (r *RuleRepository) ListAllUserPermissions(ctx context.Context) ([]UserPermissionRef, error) {
	rows, err := r.db.ListAllUserPermissions(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]UserPermissionRef, 0, len(rows))
	for _, row := range rows {
		out = append(out, UserPermissionRef{
			UserID:   row.UserID,
			Resource: row.Resource,
			Action:   row.Action,
			Scope:    row.Scope,
			Granted:  row.Granted,
		})
	}
	return out, nil
}

// UserPreset — user_id → preset (from ListUserPrincipals).
type UserPreset struct {
	UserID int64
	Preset sql.NullString
}

// UserPermissionRef — a per-user override reference (from ListAllUserPermissions).
type UserPermissionRef struct {
	UserID   int64
	Resource string
	Action   string
	Scope    string
	Granted  bool
}

func toInt8(v *int64) pgtype.Int8 {
	if v == nil {
		return pgtype.Int8{}
	}
	return pgtype.Int8{Int64: *v, Valid: true}
}

func fromInt8(v pgtype.Int8) *int64 {
	if !v.Valid {
		return nil
	}
	out := v.Int64
	return &out
}
