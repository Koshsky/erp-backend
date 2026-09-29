package service

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/timesheet/calendar/repository/sqlc"
	"github.com/Koshsky/erp-backend/pkg/date"
)

type CalendarRepository interface {
	ListResources(ctx context.Context, userID int64, viewScope string) ([]sqlc.ListResourcesRow, error)
	ListEmployeesForCalendar(ctx context.Context, start, end date.Date) ([]sqlc.ListEmployeesForCalendarRow, error)
	ListUnavailableRanges(ctx context.Context, start, end date.Date) ([]sqlc.ListUnavailableRangesRow, error)
}
