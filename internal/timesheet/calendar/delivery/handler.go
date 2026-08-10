package delivery

import (
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/timesheet/calendar/service"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/response"
	"github.com/Koshsky/erp-backend/pkg/date"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

type CalendarHandler struct {
	logger  *slog.Logger
	service CalendarService
	mw      *rbac.Middleware
}

// NewCalendarHandler builds the CalendarHandler handler.
func NewCalendarHandler(logger *slog.Logger, svc *service.CalendarService, mw *rbac.Middleware) *CalendarHandler {
	return &CalendarHandler{
		logger:  logger,
		service: svc,
		mw:      mw,
	}
}

// GetCalendar handles the request to get the resource availability calendar.
//
//	@Tags			TimesheetCalendar
//	@Summary		Get resource calendar
//	@Description	Get per-day resource availability (capacity/unavailable/available) for a date range
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			start_date	query		string	true	"Start date (YYYY-MM-DD)"
//	@Param			end_date	query		string	true	"End date (YYYY-MM-DD)"
//	@Success		200			{object}	response.SuccessResponse{data=dto.CalendarPlanning,error=nil}
//	@Failure		400			{object}	response.ErrorResponse{data=nil}
//	@Failure		500			{object}	response.ErrorResponse{data=nil}
//	@Router			/timesheet/calendar [get]
func (h *CalendarHandler) GetCalendar(c *gin.Context) {
	start, err := date.Parse(c.Query("start_date"))
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid start_date")
		return
	}
	end, err := date.Parse(c.Query("end_date"))
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid end_date")
		return
	}

	planning, err := h.service.GetCalendar(c.Request.Context(), start, end)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, planning)
}
