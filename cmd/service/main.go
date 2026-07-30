package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Koshsky/erp-backend/internal/app"
	"github.com/Koshsky/erp-backend/internal/config"
	"github.com/Koshsky/erp-backend/internal/database"
	appLogger "github.com/Koshsky/erp-backend/internal/logger"

	_ "github.com/Koshsky/erp-backend/docs/swagger"
)

const defaultShutdownTimeout = 5 * time.Second

func main() {
	// load configuration and initialize logger
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := appLogger.New(cfg.LogLevel, cfg.LogFormat)
	if err != nil {
		log.Fatal(err)
	}

	slog.SetDefault(logger)

	// create a context that is canceled on SIGINT or SIGTERM
	runCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// initialize the database pool
	pool, err := database.InitDBPool(cfg.Postgres, logger)
	if err != nil {
		logger.Error("Failed to initialize database pool", "error", err)
		os.Exit(1)
	}

	// initialize the application
	application, err := app.New(cfg, appLogger.WithComponent(logger, "app"), pool)
	if err != nil {
		logger.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}

	logger.Info("service configuration loaded")

	// start HTTP server
	if err = application.Start(); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	// graceful shutdown
	<-runCtx.Done()
	logger.Info("shutdown signal received")

	// stop server with 5 second timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	if err = application.Stop(shutdownCtx); err != nil {
		logger.Error("error during graceful shutdown", "error", err)
	}

	logger.Info("service stopped")
}
