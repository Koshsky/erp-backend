package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Projects routes
			projects := v1.Group("/projects")
			{
				projects.GET("", h.ListProjects)
				projects.GET("/:id", h.GetProject)
				projects.POST("", h.CreateProject)
				projects.DELETE("/:id", h.DeleteProject)
				projects.PATCH("/:id", h.UpdateProject)

				// Processes under projects
				projects.GET("/:id/processes", h.ListProcessesByProject)
			}

			// Processes routes
			processes := v1.Group("/processes")
			{
				processes.GET("/:id", h.GetProcess)
				processes.POST("", h.CreateProcess)
				processes.DELETE("/:id", h.DeleteProcess)
				processes.PATCH("/:id", h.UpdateProcess)

				// Milestones under processes
				processes.GET("/:id/milestones", h.ListMilestonesByProcess)
				// Tasks under processes
				processes.GET("/:id/tasks", h.ListTasksByProcess)
			}

			// Timeline routes
			timeline := v1.Group("/timeline")
			{
				timeline.GET("/processes/:id", h.GetProcessTimeline)
			}

			// Milestones routes
			milestones := v1.Group("/milestones")
			{
				milestones.GET("/:id", h.GetMilestone)
				milestones.POST("", h.CreateMilestone)
				milestones.DELETE("/:id", h.DeleteMilestone)
				milestones.PATCH("/:id", h.UpdateMilestone)
			}

			// Tasks routes
			tasks := v1.Group("/tasks")
			{
				tasks.GET("/:id", h.GetTask)
				tasks.POST("", h.CreateTask)
				tasks.DELETE("/:id", h.DeleteTask)
				tasks.PATCH("/:id", h.UpdateTask)

				// Assignments under tasks
				tasks.GET("/:id/assignments", h.ListAssignmentsByTask)
			}

			// Resources routes
			resources := v1.Group("/resources")
			{
				resources.GET("", h.ListResources)
				resources.GET("/:id", h.GetResource)
				resources.POST("", h.CreateResource)
				resources.DELETE("/:id", h.DeleteResource)
				resources.PATCH("/:id", h.UpdateResource)
				resources.GET("/usage", h.GetResourceUsage)
			}

			// Users routes
			users := v1.Group("/users")
			{
				users.GET("", h.ListUsers)
				users.GET("/:id", h.GetUser)
				users.POST("", h.CreateUser)
				users.DELETE("/:id", h.DeleteUser)
				// users.PATCH("/:id/role", h.UpdateUserRole)
				// users.PATCH("/:id/password", h.UpdateUserPassword)
			}

			// Assignments routes
			assignments := v1.Group("/assignments")
			{
				assignments.GET("/:id", h.GetAssignment)
				assignments.POST("", h.CreateAssignment)
				assignments.DELETE("/:id", h.DeleteAssignment)
				assignments.PATCH("/:id", h.UpdateAssignment)
			}
		}
	}
}
