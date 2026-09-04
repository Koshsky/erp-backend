//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/timesheet/state/domain"
	"github.com/Koshsky/erp-backend/internal/timesheet/state/repository/sqlc"
	errapi "github.com/Koshsky/erp-backend/pkg/errors"
)

type StateRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

// NewStateRepository builds the StateRepository repository.
func NewStateRepository(logger *slog.Logger, pool *pgxpool.Pool) *StateRepository {
	return &StateRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *StateRepository) CreateState(ctx context.Context, state domain.State) (*domain.State, error) {
	row, err := r.db.CreateState(ctx, sqlc.CreateStateParams{
		Code:        state.Code,
		Name:        state.Name,
		IsAvailable: state.IsAvailable,
	})
	if err != nil {
		// Idempotent create: the code already exists (ON CONFLICT
		// inserted nothing) — this is a business-key conflict, not an internal error.
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errapi.Conflict("state with this code already exists")
		}
		return nil, err
	}
	mapped := mapState(row)
	return &mapped, nil
}

func (r *StateRepository) FindState(ctx context.Context, id int64) (*domain.State, error) {
	row, err := r.db.FindState(ctx, id)
	if err != nil {
		return nil, err
	}
	mapped := mapState(row)
	return &mapped, nil
}

func (r *StateRepository) UpdateState(ctx context.Context, state domain.State) (*domain.State, error) {
	row, err := r.db.UpdateState(ctx, sqlc.UpdateStateParams{
		StateID:     state.ID,
		Code:        state.Code,
		Name:        state.Name,
		IsAvailable: state.IsAvailable,
	})
	if err != nil {
		return nil, err
	}
	mapped := mapState(row)
	return &mapped, nil
}

func (r *StateRepository) DeleteState(ctx context.Context, id int64) error {
	return r.db.DeleteState(ctx, id)
}

func (r *StateRepository) ListStates(ctx context.Context) ([]domain.State, error) {
	rows, err := r.db.ListStates(ctx)
	if err != nil {
		return nil, err
	}
	states := make([]domain.State, 0, len(rows))
	for _, row := range rows {
		states = append(states, mapState(row))
	}
	return states, nil
}

func mapState(row sqlc.State) domain.State {
	return domain.State{
		ID:          row.ID,
		Code:        row.Code,
		Name:        row.Name,
		IsAvailable: row.IsAvailable,
	}
}
