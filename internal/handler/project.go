package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Koshsky/erp/api/internal/dto"
	"github.com/gin-gonic/gin"
)

type ProjectService interface {
	ListProjects(ctx context.Context) ([]dto.ProjectResponse, error)
	GetProject(ctx context.Context, id int64) (*dto.ProjectResponse, error)
	CreateProject(ctx context.Context, project dto.CreateProjectRequest) (*dto.ProjectResponse, error)
	DeleteProject(ctx context.Context, id int64) error
	UpdateProject(ctx context.Context, id int64, project dto.UpdateProjectRequest) (*dto.ProjectResponse, error)
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

func (h *ProjectHandler) ListProjects(c *gin.Context) {
	projects, err := h.service.ListProjects(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list projects", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: projects})
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid project id"})
		return
	}

	project, err := h.service.GetProject(c.Request.Context(), id)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, response{Error: "project not found"})
			return
		}
		h.logger.Error("failed to get project", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: project})
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var project dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&project); err != nil {
		writeBindError(c, h.logger, err)
		return
	}

	created, err := h.service.CreateProject(c.Request.Context(), project)
	if err != nil {
		if isValidationError(err) {
			c.JSON(http.StatusBadRequest, response{Error: err.Error()})
			return
		}
		h.logger.Error("failed to create project", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response{Data: created})
}

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
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, response{Error: "project not found"})
			return
		}
		if isValidationError(err) {
			c.JSON(http.StatusBadRequest, response{Error: err.Error()})
			return
		}
		h.logger.Error("failed to update project", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: project})
}
