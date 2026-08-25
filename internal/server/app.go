package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/config"
	idempotencypkg "github.com/Koshsky/erp-backend/internal/idempotency"
	"github.com/Koshsky/erp-backend/internal/middleware/auth"
	"github.com/Koshsky/erp-backend/internal/middleware/cors"
	"github.com/Koshsky/erp-backend/internal/middleware/ratelimit"
	rbacpolicysvc "github.com/Koshsky/erp-backend/internal/rbacpolicy/service"
	"github.com/Koshsky/erp-backend/internal/server/maintenance"
	"github.com/Koshsky/erp-backend/internal/server/profiler"
	"github.com/Koshsky/erp-backend/internal/server/swagger"
	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"
)

const (
	defaultServerStartTimeout = 5 * time.Second
	defaultDialTimeout        = 1 * time.Second
	defaultPollInterval       = 100 * time.Millisecond
)

type App struct {
	cfg         *config.Config
	logger      *slog.Logger
	pool        *pgxpool.Pool
	httpServer  *http.Server
	profiler    *profiler.Profiler
	maintenance *maintenance.Normalizer
	authMw      *auth.Middleware
	tracer      *tracingpkg.Tracer
	idemMw      *idempotencypkg.Middleware
	policyStore *rbacpolicysvc.PolicyStore
	modules     []Module
}

// New wires the application with its injected dependencies.
func New(
	cfg *config.Config,
	logger *slog.Logger,
	pool *pgxpool.Pool,
	authMw *auth.Middleware,
	profiler *profiler.Profiler,
	tracer *tracingpkg.Tracer,
	idemMw *idempotencypkg.Middleware,
	policyStore *rbacpolicysvc.PolicyStore,
	modules []Module,
) (*App, error) {
	return &App{
		cfg:         cfg,
		logger:      logger,
		pool:        pool,
		authMw:      authMw,
		profiler:    profiler,
		maintenance: maintenance.New(cfg.Maintenance, pool, logger),
		tracer:      tracer,
		idemMw:      idemMw,
		policyStore: policyStore,
		modules:     modules,
	}, nil
}

// Logger returns the application logger.
func (a *App) Logger() *slog.Logger {
	return a.logger
}

func (a *App) Start() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Only trust the reverse proxy (nginx) as a source of client IP headers.
	// Gin defaults to trusting every peer, which lets remote clients spoof
	// X-Forwarded-For and bypass the per-IP rate limiter. An invalid/empty
	// list falls back to "trust nobody": ClientIP() then returns the proxy's
	// own IP, so all traffic shares one bucket instead of being bypassable.
	if err := router.SetTrustedProxies(a.cfg.HTTPServer.TrustedProxies); err != nil {
		a.logger.Error("invalid trusted_proxies config, rate limit will key on the proxy IP",
			"error", err, "trusted_proxies", a.cfg.HTTPServer.TrustedProxies)
		_ = router.SetTrustedProxies(nil)
	}

	a.profiler.Start()
	a.maintenance.Start()
	if a.policyStore != nil {
		a.policyStore.Start()
	}

	if a.cfg.Swagger.Enabled {
		// Swagger is outside /api/v1; keep a public per-IP wall on it.
		swag := router.Group("/swagger")
		swag.Use(ratelimit.FromConfig(a.cfg.RateLimit, a.logger))
		swagger.Register(swag)
	}

	// Register middleware. The rate limiter is intentionally NOT mounted here
	// globally: public routes are limited per-IP and protected routes per-user
	// (see registerRoutes), so a heavy user behind a shared NAT does not drain
	// a common IP bucket and block their neighbors.
	router.Use(cors.FromConfig(a.cfg.CORS))
	router.Use(gin.Recovery())
	// Корневой span запроса (trace): method/path/status/duration + user_id.
	router.Use(a.tracer.HTTPRootSpan())
	// Резервное текстовое логирование запросов (независимо от трейсинга).
	router.Use(a.requestLog())
	// Register routes
	a.registerRoutes(router)

	// Create HTTP server with configurable settings
	srvCfg := a.cfg.HTTPServer
	a.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", srvCfg.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(srvCfg.ReadTimeout),
		WriteTimeout: time.Duration(srvCfg.WriteTimeout),
		IdleTimeout:  time.Duration(srvCfg.IdleTimeout),
	}

	go func() {
		a.logger.Info("starting HTTP server",
			"addr", a.httpServer.Addr,
			"read_timeout", srvCfg.ReadTimeout,
			"write_timeout", srvCfg.WriteTimeout,
			"idle_timeout", srvCfg.IdleTimeout,
		)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("HTTP server error", "error", err)
		}
	}()

	if err := a.waitForServer(defaultServerStartTimeout); err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}
	return nil
}

func (a *App) waitForServer(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	addr := fmt.Sprintf("localhost:%d", a.cfg.HTTPServer.Port)

	dialer := &net.Dialer{
		Timeout: defaultDialTimeout,
	}

	for {
		conn, err := dialer.DialContext(ctx, "tcp", addr)
		if err == nil {
			_ = conn.Close()
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for server: %w", ctx.Err())
		default:
		}

		time.Sleep(defaultPollInterval)
	}
}

func (a *App) Stop(ctx context.Context) error {
	a.maintenance.Stop(ctx)
	if a.policyStore != nil {
		a.policyStore.Stop()
	}
	if a.pool != nil {
		a.pool.Close()
	}
	if err := a.profiler.Stop(ctx); err != nil {
		return err
	}
	if a.tracer != nil {
		if err := a.tracer.Shutdown(ctx); err != nil {
			a.logger.ErrorContext(ctx, "tracing shutdown failed", "error", err)
		}
	}
	if a.httpServer != nil {
		return a.httpServer.Shutdown(ctx)
	}
	return nil
}

// requestLog пишет в лог сводку по каждому HTTP-запросу (метод, путь, статус,
// длительность) — независимый от трейсинга резерв для текстовых логов.
func (a *App) requestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		a.logger.Info("request",
			"method", c.Request.Method,
			"path", c.Request.RequestURI,
			"status", c.Writer.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
	}
}
