package delivery

import (
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/common/ctx"
	"github.com/Koshsky/erp-backend/internal/common/response"
	"github.com/gin-gonic/gin"
)

type SchedulingHandler struct {
	logger  *slog.Logger
	service MilestoneService
}

func NewSchedulingHandler(logger *slog.Logger, service MilestoneService) *SchedulingHandler {
	return &SchedulingHandler{
		logger:  logger,
		service: service,
	}
}

// @Tags Scheduling
// @Summary Get project scheduling
// @Description Get project scheduling (project portfolio)
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /scheduling/projects [get]
func (h *SchedulingHandler) GetProjectScheduling(c *gin.Context) {
	userID := ctx.GetUserID(c)
	role := ctx.GetRole(c)
	scheduling, err := h.service.GetProjectScheduling(c.Request.Context(), userID, role)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, scheduling)
}

// @Tags Scheduling
// @Summary Get process scheduling
// @Description Get process scheduling
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /scheduling/processes [get]
func (h *SchedulingHandler) GetProcessScheduling(c *gin.Context) {
	userID := ctx.GetUserID(c)
	role := ctx.GetRole(c)
	scheduling, err := h.service.GetProcessScheduling(c.Request.Context(), userID, role)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, scheduling)
}

// @Tags Scheduling
// @Summary Get task scheduling
// @Description Get task scheduling
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /scheduling/tasks [get]
func (h *SchedulingHandler) GetTaskScheduling(c *gin.Context) {
	userID := ctx.GetUserID(c)
	role := ctx.GetRole(c)
	scheduling, err := h.service.GetTaskScheduling(c.Request.Context(), userID, role)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, scheduling)
}
