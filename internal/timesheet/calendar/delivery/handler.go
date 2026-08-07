package delivery

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/common/date"
	"github.com/Koshsky/erp-backend/internal/common/response"
)

type CalendarHandler struct {
	logger  *slog.Logger
	service CalendarService
}

func NewCalendarHandler(logger *slog.Logger, service CalendarService) *CalendarHandler {
	return &CalendarHandler{
		logger:  logger,
		service: service,
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
//	@Success		200			{object}	response.Response{data=dto.CalendarPlanning}
//	@Failure		400			{object}	response.Response{data=nil}
//	@Failure		500			{object}	response.Response{data=nil}
//	@Router			/timesheet/calendar [get]
func (h *CalendarHandler) GetCalendar(c *gin.Context) {
	start, err := date.Parse(c.Query("start_date"))
	if err != nil {
		response.BadRequest(c, "invalid start_date")
		return
	}
	end, err := date.Parse(c.Query("end_date"))
	if err != nil {
		response.BadRequest(c, "invalid end_date")
		return
	}

	planning, err := h.service.GetCalendar(c.Request.Context(), start, end)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, planning)
}
