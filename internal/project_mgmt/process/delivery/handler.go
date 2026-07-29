package delivery

import (
	"log/slog"
	"strconv"

	"github.com/Koshsky/erp-backend/internal/common/response"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/dto"
	"github.com/gin-gonic/gin"
)

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
// @Success 200 {object} response.Response{data=[]dto.ProcessResponse}
// @Failure 500 {object} response.Response
// @Router /process [get]
func (h *ProcessHandler) ListProcesses(c *gin.Context) {
	processes, err := h.service.ListProcesses(c.Request.Context())
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, processes)
}

// @Tags Processes
// @Summary Get process
// @Description Get a process by ID
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Process ID"
// @Success 200 {object} response.Response{data=dto.ProcessResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /process/{id} [get]
func (h *ProcessHandler) GetProcess(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid process id")
		return
	}

	process, err := h.service.GetProcess(c, id)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, process)
}

// @Tags Processes
// @Summary Create process
// @Description Create a new process
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param process body dto.CreateProcessRequest true "Process data"
// @Success 201 {object} response.Response{data=dto.ProcessResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /process [post]
func (h *ProcessHandler) CreateProcess(c *gin.Context) {
	var process dto.CreateProcessRequest
	if err := c.ShouldBindJSON(&process); err != nil {
		response.HandleBindError(c, h.logger, err)
		return
	}

	created, err := h.service.CreateProcess(c, process)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.Created(c, created)
}

// @Tags Processes
// @Summary Delete process
// @Description Delete a process by ID
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Process ID"
// @Success 204
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /process/{id} [delete]
func (h *ProcessHandler) DeleteProcess(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid process id")
		return
	}

	if err := h.service.DeleteProcess(c, id); err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.NoContent(c)
}

// @Tags Processes
// @Summary Update process
// @Description Update a process by ID
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Process ID"
// @Param body body dto.UpdateProcessRequest true "Process data"
// @Success 200 {object} response.Response{data=dto.ProcessResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /process/{id} [put]
func (h *ProcessHandler) UpdateProcess(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid process id")
		return
	}

	body := dto.UpdateProcessRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.HandleBindError(c, h.logger, err)
		return
	}

	process, err := h.service.UpdateProcess(c, id, body)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, process)
}
