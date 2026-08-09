package repository

import (
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/planning/repository/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProvidePlanningRepository builds the PlanningRepository repository.
func ProvidePlanningRepository(logger *slog.Logger, pool *pgxpool.Pool) *PlanningRepository {
	return &PlanningRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}
