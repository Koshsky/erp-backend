package service

import (
	"context"
	"time"

	"github.com/Koshsky/erp-backend/internal/timesheet/calendar/dto"
)

type CalendarRepository interface {
	ListResources(ctx context.Context) ([]dto.ResourceInfo, error)
	ListCapacity(ctx context.Context, start, end time.Time) ([]dto.ResourceDayCount, error)
	ListUnavailable(ctx context.Context, start, end time.Time) ([]dto.ResourceDayCount, error)
}
