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
	"github.com/Koshsky/erp-backend/internal/middleware/auth"
	"github.com/Koshsky/erp-backend/internal/middleware/cors"
	"github.com/Koshsky/erp-backend/internal/middleware/ratelimit"
	"github.com/Koshsky/erp-backend/internal/server/maintenance"
	"github.com/Koshsky/erp-backend/internal/server/profiler"
	"github.com/Koshsky/erp-backend/internal/server/swagger"
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
	modules     []Module
}

// New wires the application with its injected dependencies.
func New(
	cfg *config.Config,
	logger *slog.Logger,
	pool *pgxpool.Pool,
	authMw *auth.Middleware,
	profiler *profiler.Profiler,
	modules []Module,
) (*App, error) {
	return &App{
		cfg:         cfg,
		logger:      logger,
		pool:        pool,
		authMw:      authMw,
		profiler:    profiler,
		maintenance: maintenance.New(cfg.Maintenance, pool, logger),
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

	a.profiler.Start()
	a.maintenance.Start()

	if a.cfg.Swagger.Enabled {
		swagger.Register(router)
	}

	// Register middleware
	router.Use(cors.FromConfig(a.cfg.CORS))
	router.Use(gin.Recovery())
	router.Use(ratelimit.FromConfig(a.cfg.RateLimit, a.logger))
	router.Use(func(c *gin.Context) {
		a.logger.Info("request", "method", c.Request.Method, "path", c.Request.RequestURI)
		c.Next()
	})

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
	if a.pool != nil {
		a.pool.Close()
	}
	if err := a.profiler.Stop(ctx); err != nil {
		return err
	}
	if a.httpServer != nil {
		return a.httpServer.Shutdown(ctx)
	}
	return nil
}
