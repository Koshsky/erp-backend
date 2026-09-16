package ratelimit_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/Koshsky/erp-backend/internal/config"
	"github.com/Koshsky/erp-backend/internal/middleware/ratelimit"
)

// newRedisProvider builds a Provider over miniredis with a discard logger.
func newRedisProvider(t *testing.T) (*ratelimit.Provider, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	return ratelimit.ProvideProvider(client, config.RedisConfig{
		Enabled:   true,
		KeyPrefix: "test:rl",
	}, slog.New(slog.DiscardHandler)), mr
}

// runProvider executes the provider-built middleware once and returns the
// status and Retry-After header.
func runProvider(handler gin.HandlerFunc) (int, string) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(handler)
	r.GET("/probe", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/probe", nil))
	return w.Code, w.Header().Get("Retry-After")
}

// TestRedisLimiterBurstableAndBlocked checks the token bucket semantics: a
// burst of requests passes, the next one is throttled with Retry-After.
func TestRedisLimiterBurstableAndBlocked(t *testing.T) {
	t.Parallel()
	provider, mr := newRedisProvider(t)
	handler := provider.New(ratelimit.Config{
		RequestsPerSecond: 1,
		Burst:             2,
		Expiration:        5 * time.Minute,
	})

	for i := range 2 {
		status, _ := runProvider(handler)
		if status != http.StatusOK {
			t.Fatalf("request %d: status = %d, want 200", i+1, status)
		}
	}
	status, retryAfter := runProvider(handler)
	if status != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want 429", status)
	}
	if retryAfter == "" {
		t.Fatal("Retry-After header is missing")
	}

	keys := mr.Keys()
	if len(keys) != 1 {
		t.Fatalf("redis keys = %v, want exactly one bucket key", keys)
	}
}

// TestRedisLimiterRefills checks that time passing refills the bucket. The
// script reads Redis TIME, so the refill is driven by real elapsed time
// (miniredis FastForward does not advance the TIME command).
func TestRedisLimiterRefills(t *testing.T) {
	t.Parallel()
	provider, _ := newRedisProvider(t)
	handler := provider.New(ratelimit.Config{
		RequestsPerSecond: 100,
		Burst:             1,
		Expiration:        5 * time.Minute,
	})

	if status, _ := runProvider(handler); status != http.StatusOK {
		t.Fatalf("first request: status = %d, want 200", status)
	}
	if status, _ := runProvider(handler); status != http.StatusTooManyRequests {
		t.Fatalf("second request: status = %d, want 429", status)
	}

	time.Sleep(30 * time.Millisecond) // refills 3 tokens at 100/s
	if status, _ := runProvider(handler); status != http.StatusOK {
		t.Fatalf("request after refill: status = %d, want 200", status)
	}
}

// TestRedisLimiterKeyedByClientIP checks keys are scoped per IP (default key).
func TestRedisLimiterKeyedByClientIP(t *testing.T) {
	t.Parallel()
	provider, mr := newRedisProvider(t)
	handler := provider.New(ratelimit.Config{
		RequestsPerSecond: 1,
		Burst:             1,
		Expiration:        5 * time.Minute,
	})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(handler)
	r.GET("/probe", func(c *gin.Context) { c.Status(http.StatusOK) })

	newReq := func(ip string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/probe", nil)
		req.RemoteAddr = ip + ":1234"
		r.ServeHTTP(w, req)
		return w
	}

	// Same IP: first passes, second is throttled.
	if newReq("10.0.0.1").Code != http.StatusOK {
		t.Fatal("first request from 10.0.0.1: want 200")
	}
	if newReq("10.0.0.1").Code != http.StatusTooManyRequests {
		t.Fatal("second request from 10.0.0.1: want 429")
	}
	// Different IP has its own bucket.
	if newReq("10.0.0.2").Code != http.StatusOK {
		t.Fatal("request from 10.0.0.2: want 200")
	}

	if len(mr.Keys()) != 2 {
		t.Fatalf("redis keys = %v, want two per-IP buckets", mr.Keys())
	}
}

// TestRedisLimiterDisabledNoop checks disabled limits short-circuit.
func TestRedisLimiterDisabledNoop(t *testing.T) {
	t.Parallel()
	provider, _ := newRedisProvider(t)
	handler := provider.FromConfig(config.RateLimitConfig{Enabled: false})
	if status, _ := runProvider(handler); status != http.StatusOK {
		t.Fatalf("status = %d, want 200", status)
	}
}

// TestProviderMemoryFallback checks the in-memory path when Redis is nil.
func TestProviderMemoryFallback(t *testing.T) {
	t.Parallel()
	provider := ratelimit.ProvideProvider(nil, config.RedisConfig{KeyPrefix: "x"}, slog.New(slog.DiscardHandler))
	handler := provider.New(ratelimit.Config{
		RequestsPerSecond: 1,
		Burst:             1,
		Expiration:        time.Minute,
		CleanupInterval:   time.Minute,
	})
	if status, _ := runProvider(handler); status != http.StatusOK {
		t.Fatalf("first request: status = %d, want 200", status)
	}
	if status, _ := runProvider(handler); status != http.StatusTooManyRequests {
		t.Fatalf("second request: status = %d, want 429 (in-memory bucket)", status)
	}
}
