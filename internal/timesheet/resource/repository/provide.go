package repository

import (
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/timesheet/resource/repository/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProvideResourceRepository builds the ResourceRepository repository.
func ProvideResourceRepository(logger *slog.Logger, pool *pgxpool.Pool) *ResourceRepository {
	return &ResourceRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}
