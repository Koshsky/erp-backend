package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/timesheet/calendar/dto"
	"github.com/Koshsky/erp-backend/pkg/date"
)

type CalendarService interface {
	GetCalendar(ctx context.Context, start, end date.Date) (*dto.CalendarPlanning, error)
}
