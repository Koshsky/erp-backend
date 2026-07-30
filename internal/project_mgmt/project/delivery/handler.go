package delivery

import (
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/common/response"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/dto"
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
//	@Failure		500	{object}	response.Response
//	@Router			/project [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	projects, err := h.service.ListProjects(c.Request.Context())
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
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
//	@Failure		400	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/project/{id} [get]
func (h *ProjectHandler) FindProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid project id")
		return
	}

	project, err := h.service.FindProject(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
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
//	@Failure		400		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/project [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var project dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&project); err != nil {
		response.HandleBindError(c, h.logger, err)
		return
	}

	created, err := h.service.CreateProject(c.Request.Context(), project)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
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
//	@Failure		400	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/project/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid project id")
		return
	}

	if err = h.service.DeleteProject(c.Request.Context(), id); err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
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
//	@Failure		400		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/project/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid project id")
		return
	}

	body := dto.UpdateProjectRequest{}
	if err = c.ShouldBindJSON(&body); err != nil {
		response.HandleBindError(c, h.logger, err)
		return
	}

	project, err := h.service.UpdateProject(c.Request.Context(), id, body)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, project)
}
