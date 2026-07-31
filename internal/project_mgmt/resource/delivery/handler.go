package delivery

import (
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/common/response"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/resource/dto"
)

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

// ListResources handles the request to list all resources.
//
//	@Tags			Resources
//	@Summary		List resources
//	@Description	List all resources
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	response.Response{data=[]dto.ResourceResponse}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/resource [get]
func (h *ResourceHandler) ListResources(c *gin.Context) {
	resources, err := h.service.ListResources(c.Request.Context())
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, resources)
}

// FindResource handles the request to get a resource by id.
//
//	@Tags			Resources
//	@Summary		Get resource
//	@Description	Get resource by id
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		int	true	"Resource ID"
//	@Success		200	{object}	response.Response{data=dto.ResourceResponse}
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/resource/{id} [get]
func (h *ResourceHandler) FindResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid resource id")
		return
	}

	resource, err := h.service.FindResource(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, resource)
}

// CreateResource handles the request to create a resource.
//
//	@Tags			Resources
//	@Summary		Create resource
//	@Description	Create resource
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			resource	body		dto.CreateResourceRequest	true	"Resource"
//	@Success		201			{object}	response.Response{data=dto.ResourceResponse}
//	@Failure		400			{object}	response.Response{data=nil}
//	@Failure		500			{object}	response.Response{data=nil}
//	@Router			/resource [post]
func (h *ResourceHandler) CreateResource(c *gin.Context) {
	var resource dto.CreateResourceRequest
	if err := c.ShouldBindJSON(&resource); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	created, err := h.service.CreateResource(c.Request.Context(), resource)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.Created(c, created)
}

// DeleteResource handles the request to delete resource by id.
//
//	@Tags			Resources
//	@Summary		Delete resource
//	@Description	Delete resource by id
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path	int	true	"Resource ID"
//	@Success		204
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/resource/{id} [delete]
func (h *ResourceHandler) DeleteResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid resource id")
		return
	}

	if err = h.service.DeleteResource(c.Request.Context(), id); err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.NoContent(c)
}

// UpdateResource handles the request to update resource by id.
//
//	@Tags			Resources
//	@Summary		Update resource
//	@Description	Update resource by id
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int							true	"Resource ID"
//	@Param			resource	body		dto.UpdateResourceRequest	true	"Resource"
//	@Success		200			{object}	response.Response{data=dto.ResourceResponse}
//	@Failure		400			{object}	response.Response{data=nil}
//	@Failure		500			{object}	response.Response{data=nil}
//	@Router			/resource/{id} [put]
func (h *ResourceHandler) UpdateResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid resource id")
		return
	}

	body := dto.UpdateResourceRequest{}
	if err = c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updated, err := h.service.UpdateResource(c.Request.Context(), id, body)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, updated)
}
