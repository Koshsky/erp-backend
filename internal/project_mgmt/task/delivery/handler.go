package delivery

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/dto"
	"github.com/gin-gonic/gin"
)

// Response structures
type response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func writeBindError(c *gin.Context, logger *slog.Logger, err error) {
	logger.Warn("invalid request payload", "error", err)
	c.JSON(http.StatusBadRequest, response{Error: "invalid request payload"})
}

type TaskHandler struct {
	logger  *slog.Logger
	service TaskService
}

func NewTaskHandler(logger *slog.Logger, service TaskService) *TaskHandler {
	return &TaskHandler{
		logger:  logger,
		service: service,
	}
}

// @Tags Tasks
// @Summary List all tasks
// @Description Get a list of all tasks
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response{data=[]dto.TaskResponse}
// @Failure 500 {object} response{error=string}
// @Router /task [get]
func (h *TaskHandler) ListTasks(c *gin.Context) {
	tasks, err := h.service.ListTasks(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list tasks", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: tasks})
}

// @Tags Tasks
// @Summary Get a task by ID
// @Description Get a task by its ID
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} response{data=dto.TaskResponse}
// @Failure 400 {object} response{error=string}
// @Failure 500 {object} response{error=string}
// @Router /task/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid task id"})
		return
	}

	task, err := h.service.GetTask(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get task", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: task})
}

// @Tags Tasks
// @Summary Create a new task
// @Description Create a new task
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param task body dto.CreateTaskRequest true "Task"
// @Success 201 {object} response{data=dto.TaskResponse}
// @Failure 400 {object} response{error=string}
// @Failure 500 {object} response{error=string}
// @Router /task [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&task); err != nil {
		writeBindError(c, h.logger, err)
		return
	}

	created, err := h.service.CreateTask(c.Request.Context(), task)
	if err != nil {
		h.logger.Error("failed to create task", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response{Data: created})
}

// @Tags Tasks
// @Summary Delete a task
// @Description Delete a task by ID
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 204
// @Failure 400 {object} response{error=string}
// @Failure 500 {object} response{error=string}
// @Router /task/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid task id"})
		return
	}

	if err := h.service.DeleteTask(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete task", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// @Tags Tasks
// @Summary Update a task
// @Description Update a task by ID
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body dto.UpdateTaskRequest true "Task data"
// @Success 200 {object} response{data=dto.TaskResponse}
// @Failure 400 {object} response{error=string}
// @Failure 500 {object} response{error=string}
// @Router /task/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid task id"})
		return
	}

	var body dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		writeBindError(c, h.logger, err)
		return
	}

	updated, err := h.service.UpdateTask(c.Request.Context(), id, body)
	if err != nil {
		h.logger.Error("failed to update task", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: updated})
}
