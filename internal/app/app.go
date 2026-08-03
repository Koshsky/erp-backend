package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/config"
	"github.com/Koshsky/erp-backend/internal/middleware/auth"
	"github.com/Koshsky/erp-backend/internal/middleware/cors"
	"github.com/Koshsky/erp-backend/internal/security/jwt"
)

const (
	defaultServerStartTimeout = 5 * time.Second
	defaultDialTimeout        = 1 * time.Second
	defaultPollInterval       = 100 * time.Millisecond
	defaultHeaderReadTimeout  = 5 * time.Second
)

type App struct {
	cfg        *config.Config
	logger     *slog.Logger
	pool       *pgxpool.Pool
	httpServer *http.Server
	profiler   *http.Server
	jwtManager *jwt.Service
	authMw     *auth.Middleware
}

func New(cfg *config.Config, logger *slog.Logger, pool *pgxpool.Pool) (*App, error) {
	return &App{
		cfg:    cfg,
		logger: logger,
		pool:   pool,
	}, nil
}

func (a *App) Start() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	if a.cfg.Profiling.Enabled {
		a.startProfiler()
	}

	if a.cfg.Swagger.Enabled {
		a.runSwaggerServer(router)
	}

	// Register middleware
	router.Use(cors.Development())
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		a.logger.Info("request", "method", c.Request.Method, "path", c.Request.RequestURI)
		c.Next()
	})

	a.jwtManager = jwt.NewJWTService(a.cfg.JWT)
	a.authMw = auth.NewMiddleware(a.logger, a.jwtManager)

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

// startProfiler runs a separate HTTP server with pprof profiling.
func (a *App) startProfiler() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	a.profiler = &http.Server{
		Addr:              a.cfg.Profiling.Address,
		Handler:           mux,
		ReadHeaderTimeout: defaultHeaderReadTimeout,
	}

	go func() {
		a.logger.Info("starting profiler", "addr", a.profiler.Addr)
		if err := a.profiler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("profiler server error", "error", err)
		}
	}()
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
	if a.pool != nil {
		a.pool.Close()
	}
	if a.profiler != nil {
		if err := a.profiler.Shutdown(ctx); err != nil {
			return err
		}
	}
	if a.httpServer != nil {
		return a.httpServer.Shutdown(ctx)
	}
	return nil
}
