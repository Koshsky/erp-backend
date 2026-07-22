package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Koshsky/erp/api/internal/dto"
	"github.com/gin-gonic/gin"
)

type MilestoneService interface {
	ListMilestonesByProcessID(ctx context.Context, processID int64) ([]dto.MilestoneResponse, error)
	GetMilestone(ctx context.Context, id int64) (*dto.MilestoneResponse, error)
	CreateMilestone(ctx context.Context, milestone dto.CreateMilestoneRequest) (*dto.MilestoneResponse, error)
	DeleteMilestone(ctx context.Context, id int64) error
	UpdateMilestone(ctx context.Context, id int64, milestone dto.UpdateMilestoneRequest) (*dto.MilestoneResponse, error)
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

func (h *MilestoneHandler) ListMilestonesByProcess(c *gin.Context) {
	processID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid process id"})
		return
	}

	milestones, err := h.service.ListMilestonesByProcessID(c.Request.Context(), processID)
	if err != nil {
		h.logger.Error("failed to list milestones", "processID", processID, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: milestones})
}

func (h *MilestoneHandler) GetMilestone(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid milestone id"})
		return
	}

	milestone, err := h.service.GetMilestone(c.Request.Context(), id)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, response{Error: "milestone not found"})
			return
		}
		h.logger.Error("failed to get milestone", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: milestone})
}

func (h *MilestoneHandler) CreateMilestone(c *gin.Context) {
	var milestone dto.CreateMilestoneRequest
	if err := c.ShouldBindJSON(&milestone); err != nil {
		writeBindError(c, h.logger, err)
		return
	}

	created, err := h.service.CreateMilestone(c.Request.Context(), milestone)
	if err != nil {
		if isValidationError(err) {
			c.JSON(http.StatusBadRequest, response{Error: err.Error()})
			return
		}
		h.logger.Error("failed to create milestone", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response{Data: created})
}

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
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, response{Error: "milestone not found"})
			return
		}
		if isValidationError(err) {
			c.JSON(http.StatusBadRequest, response{Error: err.Error()})
			return
		}
		h.logger.Error("failed to update milestone", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: updated})
}
