package repository

import (
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/timesheet/state/repository/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProvideStateRepository builds the StateRepository repository.
func ProvideStateRepository(logger *slog.Logger, pool *pgxpool.Pool) *StateRepository {
	return &StateRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}
