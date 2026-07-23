package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Koshsky/erp/api/internal/dto"
	"github.com/gin-gonic/gin"
)

type TaskService interface {
	GetTask(ctx context.Context, id int64) (*dto.TaskResponse, error)
	CreateTask(ctx context.Context, task dto.CreateTaskRequest) (*dto.TaskResponse, error)
	DeleteTask(ctx context.Context, id int64) error
	UpdateTask(ctx context.Context, id int64, task dto.UpdateTaskRequest) (*dto.TaskResponse, error)
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

func (h *TaskHandler) GetTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid task id"})
		return
	}

	task, err := h.service.GetTask(c.Request.Context(), id)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, response{Error: "task not found"})
			return
		}
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
		if isValidationError(err) {
			c.JSON(http.StatusBadRequest, response{Error: err.Error()})
			return
		}
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
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, response{Error: "task not found"})
			return
		}
		if isValidationError(err) {
			c.JSON(http.StatusBadRequest, response{Error: err.Error()})
			return
		}
		h.logger.Error("failed to update task", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: updated})
}
