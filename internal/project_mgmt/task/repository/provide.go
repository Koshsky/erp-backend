package repository

import (
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/repository/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProvideTaskRepository builds the TaskRepository repository.
func ProvideTaskRepository(logger *slog.Logger, pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}
