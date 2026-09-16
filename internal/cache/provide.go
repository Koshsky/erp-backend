// Package cache provides the shared Redis connection used by infra
// middlewares (currently the rate limiter, M1).
package cache

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Koshsky/erp-backend/internal/config"
)

// pingTimeout bounds the startup connectivity check.
const pingTimeout = 3 * time.Second

// ProvideRedisClient builds the shared Redis client from the configuration
// (M1). When Redis is disabled, a nil client is returned and the rate limiter
// falls back to the in-memory implementation. A failed ping is logged and
// degraded to nil (fail-open) so the application keeps serving without
// distributed limits rather than refusing to start.
//
//nolint:nilnil // a nil client is the documented "in-memory fallback" value, not an error state
func ProvideRedisClient(cfg config.RedisConfig, logger *slog.Logger) (*redis.Client, error) {
	if !cfg.Enabled {
		logger.Info("redis disabled, rate limiter will use the in-memory fallback")
		return nil, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Address,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  time.Duration(cfg.DialTimeout),
		ReadTimeout:  time.Duration(cfg.ReadTimeout),
		WriteTimeout: time.Duration(cfg.WriteTimeout),
		PoolSize:     cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("redis ping failed, rate limiter will use the in-memory fallback",
			"address", cfg.Address, "error", err)
		_ = client.Close()
		return nil, nil
	}

	logger.Info("redis connected", "address", cfg.Address, "db", cfg.DB)
	return client, nil
}
