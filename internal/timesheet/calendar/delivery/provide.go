package delivery

import (
	"log/slog"

	rbac "github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/timesheet/calendar/service"
)

// ProvideCalendarHandler builds the CalendarHandler handler.
func ProvideCalendarHandler(logger *slog.Logger, svc *service.CalendarService, mw *rbac.Middleware) *CalendarHandler {
	return &CalendarHandler{
		logger:  logger,
		service: svc,
		mw:      mw,
	}
}
