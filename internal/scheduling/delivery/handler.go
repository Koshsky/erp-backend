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

func (h *SchedulingHandler) GetProjectScheduling(c *gin.Context) {
	sheduling, err := h.service.GetProjectScheduling(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get scheduling", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: sheduling})
}
func (h *SchedulingHandler) GetProcessScheduling(c *gin.Context) {
	sheduling, err := h.service.GetProcessScheduling(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get scheduling", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: sheduling})
}

func (h *SchedulingHandler) GetTaskScheduling(c *gin.Context) {
	sheduling, err := h.service.GetTaskScheduling(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get scheduling", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: sheduling})
}
