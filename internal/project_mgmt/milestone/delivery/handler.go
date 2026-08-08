package delivery

import (
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/common/response"
	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/dto"
)

type MilestoneHandler struct {
	logger  *slog.Logger
	service MilestoneService
	mw      *rbac.Middleware
}

func NewMilestoneHandler(logger *slog.Logger, service MilestoneService, mw *rbac.Middleware) *MilestoneHandler {
	return &MilestoneHandler{
		logger:  logger,
		service: service,
		mw:      mw,
	}
}

// ListMilestones handles the request to list all milestones.
//
//	@Tags			Milestones
//	@Summary		List milestones
//	@Description	List all milestones
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	response.Response{data=[]dto.MilestoneResponse}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/milestone [get]
func (h *MilestoneHandler) ListMilestones(c *gin.Context) {
	milestones, err := h.service.ListMilestones(c.Request.Context())
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, milestones)
}

// FindMilestone handles the request to get a milestone by id.
//
//	@Tags			Milestones
//	@Summary		Get milestone
//	@Description	Get milestone by id
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		int	true	"Milestone ID"
//	@Success		200	{object}	response.Response{data=dto.MilestoneResponse}
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/milestone/{id} [get]
func (h *MilestoneHandler) FindMilestone(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid milestone id")
		return
	}

	milestone, err := h.service.FindMilestone(c.Request.Context(), id)
	if err != nil {
		response.HandleError(c, h.logger, err)
		return
	}
	response.OK(c, milestone)
}

// CreateMilestone handles the request to create a milestone.
//
//	@Tags			Milestones
//	@Summary		Create milestone
//	@Description	Create milestone with the input payload
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateMilestoneRequest	true	"Milestone data"
//	@Success		201		{object}	response.Response{data=dto.MilestoneResponse}
//	@Failure		400		{object}	response.Response{data=nil}
//	@Failure		500		{object}	response.Response{data=nil}
//	@Router			/milestone [post]
func (h *MilestoneHandler) CreateMilestone(c *gin.Context) {
	var milestone dto.CreateMilestoneRequest
	if err := c.ShouldBindJSON(&milestone); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	created, err := h.service.CreateMilestone(c.Request.Context(), milestone)
	if err != nil {
		response.HandleError(c, h.logger, err)
		return
	}
	response.Created(c, created)
}

// DeleteMilestone handles the request to delete a milestone.
//
//	@Tags			Milestones
//	@Summary		Delete milestone
//	@Description	Delete milestone by ID
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path	int	true	"Milestone ID"
//	@Success		204
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/milestone/{id} [delete]
func (h *MilestoneHandler) DeleteMilestone(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid milestone id")
		return
	}

	if err = h.service.DeleteMilestone(c.Request.Context(), id); err != nil {
		response.HandleError(c, h.logger, err)
		return
	}
	response.NoContent(c)
}

// UpdateMilestone handles the request to update a milestone.
//
//	@Tags			Milestones
//	@Summary		Update milestone
//	@Description	Update milestone by ID
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		int							true	"Milestone ID"
//	@Param			body	body		dto.UpdateMilestoneRequest	true	"Milestone data"
//	@Success		200		{object}	response.Response{data=dto.MilestoneResponse}
//	@Failure		400		{object}	response.Response{data=nil}
//	@Failure		500		{object}	response.Response{data=nil}
//	@Router			/milestone/{id} [put]
func (h *MilestoneHandler) UpdateMilestone(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid milestone id")
		return
	}

	body := dto.UpdateMilestoneRequest{}
	if err = c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updated, err := h.service.UpdateMilestone(c.Request.Context(), id, body)
	if err != nil {
		response.HandleError(c, h.logger, err)
		return
	}
	response.OK(c, updated)
}
