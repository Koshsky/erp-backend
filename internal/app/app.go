package app

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
	"github.com/Koshsky/erp-backend/internal/security/jwt"
)

type App struct {
	cfg        *config.Config
	logger     *slog.Logger
	pool       *pgxpool.Pool
	httpServer *http.Server
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

	if a.cfg.Swagger.Enabled {
		a.runSwaggerServer(router)
	}

	// Register middleware
	router.Use(gin.Recovery())
	a.jwtManager = jwt.NewJWTService(a.cfg.JWT)
	a.authMw = auth.NewMiddleware(a.logger, a.jwtManager)
	router.Use(a.authMw.Middleware())
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
		ReadTimeout:  srvCfg.ReadTimeout,
		WriteTimeout: srvCfg.WriteTimeout,
		IdleTimeout:  srvCfg.IdleTimeout,
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

	if err := a.waitForServer(5 * time.Second); err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}

	return nil
}

func (a *App) waitForServer(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	addr := fmt.Sprintf("localhost:%d", a.cfg.HTTPServer.Port)

	dialer := &net.Dialer{
		Timeout: time.Second,
	}

	for {
		conn, err := dialer.DialContext(ctx, "tcp", addr)
		if err == nil {
			conn.Close()
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for server: %w", ctx.Err())
		default:
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (a *App) Stop(ctx context.Context) error {
	if a.pool != nil {
		a.pool.Close()
	}
	if a.httpServer != nil {
		return a.httpServer.Shutdown(ctx)
	}
	return nil
}
