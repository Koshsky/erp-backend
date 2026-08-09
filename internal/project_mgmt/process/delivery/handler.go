package delivery

import (
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/dto"
	"github.com/Koshsky/erp-backend/internal/response"
)

type ProcessHandler struct {
	logger  *slog.Logger
	service ProcessService
	mw      *rbac.Middleware
}

// ListProcesses handles the request to list all processes.
//
//	@Tags			Processes
//	@Summary		List processes
//	@Description	Get a list of all processes
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	response.Response{data=[]dto.ProcessResponse}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/process [get]
func (h *ProcessHandler) ListProcesses(c *gin.Context) {
	processes, err := h.service.ListProcesses(c.Request.Context())
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, processes)
}

// FindProcess handles the request to get a process by ID.
//
//	@Tags			Processes
//	@Summary		Get process
//	@Description	Get a process by ID
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Process ID"
//	@Success		200	{object}	response.Response{data=dto.ProcessResponse}
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/process/{id} [get]
func (h *ProcessHandler) FindProcess(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid process id")
		return
	}

	process, err := h.service.FindProcess(c.Request.Context(), id)
	if err != nil {
		response.Error(c, h.logger, err)
		return
	}
	response.OK(c, process)
}

// CreateProcess handles the request to create a new process.
//
//	@Tags			Processes
//	@Summary		Create process
//	@Description	Create a new process
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			process	body		dto.CreateProcessRequest	true	"Process data"
//	@Success		201		{object}	response.Response{data=dto.ProcessResponse}
//	@Failure		400		{object}	response.Response{data=nil}
//	@Failure		500		{object}	response.Response{data=nil}
//	@Router			/process [post]
func (h *ProcessHandler) CreateProcess(c *gin.Context) {
	var process dto.CreateProcessRequest
	if err := c.ShouldBindJSON(&process); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	created, err := h.service.CreateProcess(c.Request.Context(), process)
	if err != nil {
		response.Error(c, h.logger, err)
		return
	}
	response.Created(c, created)
}

// DeleteProcess handles the request to delete a process.
//
//	@Tags			Processes
//	@Summary		Delete process
//	@Description	Delete a process by ID
//	@Security		ApiKeyAuth
//	@Param			id	path	int	true	"Process ID"
//	@Success		204
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/process/{id} [delete]
func (h *ProcessHandler) DeleteProcess(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid process id")
		return
	}

	if err = h.service.DeleteProcess(c.Request.Context(), id); err != nil {
		response.Error(c, h.logger, err)
		return
	}
	response.NoContent(c)
}

// UpdateProcess handles the request to update a process.
//
//	@Tags			Processes
//	@Summary		Update process
//	@Description	Update a process by ID
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Process ID"
//	@Param			body	body		dto.UpdateProcessRequest	true	"Process data"
//	@Success		200		{object}	response.Response{data=dto.ProcessResponse}
//	@Failure		400		{object}	response.Response{data=nil}
//	@Failure		500		{object}	response.Response{data=nil}
//	@Router			/process/{id} [put]
func (h *ProcessHandler) UpdateProcess(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid process id")
		return
	}

	body := dto.UpdateProcessRequest{}
	if err = c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	process, err := h.service.UpdateProcess(c.Request.Context(), id, body)
	if err != nil {
		response.Error(c, h.logger, err)
		return
	}
	response.OK(c, process)
}
