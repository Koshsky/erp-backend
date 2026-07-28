package delivery

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: move response to common DTO?

type response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

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
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /scheduling/projects [get]
func (h *SchedulingHandler) GetProjectScheduling(c *gin.Context) {
	sheduling, err := h.service.GetProjectScheduling(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get scheduling", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: sheduling})
}

// @Tags Scheduling
// @Summary Get process scheduling
// @Description Get process scheduling
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /scheduling/processes [get]
func (h *SchedulingHandler) GetProcessScheduling(c *gin.Context) {
	sheduling, err := h.service.GetProcessScheduling(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get scheduling", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: sheduling})
}

// @Tags Scheduling
// @Summary Get task scheduling
// @Description Get task scheduling
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /scheduling/tasks [get]
func (h *SchedulingHandler) GetTaskScheduling(c *gin.Context) {
	sheduling, err := h.service.GetTaskScheduling(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get scheduling", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: sheduling})
}
