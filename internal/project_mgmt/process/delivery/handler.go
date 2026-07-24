package delivery

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/dto"
	"github.com/gin-gonic/gin"
)

// TODO: move response to common DTO

type response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func writeBindError(c *gin.Context, logger *slog.Logger, err error) {
	logger.Warn("invalid request payload", "error", err)
	c.JSON(http.StatusBadRequest, response{Error: "invalid request payload"})
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

func (h *ProcessHandler) ListProcesses(c *gin.Context) {
	processes, err := h.service.ListProcesses(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list processes", "error", err)
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

	process, err := h.service.GetProcess(c, id)
	if err != nil {
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

	created, err := h.service.CreateProcess(c, process)
	if err != nil {
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

	if err := h.service.DeleteProcess(c, id); err != nil {
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

	process, err := h.service.UpdateProcess(c, id, body)
	if err != nil {
		h.logger.Error("failed to update process", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: process})
}
