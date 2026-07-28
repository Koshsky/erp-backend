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

// @Tags Processes
// @Summary List processes
// @Description Get a list of all processes
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response{data=[]dto.ProcessResponse}
// @Failure 500 {object} response
// @Router /process [get]
func (h *ProcessHandler) ListProcesses(c *gin.Context) {
	processes, err := h.service.ListProcesses(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list processes", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: processes})
}

// @Tags Processes
// @Summary Get process
// @Description Get a process by ID
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Process ID"
// @Success 200 {object} response{data=dto.ProcessResponse}
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /process/{id} [get]
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

// @Tags Processes
// @Summary Create process
// @Description Create a new process
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param process body dto.CreateProcessRequest true "Process data"
// @Success 201 {object} response{data=dto.ProcessResponse}
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /process [post]
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

// @Tags Processes
// @Summary Delete process
// @Description Delete a process by ID
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Process ID"
// @Success 204
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /process/{id} [delete]
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

// @Tags Processes
// @Summary Update process
// @Description Update a process by ID
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Process ID"
// @Param body body dto.UpdateProcessRequest true "Process data"
// @Success 200 {object} response{data=dto.ProcessResponse}
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /process/{id} [put]
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
