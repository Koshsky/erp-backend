package delivery

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/dto"
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

// @Tags Projects
// @Summary List projects
// @Description Get a list of all projects
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response{data=[]dto.ProjectResponse}
// @Failure 500 {object} response{error=string}
// @Router /project [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	projects, err := h.service.ListProjects(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list projects", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: projects})
}

// @Tags Projects
// @Summary Get project
// @Description Get a project by ID
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} response{data=dto.ProjectResponse}
// @Failure 400 {object} response{error=string}
// @Failure 500 {object} response{error=string}
// @Router /project/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid project id"})
		return
	}

	project, err := h.service.GetProject(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get project", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: project})
}

// @Tags Projects
// @Summary Create a new project
// @Description Create a new project with the provided details
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param project body dto.CreateProjectRequest true "Project details"
// @Success 201 {object} response{data=dto.ProjectResponse}
// @Failure 400 {object} response{error=string}
// @Failure 500 {object} response{error=string}
// @Router /project [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var project dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&project); err != nil {
		writeBindError(c, h.logger, err)
		return
	}

	created, err := h.service.CreateProject(c.Request.Context(), project)
	if err != nil {
		h.logger.Error("failed to create project", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response{Data: created})
}

// @Tags Projects
// @Summary Delete project
// @Description Delete a project by its ID
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 204
// @Failure 400 {object} response{error=string}
// @Failure 500 {object} response{error=string}
// @Router /project/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid project id"})
		return
	}

	if err := h.service.DeleteProject(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete project", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// @Tags Projects
// @Summary Update project
// @Description Update a project by its ID
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param project body dto.UpdateProjectRequest true "Project data"
// @Success 200 {object} response{data=dto.ProjectResponse}
// @Failure 400 {object} response{error=string}
// @Failure 500 {object} response{error=string}
// @Router /project/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid project id"})
		return
	}

	body := dto.UpdateProjectRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		writeBindError(c, h.logger, err)
		return
	}

	project, err := h.service.UpdateProject(c.Request.Context(), id, body)
	if err != nil {
		h.logger.Error("failed to update project", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: project})
}
