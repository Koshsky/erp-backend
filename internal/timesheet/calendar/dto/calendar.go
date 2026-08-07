package dto

import "github.com/Koshsky/erp-backend/internal/common/date"

type CalendarPlanning struct {
	Resources []ResourceCalendar `json:"resources"`
}

type ResourceCalendar struct {
	ResourceID int64             `json:"resource_id" example:"1"`
	Title      string            `json:"title"       example:"Инженер"`
	Code       string            `json:"code"        example:"И"`
	Days       []DayAvailability `json:"days"`
}

type DayAvailability struct {
	Date        date.Date `json:"date"        example:"2026-07-15" format:"date"`
	Capacity    int       `json:"capacity"    example:"7"`
	Unavailable int       `json:"unavailable" example:"2"`
	Available   int       `json:"available"   example:"5"`
}

// ResourceInfo — справочная информация категории для календаря.
type ResourceInfo struct {
	ID    int64  `json:"id"    example:"1"`
	Title string `json:"title" example:"Инженер"`
	Code  string `json:"code"  example:"И"`
}

// ResourceDayCount — промежуточный результат агрегации календаря (ресурс x день).
type ResourceDayCount struct {
	ResourceID int64
	Date       date.Date
	Count      int
}
