package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/Koshsky/erp/api/internal/handler"
	"github.com/Koshsky/erp/api/internal/service"
	"github.com/gin-gonic/gin"
)

const APP_PORT = 8080

type App struct {
	enableSwagger bool
	logger        *slog.Logger
	repository    service.Repository
	service       *service.Service
	handler       *handler.Handler
	httpServer    *http.Server
}

func New(ctx context.Context, enableSwagger bool, logger *slog.Logger, repository service.Repository) (*App, error) {
	serviceLayer := service.New(logger, repository)
	handlerLayer := handler.New(logger, serviceLayer)
	return &App{
		enableSwagger: enableSwagger,
		logger:        logger,
		repository:    repository,
		service:       serviceLayer,
		handler:       handlerLayer,
	}, nil
}

func (a *App) Start() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	if a.enableSwagger {
		a.runSwaggerServer(router)
	}

	// Register middleware
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		a.logger.Info("request", "method", c.Request.Method, "path", c.Request.RequestURI)
		c.Next()
	})

	// Register routes
	a.handler.RegisterRoutes(router)

	// Create HTTP server
	a.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", APP_PORT),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		a.logger.Info("starting HTTP server", "addr", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("HTTP server error", "error", err)
		}
	}()

	// Wait for server to be ready
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
	if a != nil && a.repository != nil {
		a.repository.Close()
	}
	if a.httpServer != nil {
		return a.httpServer.Shutdown(ctx)
	}

	return nil
}
