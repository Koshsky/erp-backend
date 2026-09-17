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
func (r *RuleRepository) ListActivePresets(ctx context.Context) ([]sqlc.ListActivePresetsRow, error) {
	return r.db.ListActivePresets(ctx)
}

// ListActiveRules returns all active matrix rows.
func (r *RuleRepository) ListActiveRules(ctx context.Context) ([]sqlc.ListActivePresetRulesRow, error) {
	return r.db.ListActivePresetRules(ctx)
}

// UpsertRule writes (or updates) a matrix row by the unique key
// (preset, resource, action) and returns the stored row.
func (r *RuleRepository) UpsertRule(
	ctx context.Context,
	preset, resource, action, scope string,
	updatedBy *int64,
) (sqlc.UpsertPresetRuleRow, error) {
	return r.db.UpsertPresetRule(ctx, sqlc.UpsertPresetRuleParams{
		Preset:    preset,
		Resource:  resource,
		Action:    action,
		Scope:     scope,
		UpdatedBy: toInt8(updatedBy),
	})
}

// DeleteRule removes a matrix row (archived by the DB trigger).
func (r *RuleRepository) DeleteRule(ctx context.Context, id int64) error {
	return r.db.DeletePresetRule(ctx, id)
}

// DeleteAllRules removes every matrix row (for reset); each row is archived.
func (r *RuleRepository) DeleteAllRules(ctx context.Context) error {
	return r.db.DeleteAllPresetRules(ctx)
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

// DeleteRoutePolicy removes a route policy by name (archived).
func (r *RuleRepository) DeleteRoutePolicy(ctx context.Context, name string) error {
	return r.db.DeleteRoutePolicy(ctx, name)
}

// DeleteAllRoutePolicies removes every route policy (for reset); each is archived.
func (r *RuleRepository) DeleteAllRoutePolicies(ctx context.Context) error {
	return r.db.DeleteAllRoutePolicies(ctx)
}

// UpsertPreset creates a preset (or updates it by name).
func (r *RuleRepository) UpsertPreset(ctx context.Context, name, description string) (sqlc.UpsertPresetRow, error) {
	return r.db.UpsertPreset(ctx, sqlc.UpsertPresetParams{Name: name, Description: description})
}

// UpdatePresetDescription updates the preset description.
func (r *RuleRepository) UpdatePresetDescription(
	ctx context.Context,
	name, description string,
) (sqlc.UpdatePresetDescriptionRow, error) {
	return r.db.UpdatePresetDescription(ctx, sqlc.UpdatePresetDescriptionParams{
		Name: name, Description: description,
	})
}

// DeletePreset deletes a preset together with its rules and clears the preset
// ref on users (base rights vanish; individual overrides survive). The users
// ref must be cleared before the preset row is removed (FK users.preset →
// rbac_presets(name) is RESTRICT).
func (r *RuleRepository) DeletePreset(ctx context.Context, name string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := sqlc.New(tx)
	if err = q.ClearPresetOnUsers(ctx, name); err != nil {
		return err
	}
	if err = q.DeletePresetRulesByPreset(ctx, name); err != nil {
		return err
	}
	if err = q.DeletePreset(ctx, name); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// ListUserPermissions returns the active per-user overrides of a user.
func (r *RuleRepository) ListUserPermissions(ctx context.Context, userID int64) ([]sqlc.ListUserPermissionsRow, error) {
	return r.db.ListUserPermissions(ctx, userID)
}

// ReplaceUserPermissions atomically replaces the user's override set:
// deletes (archives) the previous rows and inserts the new ones.
func (r *RuleRepository) ReplaceUserPermissions(
	ctx context.Context,
	userID int64,
	perms []sqlc.InsertUserPermissionParams,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := sqlc.New(tx)
	if err = q.DeleteAllUserPermissions(ctx, userID); err != nil {
		return err
	}
	for _, p := range perms {
		if _, err = q.InsertUserPermission(ctx, p); err != nil {
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
