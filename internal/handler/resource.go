package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Koshsky/erp/api/internal/dto"
	"github.com/gin-gonic/gin"
)

type ResourceService interface {
	ListResources(ctx context.Context) ([]dto.ResourceResponse, error)
	GetResource(ctx context.Context, id int64) (*dto.ResourceResponse, error)
	CreateResource(ctx context.Context, resource dto.CreateResourceRequest) (*dto.ResourceResponse, error)
	DeleteResource(ctx context.Context, id int64) error
	UpdateResource(ctx context.Context, id int64, resource dto.UpdateResourceRequest) (*dto.ResourceResponse, error)
	GetResourceUsage(ctx context.Context, date time.Time) ([]dto.ResourceUsageResponse, error)
}

type ResourceHandler struct {
	logger  *slog.Logger
	service ResourceService
}

func NewResourceHandler(logger *slog.Logger, service ResourceService) *ResourceHandler {
	return &ResourceHandler{
		logger:  logger,
		service: service,
	}
}

func (h *ResourceHandler) ListResources(c *gin.Context) {
	resources, err := h.service.ListResources(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list resources", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: resources})
}

func (h *ResourceHandler) GetResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid resource id"})
		return
	}

	resource, err := h.service.GetResource(c.Request.Context(), id)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, response{Error: "resource not found"})
			return
		}
		h.logger.Error("failed to get resource", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: resource})
}

func (h *ResourceHandler) CreateResource(c *gin.Context) {
	var resource dto.CreateResourceRequest
	if err := c.ShouldBindJSON(&resource); err != nil {
		writeBindError(c, h.logger, err)
		return
	}

	created, err := h.service.CreateResource(c.Request.Context(), resource)
	if err != nil {
		if isValidationError(err) {
			c.JSON(http.StatusBadRequest, response{Error: err.Error()})
			return
		}
		h.logger.Error("failed to create resource", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response{Data: created})
}

func (h *ResourceHandler) DeleteResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid resource id"})
		return
	}

	if err := h.service.DeleteResource(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete resource", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *ResourceHandler) UpdateResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid resource id"})
		return
	}

	body := dto.UpdateResourceRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		writeBindError(c, h.logger, err)
		return
	}

	updated, err := h.service.UpdateResource(c.Request.Context(), id, body)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, response{Error: "resource not found"})
			return
		}
		if isValidationError(err) {
			c.JSON(http.StatusBadRequest, response{Error: err.Error()})
			return
		}
		h.logger.Error("failed to update resource", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: updated})
}

func (h *ResourceHandler) GetResourceUsage(c *gin.Context) {
	dateRaw := c.Query("date")
	if dateRaw == "" {
		c.JSON(http.StatusBadRequest, response{Error: "query parameter 'date' is required (YYYY-MM-DD)"})
		return
	}

	targetDate, err := time.Parse("2006-01-02", dateRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid date format, expected YYYY-MM-DD"})
		return
	}

	usage, err := h.service.GetResourceUsage(c.Request.Context(), targetDate)
	if err != nil {
		h.logger.Error("failed to get resource usage", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: usage})
}
