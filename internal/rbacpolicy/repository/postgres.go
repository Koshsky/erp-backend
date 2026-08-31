//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/rbacpolicy/domain"
	"github.com/Koshsky/erp-backend/internal/rbacpolicy/repository/sqlc"
)

// RuleRepository — доступ к конфигурируемым RBAC-политикам в Postgres.
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

// ListActiveRoles returns the role catalog.
func (r *RuleRepository) ListActiveRoles(ctx context.Context) ([]domain.Role, error) {
	rows, err := r.db.ListActiveRoles(ctx)
	if err != nil {
		return nil, err
	}
	roles := make([]domain.Role, 0, len(rows))
	for _, row := range rows {
		roles = append(roles, domain.Role{ID: row.ID, Name: row.Name, Description: row.Description})
	}
	return roles, nil
}

// ListActiveRules returns all active matrix rows.
func (r *RuleRepository) ListActiveRules(ctx context.Context) ([]domain.Rule, error) {
	rows, err := r.db.ListActiveRules(ctx)
	if err != nil {
		return nil, err
	}
	rules := make([]domain.Rule, 0, len(rows))
	for _, row := range rows {
		rules = append(rules, domain.Rule{
			ID:        row.ID,
			Role:      row.Role,
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
// (role, resource, action) and returns the stored row.
func (r *RuleRepository) UpsertRule(ctx context.Context, rule domain.Rule) (domain.Rule, error) {
	row, err := r.db.UpsertRule(ctx, sqlc.UpsertRuleParams{
		Role:      rule.Role,
		Resource:  rule.Resource,
		Action:    rule.Action,
		Scope:     rule.Scope,
		UpdatedBy: toInt8(rule.UpdatedBy),
	})
	if err != nil {
		return domain.Rule{}, err
	}
	return domain.Rule{
		ID:        row.ID,
		Role:      row.Role,
		Resource:  row.Resource,
		Action:    row.Action,
		Scope:     row.Scope,
		UpdatedBy: fromInt8(row.UpdatedBy),
		UpdatedAt: row.UpdatedAt,
	}, nil
}

// SoftDeleteRule marks a matrix row as deleted (no-op if already deleted).
func (r *RuleRepository) SoftDeleteRule(ctx context.Context, id int64) error {
	return r.db.SoftDeleteRule(ctx, id)
}

// SoftDeleteAllRules marks every matrix row as deleted (для reset).
func (r *RuleRepository) SoftDeleteAllRules(ctx context.Context) error {
	return r.db.SoftDeleteAllRules(ctx)
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

// SoftDeleteAllRoutePolicies marks every route policy as deleted (для reset).
func (r *RuleRepository) SoftDeleteAllRoutePolicies(ctx context.Context) error {
	return r.db.SoftDeleteAllRoutePolicies(ctx)
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

// UpsertRole создаёт роль (или «оживляет» soft-deleted по имени).
func (r *RuleRepository) UpsertRole(ctx context.Context, name, description string) (domain.Role, error) {
	row, err := r.db.UpsertRole(ctx, sqlc.UpsertRoleParams{Name: name, Description: description})
	if err != nil {
		return domain.Role{}, err
	}
	return domain.Role{ID: row.ID, Name: row.Name, Description: row.Description}, nil
}

// UpdateRoleDescription обновляет описание роли.
func (r *RuleRepository) UpdateRoleDescription(ctx context.Context, name, description string) (domain.Role, error) {
	row, err := r.db.UpdateRoleDescription(ctx, sqlc.UpdateRoleDescriptionParams{Name: name, Description: description})
	if err != nil {
		return domain.Role{}, err
	}
	return domain.Role{ID: row.ID, Name: row.Name, Description: row.Description}, nil
}

// SoftDeleteRole мягко удаляет роль вместе с её правилами.
func (r *RuleRepository) SoftDeleteRole(ctx context.Context, name string) error {
	if err := r.db.SoftDeleteRole(ctx, name); err != nil {
		return err
	}
	return r.db.SoftDeleteRulesByRole(ctx, name)
}
