package delivery

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/dto"
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

type MilestoneHandler struct {
	logger  *slog.Logger
	service MilestoneService
}

func NewMilestoneHandler(logger *slog.Logger, service MilestoneService) *MilestoneHandler {
	return &MilestoneHandler{
		logger:  logger,
		service: service,
	}
}

// @Tags Milestones
// @Summary List milestones
// @Description List all milestones
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response{data=[]dto.MilestoneResponse}
// @Failure 500 {object} response
// @Router /milestone [get]
func (h *MilestoneHandler) ListMilestones(c *gin.Context) {
	milestones, err := h.service.ListMilestones(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list milestones", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: milestones})
}

// @Tags Milestones
// @Summary Get milestone
// @Description Get milestone by id
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Milestone ID"
// @Success 200 {object} response{data=dto.MilestoneResponse}
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /milestone/{id} [get]
func (h *MilestoneHandler) GetMilestone(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid milestone id"})
		return
	}

	milestone, err := h.service.GetMilestone(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get milestone", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: milestone})
}

// @Tags Milestones
// @Summary Create milestone
// @Description Create milestone with the input payload
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.CreateMilestoneRequest true "Milestone data"
// @Success 201 {object} response{data=dto.MilestoneResponse}
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /milestone [post]
func (h *MilestoneHandler) CreateMilestone(c *gin.Context) {
	var milestone dto.CreateMilestoneRequest
	if err := c.ShouldBindJSON(&milestone); err != nil {
		writeBindError(c, h.logger, err)
		return
	}

	created, err := h.service.CreateMilestone(c.Request.Context(), milestone)
	if err != nil {
		h.logger.Error("failed to create milestone", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response{Data: created})
}

// @Tags Milestones
// @Summary Delete milestone
// @Description Delete milestone by ID
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Milestone ID"
// @Success 204
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /milestone/{id} [delete]
func (h *MilestoneHandler) DeleteMilestone(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid milestone id"})
		return
	}

	if err := h.service.DeleteMilestone(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete milestone", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// @Tags Milestones
// @Summary Update milestone
// @Description Update milestone by ID
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Milestone ID"
// @Param body body dto.UpdateMilestoneRequest true "Milestone data"
// @Success 200 {object} response{data=dto.MilestoneResponse}
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /milestone/{id} [put]
func (h *MilestoneHandler) UpdateMilestone(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid milestone id"})
		return
	}

	body := dto.UpdateMilestoneRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		writeBindError(c, h.logger, err)
		return
	}

	updated, err := h.service.UpdateMilestone(c.Request.Context(), id, body)
	if err != nil {
		h.logger.Error("failed to update milestone", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: updated})
}
