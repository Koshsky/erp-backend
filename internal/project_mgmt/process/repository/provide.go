package repository

import (
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/repository/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProvideProcessRepository builds the ProcessRepository repository.
func ProvideProcessRepository(logger *slog.Logger, pool *pgxpool.Pool) *ProcessRepository {
	return &ProcessRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}
