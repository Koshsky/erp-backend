package service

import (
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/timesheet/calendar/repository"
)

// ProvideCalendarService builds the CalendarService service.
func ProvideCalendarService(logger *slog.Logger, r *repo.CalendarRepository) *CalendarService {
	return &CalendarService{
		logger:     logger,
		repository: r,
	}
}
