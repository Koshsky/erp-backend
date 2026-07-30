package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/config"
)

func InitDBPool(pgCfg config.PostgresConfig, logger *slog.Logger) (*pgxpool.Pool, error) {
	const op = "initDBPool"

	cfg, err := pgxpool.ParseConfig(pgCfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("%s: parse config: %w", op, err)
	}

	cfg.MaxConns = pgCfg.MaxConns
	cfg.MinConns = pgCfg.MinConns
	cfg.MaxConnLifetime = pgCfg.MaxConnLifetime
	cfg.MaxConnIdleTime = pgCfg.MaxConnIdleTime
	cfg.HealthCheckPeriod = pgCfg.HealthCheckPeriod
	cfg.ConnConfig.ConnectTimeout = pgCfg.ConnectTimeout

	ctx, cancel := context.WithTimeout(context.Background(), pgCfg.ConnectTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: create pool: %w", op, err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: ping: %w", op, err)
	}

	stats := pool.Stat()
	logger.Info("database connected",
		"max_conns", cfg.MaxConns,
		"min_conns", cfg.MinConns,
		"acquired_conns", stats.AcquiredConns(),
		"total_conns", stats.TotalConns(),
	)

	return pool, nil
}
