package delivery

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Koshsky/erp/api/internal/task/dto"
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

func (h *TaskHandler) ListTasks(c *gin.Context) {
	tasks, err := h.service.ListTasks(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list tasks", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: tasks})
}

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
