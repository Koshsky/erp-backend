package delivery

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/common/helpers"
	"github.com/Koshsky/erp-backend/internal/common/response"
)

type PlanningHandler struct {
	logger  *slog.Logger
	service MilestoneService
}

func NewPlanningHandler(logger *slog.Logger, service MilestoneService) *PlanningHandler {
	return &PlanningHandler{
		logger:  logger,
		service: service,
	}
}

// GetProjectPlanning handles the request to get project planning (project portfolio).
//
//	@Tags			Planning
//	@Summary		Get project planning
//	@Description	Get project planning (project portfolio)
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	response.Response{data=dto.ProjectPlanning}
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/planning/projects [get]
func (h *PlanningHandler) GetProjectPlanning(c *gin.Context) {
	user, err := helpers.GetUser(c)
	if err != nil {
		response.Unauthorized(c, "authentication required")
		return
	}

	planning, err := h.service.GetProjectPlanning(c.Request.Context(), user.ID, user.Role)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, planning)
}

// GetProcessPlanning handles the request to get process planning.
//
//	@Tags			Planning
//	@Summary		Get process planning
//	@Description	Get process planning
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	response.Response{data=dto.ProcessPlanning}
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/planning/processes [get]
func (h *PlanningHandler) GetProcessPlanning(c *gin.Context) {
	user, err := helpers.GetUser(c)
	if err != nil {
		response.Unauthorized(c, "authentication required")
		return
	}

	planning, err := h.service.GetProcessPlanning(c.Request.Context(), user.ID, user.Role)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, planning)
}

// GetTaskPlanning handles the request to get task planning.
//
//	@Tags			Planning
//	@Summary		Get task planning
//	@Description	Get task planning
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	response.Response{data=dto.TaskPlanning}
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/planning/tasks [get]
func (h *PlanningHandler) GetTaskPlanning(c *gin.Context) {
	user, err := helpers.GetUser(c)
	if err != nil {
		response.Unauthorized(c, "authentication required")
		return
	}

	planning, err := h.service.GetTaskPlanning(c.Request.Context(), user.ID, user.Role)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, planning)
}
