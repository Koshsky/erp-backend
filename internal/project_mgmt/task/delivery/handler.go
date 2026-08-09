package delivery

import (
	"log/slog"
	"strconv"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/service"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/dto"
	"github.com/Koshsky/erp-backend/internal/response"
)

type TaskHandler struct {
	logger  *slog.Logger
	service TaskService
	mw      *rbac.Middleware
}

// NewTaskHandler builds the TaskHandler handler.
func NewTaskHandler(logger *slog.Logger, svc *service.TaskService, mw *rbac.Middleware) *TaskHandler {
	return &TaskHandler{
		logger:  logger,
		service: svc,
		mw:      mw,
	}
}

// ListTasks handles the request to list all tasks.
//
//	@Tags			Tasks
//	@Summary		List all tasks
//	@Description	Get a list of all tasks
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	response.Response{data=[]dto.TaskResponse}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/task [get]
func (h *TaskHandler) ListTasks(c *gin.Context) {
	tasks, err := h.service.ListTasks(c.Request.Context())
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, tasks)
}

// FindTask handles the request to get a task by ID.
//
//	@Tags			Tasks
//	@Summary		Get a task by ID
//	@Description	Get a task by its ID
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		int	true	"Task ID"
//	@Success		200	{object}	response.Response{data=dto.TaskResponse}
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/task/{id} [get]
func (h *TaskHandler) FindTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	task, err := h.service.FindTask(c.Request.Context(), id)
	if err != nil {
		response.Error(c, h.logger, err)
		return
	}
	response.OK(c, task)
}

// CreateTask handles the request to create a new task.
//
//	@Tags			Tasks
//	@Summary		Create a new task
//	@Description	Create a new task
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			task	body		dto.CreateTaskRequest	true	"Task"
//	@Success		201		{object}	response.Response{data=dto.TaskResponse}
//	@Failure		400		{object}	response.Response{data=nil}
//	@Failure		500		{object}	response.Response{data=nil}
//	@Router			/task [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&task); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	created, err := h.service.CreateTask(c.Request.Context(), task)
	if err != nil {
		response.Error(c, h.logger, err)
		return
	}
	response.Created(c, created)
}

// DeleteTask handles the request to delete a task.
//
//	@Tags			Tasks
//	@Summary		Delete a task
//	@Description	Delete a task by ID
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Task ID"
//	@Success		204
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/task/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	if err = h.service.DeleteTask(c.Request.Context(), id); err != nil {
		response.Error(c, h.logger, err)
		return
	}
	response.NoContent(c)
}

// UpdateTask handles the request to update a task.
//
//	@Tags			Tasks
//	@Summary		Update a task
//	@Description	Update a task by ID
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Task ID"
//	@Param			task	body		dto.UpdateTaskRequest	true	"Task data"
//	@Success		200		{object}	response.Response{data=dto.TaskResponse}
//	@Failure		400		{object}	response.Response{data=nil}
//	@Failure		500		{object}	response.Response{data=nil}
//	@Router			/task/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	var body dto.UpdateTaskRequest
	if err = c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updated, err := h.service.UpdateTask(c.Request.Context(), id, body)
	if err != nil {
		response.Error(c, h.logger, err)
		return
	}
	response.OK(c, updated)
}
