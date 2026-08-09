package repository

import (
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/repository/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProvideProjectRepository builds the ProjectRepository repository.
func ProvideProjectRepository(logger *slog.Logger, pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}
