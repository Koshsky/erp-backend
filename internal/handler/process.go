package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Koshsky/erp/api/internal/dto"
	"github.com/gin-gonic/gin"
)

type ProcessService interface {
	GetProcess(ctx context.Context, id int64) (*dto.ProcessResponse, error)
	CreateProcess(ctx context.Context, process dto.CreateProcessRequest) (*dto.ProcessResponse, error)
	DeleteProcess(ctx context.Context, id int64) error
	UpdateProcess(ctx context.Context, id int64, process dto.UpdateProcessRequest) (*dto.ProcessResponse, error)
	ListProcesses(ctx context.Context, projectID int64) ([]dto.ProcessResponse, error)
	GetDetailedProcess(ctx context.Context, id int64) (*dto.ProcessDetailResponse, error)
}

type ProcessHandler struct {
	logger  *slog.Logger
	service ProcessService
}

func NewProcessHandler(logger *slog.Logger, service ProcessService) *ProcessHandler {
	return &ProcessHandler{
		logger:  logger,
		service: service,
	}
}

func (h *ProcessHandler) GetDetailedProcess(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid process id",
		})
		return
	}

	process, err := h.service.GetDetailedProcess(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, process)
}

func (h *ProcessHandler) ListProcesses(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid project id"})
		return
	}

	processes, err := h.service.ListProcesses(c.Request.Context(), projectID)
	if err != nil {
		h.logger.Error("failed to list processes", "projectID", projectID, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: processes})
}

func (h *ProcessHandler) GetProcess(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid process id"})
		return
	}

	process, err := h.service.GetProcess(c.Request.Context(), id)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, response{Error: "process not found"})
			return
		}
		h.logger.Error("failed to get process", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: process})
}

func (h *ProcessHandler) CreateProcess(c *gin.Context) {
	var process dto.CreateProcessRequest
	if err := c.ShouldBindJSON(&process); err != nil {
		writeBindError(c, h.logger, err)
		return
	}

	created, err := h.service.CreateProcess(c.Request.Context(), process)
	if err != nil {
		if isValidationError(err) {
			c.JSON(http.StatusBadRequest, response{Error: err.Error()})
			return
		}
		h.logger.Error("failed to create process", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response{Data: created})
}

func (h *ProcessHandler) DeleteProcess(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid process id"})
		return
	}

	if err := h.service.DeleteProcess(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete process", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *ProcessHandler) UpdateProcess(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid process id"})
		return
	}

	body := dto.UpdateProcessRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		writeBindError(c, h.logger, err)
		return
	}

	process, err := h.service.UpdateProcess(c.Request.Context(), id, body)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, response{Error: "process not found"})
			return
		}
		if isValidationError(err) {
			c.JSON(http.StatusBadRequest, response{Error: err.Error()})
			return
		}
		h.logger.Error("failed to update process", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: process})
}
