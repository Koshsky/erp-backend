package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Koshsky/erp-backend/internal/common/date"
	"github.com/Koshsky/erp-backend/internal/timesheet/calendar/dto"
)

// maxCalendarRange — максимальная ширина диапазона календаря за один запрос (в днях).
const (
	maxCalendarRange = 730
	hoursPerDay      = 24
)

type CalendarService struct {
	logger     *slog.Logger
	repository CalendarRepository
}

func NewCalendarService(logger *slog.Logger, repository CalendarRepository) *CalendarService {
	return &CalendarService{
		logger:     logger,
		repository: repository,
	}
}

// GetCalendar возвращает по-дневную доступность ресурсов (мощность, недоступные, доступные)
// в заданном диапазоне — данные для бесконечного календаря на фронтенде.
func (s *CalendarService) GetCalendar(
	ctx context.Context,
	start, end date.Date,
) (*dto.CalendarPlanning, error) {
	startT, endT := start.Time(), end.Time()
	if endT.Before(startT) {
		return nil, fmt.Errorf("end_date must be greater than or equal to start_date")
	}
	if int(endT.Sub(startT).Hours()/hoursPerDay) > maxCalendarRange-1 {
		return nil, fmt.Errorf("date range must not exceed %d days", maxCalendarRange)
	}

	resources, err := s.repository.ListResources(ctx)
	if err != nil {
		return nil, err
	}

	capacityRows, err := s.repository.ListCapacity(ctx, startT, endT)
	if err != nil {
		return nil, err
	}
	unavailableRows, err := s.repository.ListUnavailable(ctx, startT, endT)
	if err != nil {
		return nil, err
	}

	capacityMap := buildDayCountMap(capacityRows)
	unavailableMap := buildDayCountMap(unavailableRows)

	planning := &dto.CalendarPlanning{
		Resources: make([]dto.ResourceCalendar, 0, len(resources)),
	}
	for _, resource := range resources {
		days := make([]dto.DayAvailability, 0, daysInRange(startT, endT))
		for d := startT; !d.After(endT); d = d.AddDate(0, 0, 1) {
			day := date.From(d)
			capacity := capacityMap[resource.ID][day]
			unavailable := unavailableMap[resource.ID][day]
			days = append(days, dto.DayAvailability{
				Date:        day,
				Capacity:    capacity,
				Unavailable: unavailable,
				Available:   capacity - unavailable,
			})
		}
		planning.Resources = append(planning.Resources, dto.ResourceCalendar{
			ResourceID: resource.ID,
			Title:      resource.Title,
			Code:       resource.Code,
			Days:       days,
		})
	}

	return planning, nil
}

// buildDayCountMap собирает агрегаты "ресурс -> дата -> количество" в карту.
func buildDayCountMap(rows []dto.ResourceDayCount) map[int64]map[date.Date]int {
	result := make(map[int64]map[date.Date]int, len(rows))
	for _, row := range rows {
		if result[row.ResourceID] == nil {
			result[row.ResourceID] = make(map[date.Date]int)
		}
		result[row.ResourceID][row.Date] = row.Count
	}
	return result
}

// daysInRange возвращает число дней в закрытом диапазоне [start, end].
func daysInRange(start, end time.Time) int {
	return int(end.Sub(start).Hours()/hoursPerDay) + 1
}
