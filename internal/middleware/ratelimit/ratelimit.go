package ratelimit

import (
	"log/slog"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/Koshsky/erp-backend/internal/common/response"
	"github.com/Koshsky/erp-backend/internal/config"
)

// retryAfterSeconds is the Retry-After value sent with a 429 response.
const retryAfterSeconds = "1"

// Config is the rate limiting settings.
type Config struct {
	RequestsPerSecond float64
	Burst             int
	CleanupInterval   time.Duration
	Expiration        time.Duration
}

// FromConfig builds a rate limiting handler from the application configuration.
//
// When rate limiting is disabled, a no-op handler is returned.
func FromConfig(cfg config.RateLimitConfig, logger *slog.Logger) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return New(Config{
		RequestsPerSecond: cfg.RequestsPerSecond,
		Burst:             cfg.Burst,
		CleanupInterval:   time.Duration(cfg.CleanupInterval),
		Expiration:        time.Duration(cfg.Expiration),
	}, logger)
}

// New builds a per-client-IP token bucket middleware.
//
// The returned handler runs a background cleanup goroutine for the lifetime
// of the process; it prunes limiters that were not used for the expiration
// period so the client map does not grow unbounded.
func New(config Config, logger *slog.Logger) gin.HandlerFunc {
	if config.RequestsPerSecond <= 0 {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	if logger == nil {
		logger = slog.Default()
	}

	clientLimit := rate.Limit(config.RequestsPerSecond)

	var mu sync.Mutex
	limiters := make(map[string]*clientLimiter)

	go func() {
		ticker := time.NewTicker(config.CleanupInterval)
		defer ticker.Stop()

		for range ticker.C {
			pruneExpired(limiters, &mu, config.Expiration, time.Now())
		}
	}()

	return func(c *gin.Context) {
		key := c.ClientIP()

		mu.Lock()
		cl, ok := limiters[key]
		if !ok {
			cl = &clientLimiter{
				limiter: rate.NewLimiter(clientLimit, config.Burst),
			}
			limiters[key] = cl
		}
		cl.lastSeen = time.Now()
		mu.Unlock()

		if !cl.limiter.Allow() {
			c.Header("Retry-After", retryAfterSeconds)
			logger.Warn("rate limit exceeded", "client_ip", key)
			response.TooManyRequests(c, "too many requests")
			c.Abort()
			return
		}

		c.Next()
	}
}

// clientLimiter is a per-client token bucket with its last seen timestamp.
type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// pruneExpired removes limiters that were not used for the expiration period.
func pruneExpired(limiters map[string]*clientLimiter, mu *sync.Mutex, expiration time.Duration, now time.Time) {
	mu.Lock()
	defer mu.Unlock()

	for key, cl := range limiters {
		if now.Sub(cl.lastSeen) > expiration {
			delete(limiters, key)
		}
	}
}
