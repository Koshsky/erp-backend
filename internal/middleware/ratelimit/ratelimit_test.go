package ratelimit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/config"
	"github.com/Koshsky/erp-backend/internal/middleware/ratelimit"
)

func doRequest(handler gin.HandlerFunc, clientIP string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(handler)
	router.Any("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = clientIP
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestAllowsRequestsWithinLimit(t *testing.T) {
	t.Parallel()
	cfg := config.RateLimitConfig{
		Enabled:           true,
		RequestsPerSecond: 10,
		Burst:             5,
		CleanupInterval:   config.Duration(time.Minute),
		Expiration:        config.Duration(time.Minute),
	}

	handler := ratelimit.FromConfig(cfg, nil)

	for i := range 5 {
		rec := doRequest(handler, "192.168.0.1:1234")
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d status = %d, want %d", i+1, rec.Code, http.StatusOK)
		}
	}
}

func TestRejectsRequestsOverBurst(t *testing.T) {
	t.Parallel()
	cfg := config.RateLimitConfig{
		Enabled:           true,
		RequestsPerSecond: 10,
		Burst:             2,
		CleanupInterval:   config.Duration(time.Minute),
		Expiration:        config.Duration(time.Minute),
	}

	handler := ratelimit.FromConfig(cfg, nil)

	for i := range 3 {
		rec := doRequest(handler, "192.168.0.2:1234")
		if i < 2 && rec.Code != http.StatusOK {
			t.Fatalf("request %d status = %d, want %d", i+1, rec.Code, http.StatusOK)
		}
		if i == 2 && rec.Code != http.StatusTooManyRequests {
			t.Fatalf("request %d status = %d, want %d", i+1, rec.Code, http.StatusTooManyRequests)
		}
		if i == 2 && rec.Header().Get("Retry-After") != "1" {
			t.Errorf("Retry-After = %q, want \"1\"", rec.Header().Get("Retry-After"))
		}
	}
}

func TestLimitIsPerClient(t *testing.T) {
	t.Parallel()
	cfg := config.RateLimitConfig{
		Enabled:           true,
		RequestsPerSecond: 10,
		Burst:             1,
		CleanupInterval:   config.Duration(time.Minute),
		Expiration:        config.Duration(time.Minute),
	}

	handler := ratelimit.FromConfig(cfg, nil)

	if rec := doRequest(handler, "192.168.0.3:1234"); rec.Code != http.StatusOK {
		t.Fatalf("first client request status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec := doRequest(handler, "192.168.0.3:1234"); rec.Code != http.StatusTooManyRequests {
		t.Fatalf("first client repeat status = %d, want %d", rec.Code, http.StatusTooManyRequests)
	}
	if rec := doRequest(handler, "192.168.0.4:1234"); rec.Code != http.StatusOK {
		t.Fatalf("second client request status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestDisabledConfigPassesEverything(t *testing.T) {
	t.Parallel()
	cfg := config.RateLimitConfig{Enabled: false}

	handler := ratelimit.FromConfig(cfg, nil)

	for i := range 10 {
		rec := doRequest(handler, "192.168.0.5:1234")
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d status = %d, want %d", i+1, rec.Code, http.StatusOK)
		}
	}
}

func TestZeroRateIsNoOp(t *testing.T) {
	t.Parallel()
	cfg := config.RateLimitConfig{
		Enabled:           true,
		RequestsPerSecond: 0,
	}

	handler := ratelimit.FromConfig(cfg, nil)

	for i := range 3 {
		rec := doRequest(handler, "192.168.0.6:1234")
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d status = %d, want %d", i+1, rec.Code, http.StatusOK)
		}
	}
}
