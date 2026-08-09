package repository

import (
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/timesheet/employee/repository/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProvideEmployeeRepository builds the EmployeeRepository repository.
func ProvideEmployeeRepository(logger *slog.Logger, pool *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{
		logger: logger,
		pool:   pool,
		db:     sqlc.New(pool),
	}
}
