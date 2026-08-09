package database

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/config"
)

// ProvidePostgresDB initializes the database pool.
func ProvidePostgresDB(pgCfg config.PostgresConfig, logger *slog.Logger) (*pgxpool.Pool, error) {
	return InitDBPool(pgCfg, logger)
}
