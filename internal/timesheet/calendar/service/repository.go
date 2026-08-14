package service

import (
	"context"
	"time"

	"github.com/Koshsky/erp-backend/internal/timesheet/calendar/dto"
)

type CalendarRepository interface {
	ListResources(ctx context.Context) ([]dto.ResourceInfo, error)
	ListEmployeesForCalendar(ctx context.Context, start, end time.Time) ([]dto.CalendarMember, error)
	ListUnavailableRanges(ctx context.Context, start, end time.Time) ([]dto.UnavailableRange, error)
}
