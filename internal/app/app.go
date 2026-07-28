package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/Koshsky/erp-backend/internal/config"
	"github.com/Koshsky/erp-backend/internal/middleware/auth"
	"github.com/Koshsky/erp-backend/internal/security/jwt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

const APP_PORT = 8080

type App struct {
	cfg        *config.Config
	logger     *slog.Logger
	pool       *pgxpool.Pool
	httpServer *http.Server
	jwtManager *jwt.JWTService
	authMw     *auth.AuthMiddleware
}

func New(ctx context.Context, cfg *config.Config, logger *slog.Logger, pool *pgxpool.Pool) (*App, error) {
	return &App{
		cfg:    cfg,
		logger: logger,
		pool:   pool,
	}, nil
}

func (a *App) Start() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	if a.cfg.SwaggerEnable {
		a.runSwaggerServer(router)
	}

	// Register middleware
	router.Use(gin.Recovery())
	a.jwtManager = jwt.NewJWTService(a.cfg.JWT)
	a.authMw = auth.NewAuthMiddleware(a.logger, a.jwtManager)
	router.Use(a.authMw.Middleware())
	router.Use(func(c *gin.Context) {
		a.logger.Info("request", "method", c.Request.Method, "path", c.Request.RequestURI)
		c.Next()
	})

	// Register routes
	a.registerRoutes(router)

	// Create HTTP server
	a.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", APP_PORT),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		a.logger.Info("starting HTTP server", "addr", a.httpServer.Addr)
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
	deadline := time.Now().Add(timeout)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", APP_PORT), time.Second)
		if err == nil {
			conn.Close()
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for server")
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
