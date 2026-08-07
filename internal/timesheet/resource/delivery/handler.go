package delivery

import (
	"errors"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/common/helpers"
	"github.com/Koshsky/erp-backend/internal/common/response"
	"github.com/Koshsky/erp-backend/internal/timesheet/resource/dto"
	"github.com/Koshsky/erp-backend/internal/timesheet/resource/service"
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

// handleError маппит доменные ошибки сервиса ресурсов в HTTP-статусы.
func (h *ResourceHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrForbidden):
		response.Forbidden(c, err.Error())
	case errors.Is(err, service.ErrNotFound):
		response.NotFound(c, err.Error())
	default:
		response.InternalError(c, h.logger, err.Error(), err)
	}
}

// ListResources handles the request to list all resources.
//
//	@Tags			TimesheetResources
//	@Summary		List resources
//	@Description	List all resources
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	response.Response{data=[]dto.ResourceResponse}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/timesheet/resources [get]
func (h *ResourceHandler) ListResources(c *gin.Context) {
	user, err := helpers.GetUser(c)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}

	resources, err := h.service.ListResources(c.Request.Context(), user.ID, user.Role)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, resources)
}

// FindResource handles the request to get a resource by id.
//
//	@Tags			TimesheetResources
//	@Summary		Get resource
//	@Description	Get resource by id
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		int	true	"Resource ID"
//	@Success		200	{object}	response.Response{data=dto.ResourceResponse}
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/timesheet/resources/{id} [get]
func (h *ResourceHandler) FindResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid resource id")
		return
	}

	user, err := helpers.GetUser(c)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}

	resource, err := h.service.FindResource(c.Request.Context(), id, user.ID, user.Role)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, resource)
}

// CreateResource handles the request to create a resource.
//
//	@Tags			TimesheetResources
//	@Summary		Create resource
//	@Description	Create resource
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			resource	body		dto.CreateResourceRequest	true	"Resource"
//	@Success		201			{object}	response.Response{data=dto.ResourceResponse}
//	@Failure		400			{object}	response.Response{data=nil}
//	@Failure		500			{object}	response.Response{data=nil}
//	@Router			/timesheet/resources [post]
func (h *ResourceHandler) CreateResource(c *gin.Context) {
	var resource dto.CreateResourceRequest
	if err := c.ShouldBindJSON(&resource); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := helpers.GetUser(c)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}

	created, err := h.service.CreateResource(c.Request.Context(), resource, user.ID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.Created(c, created)
}

// DeleteResource handles the request to delete resource by id.
//
//	@Tags			TimesheetResources
//	@Summary		Delete resource
//	@Description	Delete resource by id
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path	int	true	"Resource ID"
//	@Success		204
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/timesheet/resources/{id} [delete]
func (h *ResourceHandler) DeleteResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid resource id")
		return
	}

	user, err := helpers.GetUser(c)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}

	if err = h.service.DeleteResource(c.Request.Context(), id, user.ID, user.Role); err != nil {
		h.handleError(c, err)
		return
	}
	response.NoContent(c)
}

// UpdateResource handles the request to update resource by id.
//
//	@Tags			TimesheetResources
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
//	@Router			/timesheet/resources/{id} [put]
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

	user, err := helpers.GetUser(c)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}

	updated, err := h.service.UpdateResource(c.Request.Context(), id, body, user.ID, user.Role)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, updated)
}
