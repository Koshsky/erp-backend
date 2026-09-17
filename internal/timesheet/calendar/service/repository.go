package service

import (
	"context"
	"time"

	"github.com/Koshsky/erp-backend/internal/timesheet/calendar/repository/sqlc"
)

type CalendarRepository interface {
	ListResources(ctx context.Context) ([]sqlc.ListResourcesRow, error)
	ListEmployeesForCalendar(ctx context.Context, start, end time.Time) ([]sqlc.ListEmployeesForCalendarRow, error)
	ListUnavailableRanges(ctx context.Context, start, end time.Time) ([]sqlc.ListUnavailableRangesRow, error)
}
