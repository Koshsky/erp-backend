package repository

import (
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/user/repository/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProvideUserRepository builds the UserRepository repository.
func ProvideUserRepository(logger *slog.Logger, pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}
