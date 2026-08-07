package dto

import (
	"time"

	"github.com/Koshsky/erp-backend/internal/common/date"
)

type CalendarPlanning struct {
	Resources []ResourceCalendar `json:"resources"`
}

type ResourceCalendar struct {
	ResourceID int64                `json:"resource_id" example:"1"`
	Title      string               `json:"title"       example:"Инженер"`
	Code       string               `json:"code"        example:"И"`
	Periods    []AvailabilityPeriod `json:"periods"`
}

// AvailabilityPeriod — диапазон постоянной доступности категории: между границами
// (наймы/увольнения/интервалы состояний) наборы активных и недоступных не меняются.
type AvailabilityPeriod struct {
	StartDate   date.Date `json:"start_date"  example:"2026-07-15" format:"date"`
	EndDate     date.Date `json:"end_date"    example:"2026-07-19" format:"date"`
	Capacity    int       `json:"capacity"    example:"7"`
	Unavailable int       `json:"unavailable" example:"0"`
	Available   int       `json:"available"   example:"7"`
}

// ResourceInfo — справочная информация категории для календаря.
type ResourceInfo struct {
	ID      int64  `json:"id"       example:"1"`
	Title   string `json:"title"    example:"Инженер"`
	Code    string `json:"code"     example:"И"`
	OwnerID *int64 `json:"owner_id" example:"3"`
}

// CalendarEmployee — сотрудник (со своим интервалом работы) для расчёта календаря.
type CalendarEmployee struct {
	EmployeeID      int64
	ResourceID      int64
	HireDate        *time.Time
	TerminationDate *time.Time
}

// UnavailableRange — интервал отсутствия (is_available = false) для расчёта календаря.
type UnavailableRange struct {
	ResourceID int64
	StartDate  time.Time
	EndDate    time.Time
}
