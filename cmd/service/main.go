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
)

// @title           Enterprise Resource Planning
// @version         1.0
// @description     For managing the enterprise's universal resources
// @termsOfService  http://swagger.io/terms/

// @contact.name   Shmonov Matvey
// @contact.url    https://t.me/Koshsky
// @contact.email  shmonov.mv@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description                 "Введите JWT токен в формате: Bearer {token}"

// @externalDocs.description  Документация ERP (заглушка)
// @externalDocs.url          https://swagger.io/resources/open-api/
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
	runCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// initialize the database pool
	pool, err := database.InitDBPool(cfg.DatabaseURL, logger)
	if err != nil {
		log.Fatal(err)
	}

	// initialize the application
	application, err := app.New(runCtx, cfg, appLogger.WithComponent(logger, "app"), pool)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("service configuration loaded")

	// start HTTP server
	if err := application.Start(); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	// graceful shutdown
	<-runCtx.Done()
	logger.Info("shutdown signal received")

	// stop server with 5 second timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := application.Stop(shutdownCtx); err != nil {
		logger.Error("error during graceful shutdown", "error", err)
	}

	logger.Info("service stopped")
}
