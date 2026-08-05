package delivery

import (
	"errors"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/common/helpers"
	"github.com/Koshsky/erp-backend/internal/common/response"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/dto"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/service"
)

type ProjectHandler struct {
	logger  *slog.Logger
	service ProjectService
}

func NewProjectHandler(logger *slog.Logger, service ProjectService) *ProjectHandler {
	return &ProjectHandler{
		logger:  logger,
		service: service,
	}
}

// ListProjects handles the request to list all projects.
//
//	@Tags			Projects
//	@Summary		List projects
//	@Description	Get a list of all projects
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	response.Response{data=[]dto.ProjectResponse}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/project [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	user, err := helpers.GetUser(c)
	if err != nil {
		response.Unauthorized(c, "authentication required")
		return
	}

	projects, err := h.service.ListProjects(c.Request.Context(), user.ID, user.Role)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, projects)
}

// FindProject handles the request to get a project by ID.
//
//	@Tags			Projects
//	@Summary		Get project
//	@Description	Get a project by ID
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		int	true	"Project ID"
//	@Success		200	{object}	response.Response{data=dto.ProjectResponse}
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/project/{id} [get]
func (h *ProjectHandler) FindProject(c *gin.Context) {
	user, err := helpers.GetUser(c)
	if err != nil {
		response.Unauthorized(c, "authentication required")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid project id")
		return
	}

	project, err := h.service.FindProject(c.Request.Context(), id, user.ID, user.Role)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, project)
}

// CreateProject handles the request to create a new project.
//
//	@Tags			Projects
//	@Summary		Create a new project
//	@Description	Create a new project with the provided details
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			project	body		dto.CreateProjectRequest	true	"Project details"
//	@Success		201		{object}	response.Response{data=dto.ProjectResponse}
//	@Failure		400		{object}	response.Response{data=nil}
//	@Failure		500		{object}	response.Response{data=nil}
//	@Router			/project [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	user, err := helpers.GetUser(c)
	if err != nil {
		response.Unauthorized(c, "authentication required")
		return
	}

	var project dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&project); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	created, err := h.service.CreateProject(c.Request.Context(), project, user.ID, user.Role)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.Created(c, created)
}

// DeleteProject handles the request to delete a project by its ID.
//
//	@Tags			Projects
//	@Summary		Delete project
//	@Description	Delete a project by its ID
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Project ID"
//	@Success		204
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/project/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	user, err := helpers.GetUser(c)
	if err != nil {
		response.Unauthorized(c, "authentication required")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid project id")
		return
	}

	if err = h.service.DeleteProject(c.Request.Context(), id, user.ID, user.Role); err != nil {
		h.handleError(c, err)
		return
	}
	response.NoContent(c)
}

// UpdateProject handles the request to update a project by its ID.
//
//	@Tags			Projects
//	@Summary		Update project
//	@Description	Update a project by its ID
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Project ID"
//	@Param			project	body		dto.UpdateProjectRequest	true	"Project data"
//	@Success		200		{object}	response.Response{data=dto.ProjectResponse}
//	@Failure		400		{object}	response.Response{data=nil}
//	@Failure		500		{object}	response.Response{data=nil}
//	@Router			/project/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	user, err := helpers.GetUser(c)
	if err != nil {
		response.Unauthorized(c, "authentication required")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid project id")
		return
	}

	body := dto.UpdateProjectRequest{}
	if err = c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	project, err := h.service.UpdateProject(c.Request.Context(), id, body, user.ID, user.Role)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, project)
}

// handleError маппит доменные ошибки сервиса проектов в HTTP-статусы.
func (h *ProjectHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrForbidden):
		response.Forbidden(c, "no rights to perform this operation on the project")
	case errors.Is(err, service.ErrNotFound):
		response.NotFound(c, "project not found")
	default:
		response.InternalError(c, h.logger, err.Error(), err)
	}
}
