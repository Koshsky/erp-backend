package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	assignmentDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/delivery"
	assignmentDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/domain"
	assignmentRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/repository"
	milestoneDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/delivery"
	milestoneDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/domain"
	milestoneRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/repository"
	processDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/process/delivery"
	processDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/process/domain"
	processRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/process/repository"
	projectDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/project/delivery"
	projectDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/project/domain"
	projectRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/project/repository"
	resourceDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/delivery"
	resourceDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/domain"
	resourceRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/repository"
	taskDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/task/delivery"
	taskDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/task/domain"
	taskRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/task/repository"
	schedulingDelivery "github.com/Koshsky/erp-backend/internal/scheduling/delivery"
	schedulingDomain "github.com/Koshsky/erp-backend/internal/scheduling/domain"
	schedulingRepo "github.com/Koshsky/erp-backend/internal/scheduling/repository"
	"github.com/Koshsky/erp-backend/internal/security/auth"
	"github.com/Koshsky/erp-backend/internal/security/password"
	userDelivery "github.com/Koshsky/erp-backend/internal/user/delivery"
	userDomain "github.com/Koshsky/erp-backend/internal/user/domain"
	userRepo "github.com/Koshsky/erp-backend/internal/user/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

const APP_PORT = 8080

type App struct {
	enableSwagger bool
	logger        *slog.Logger
	pool          *pgxpool.Pool
	httpServer    *http.Server
}

func New(ctx context.Context, enableSwagger bool, logger *slog.Logger, pool *pgxpool.Pool) (*App, error) {
	return &App{
		enableSwagger: enableSwagger,
		logger:        logger,
		pool:          pool,
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
	authMiddleware := auth.NewAuthMiddleware()
	router.Use(authMiddleware.Middleware())
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

func (a *App) registerRoutes(router *gin.Engine) {
	// --- Scheduling ---
	schedulingQueries := schedulingRepo.NewSchedulingRepository(a.logger, a.pool)
	schedulingSvc := schedulingDomain.NewSchedulingService(a.logger, schedulingQueries)
	schedulingHandler := schedulingDelivery.NewSchedulingHandler(a.logger, schedulingSvc)

	// --- User ---
	userQueries := userRepo.NewUserRepository(a.logger, a.pool)
	userHasher := password.NewBcryptHasher()
	userSvc := userDomain.NewUserService(a.logger, userQueries, userHasher)
	userHandler := userDelivery.NewUserHandler(a.logger, userSvc)

	// --- Task ---
	taskQueries := taskRepo.NewTaskRepository(a.logger, a.pool)
	taskSvc := taskDomain.NewTaskService(a.logger, taskQueries)
	taskHandler := taskDelivery.NewTaskHandler(a.logger, taskSvc)

	// --- Resource ---
	resourceQueries := resourceRepo.NewResourceRepository(a.logger, a.pool)
	resourceSvc := resourceDomain.NewResourceService(a.logger, resourceQueries)
	resourceHandler := resourceDelivery.NewResourceHandler(a.logger, resourceSvc)

	// --- Project ---
	projectQueries := projectRepo.NewProjectRepository(a.logger, a.pool)
	projectSvc := projectDomain.NewProjectService(a.logger, projectQueries)
	projectHandler := projectDelivery.NewProjectHandler(a.logger, projectSvc)

	// --- Process ---
	processQueries := processRepo.NewProcessRepository(a.logger, a.pool)
	processSvc := processDomain.NewProcessService(a.logger, processQueries)
	processHandler := processDelivery.NewProcessHandler(a.logger, processSvc)

	// --- Milestone ---
	milestoneQueries := milestoneRepo.NewMilestoneRepository(a.logger, a.pool)
	milestoneSvc := milestoneDomain.NewMilestoneService(a.logger, milestoneQueries)
	milestoneHandler := milestoneDelivery.NewMilestoneHandler(a.logger, milestoneSvc)

	// --- Assignment ---
	assignmentQueries := assignmentRepo.NewAssignmentRepository(a.logger, a.pool)
	assignmentSvc := assignmentDomain.NewAssignmentService(a.logger, assignmentQueries)
	assignmentHandler := assignmentDelivery.NewAssignmentHandler(a.logger, assignmentSvc)

	// Register routes
	api := router.Group("/api")
	schedulingHandler.RegisterRoutes(api)
	userHandler.RegisterRoutes(api)
	taskHandler.RegisterRoutes(api)
	resourceHandler.RegisterRoutes(api)
	projectHandler.RegisterRoutes(api)
	processHandler.RegisterRoutes(api)
	milestoneHandler.RegisterRoutes(api)
	assignmentHandler.RegisterRoutes(api)
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
