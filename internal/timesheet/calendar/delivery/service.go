package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/timesheet/calendar/dto"
	"github.com/Koshsky/erp-backend/pkg/date"
)

type CalendarService interface {
	GetCalendar(
		ctx context.Context,
		userID int64,
		viewScope string,
		start, end date.Date,
	) (*dto.CalendarPlanning, error)
}
