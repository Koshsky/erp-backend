package delivery

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Koshsky/erp-backend/internal/resource/dto"
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
	c.JSON(http.StatusOK, resources)
}

func (h *ResourceHandler) GetResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid resource id"})
		return
	}

	resource, err := h.service.GetResource(c.Request.Context(), id)
	if err != nil {
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
		h.logger.Error("failed to update resource", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: updated})
}
