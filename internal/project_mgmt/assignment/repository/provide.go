package repository

import (
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/repository/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProvideAssignmentRepository builds the AssignmentRepository repository.
func ProvideAssignmentRepository(logger *slog.Logger, pool *pgxpool.Pool) *AssignmentRepository {
	return &AssignmentRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}
