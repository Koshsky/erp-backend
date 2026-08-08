package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/common/date"
	"github.com/Koshsky/erp-backend/internal/timesheet/calendar/dto"
)

type CalendarService interface {
	GetCalendar(ctx context.Context, start, end date.Date) (*dto.CalendarPlanning, error)
}
