package service

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/Koshsky/erp-backend/internal/common/date"
	"github.com/Koshsky/erp-backend/internal/timesheet/calendar/dto"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
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

// GetCalendar возвращает доступность ресурсов диапазонами (сегменты постоянной доступности):
// мощность (штат), недоступные и доступные. Сложность O((E+S) log(E+S)) — зависит от числа
// сотрудников и интервалов состояний, а не от числа дней в диапазоне.
func (s *CalendarService) GetCalendar(
	ctx context.Context,
	start, end date.Date,
	userID int64,
	role string,
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
	// vp видит в планировщике только свои ресурсы.
	if role == userdomain.ProcessOwner {
		resources = ownedResources(resources, userID)
	}
	employees, err := s.repository.ListEmployeesForCalendar(ctx, startT, endT)
	if err != nil {
		return nil, err
	}
	ranges, err := s.repository.ListUnavailableRanges(ctx, startT, endT)
	if err != nil {
		return nil, err
	}

	employeesByResource := groupEmployees(employees)
	rangesByResource := groupRanges(ranges)

	planning := &dto.CalendarPlanning{
		Resources: make([]dto.ResourceCalendar, 0, len(resources)),
	}
	for _, resource := range resources {
		periods := buildPeriods(
			startT,
			endT,
			employeesByResource[resource.ID],
			rangesByResource[resource.ID],
		)
		planning.Resources = append(planning.Resources, dto.ResourceCalendar{
			ResourceID: resource.ID,
			Title:      resource.Title,
			Code:       resource.Code,
			Periods:    periods,
		})
	}

	return planning, nil
}

// ownedResources оставляет только ресурсы, принадлежащие пользователю.
func ownedResources(resources []dto.ResourceInfo, userID int64) []dto.ResourceInfo {
	result := resources[:0]
	for _, resource := range resources {
		if resource.OwnerID != nil && *resource.OwnerID == userID {
			result = append(result, resource)
		}
	}
	return result
}

func groupEmployees(employees []dto.CalendarEmployee) map[int64][]dto.CalendarEmployee {
	result := make(map[int64][]dto.CalendarEmployee, len(employees))
	for _, employee := range employees {
		result[employee.ResourceID] = append(result[employee.ResourceID], employee)
	}
	return result
}

func groupRanges(ranges []dto.UnavailableRange) map[int64][]dto.UnavailableRange {
	result := make(map[int64][]dto.UnavailableRange, len(ranges))
	for _, r := range ranges {
		result[r.ResourceID] = append(result[r.ResourceID], r)
	}
	return result
}

// availabilityEvent — изменение числа активных и недоступных в конкретный день.
type availabilityEvent struct {
	active int
	absent int
}

// buildPeriods вычисляет сегменты [start, end] постоянной доступности по событийному sweep:
// границы сегментов — даты наймов/увольнений и начала/концы интервалов отсутствий.
func buildPeriods(
	start, end time.Time,
	employees []dto.CalendarEmployee,
	ranges []dto.UnavailableRange,
) []dto.AvailabilityPeriod {
	active := countActiveAt(start, employees)
	absent := countAbsentAt(start, ranges)

	events := make(map[time.Time]availabilityEvent)
	addEvent := func(day time.Time, deltaActive, deltaAbsent int) {
		if day.After(start) && !day.After(end) {
			event := events[day]
			event.active += deltaActive
			event.absent += deltaAbsent
			events[day] = event
		}
	}
	for _, employee := range employees {
		if employee.HireDate != nil {
			addEvent(*employee.HireDate, 1, 0)
		}
		if employee.TerminationDate != nil {
			addEvent(employee.TerminationDate.AddDate(0, 0, 1), -1, 0)
		}
	}
	for _, r := range ranges {
		addEvent(r.StartDate, 0, 1)
		addEvent(r.EndDate.AddDate(0, 0, 1), 0, -1)
	}

	days := make([]time.Time, 0, len(events))
	for day := range events {
		days = append(days, day)
	}
	slices.SortFunc(days, func(a, b time.Time) int { return a.Compare(b) })

	periods := make([]dto.AvailabilityPeriod, 0, len(days)+1)
	prev := start
	appendPeriod := func(segStart, segEnd time.Time) {
		if segStart.After(segEnd) {
			return
		}
		if n := len(periods); n > 0 &&
			periods[n-1].Capacity == active && periods[n-1].Unavailable == absent {
			periods[n-1].EndDate = date.From(segEnd)
			return
		}
		periods = append(periods, dto.AvailabilityPeriod{
			StartDate:   date.From(segStart),
			EndDate:     date.From(segEnd),
			Capacity:    active,
			Unavailable: absent,
			Available:   active - absent,
		})
	}

	for _, day := range days {
		appendPeriod(prev, day.AddDate(0, 0, -1))
		event := events[day]
		active += event.active
		absent += event.absent
		prev = day
	}
	appendPeriod(prev, end)

	return periods
}

// countActiveAt считает сотрудников, активных в день day.
func countActiveAt(day time.Time, employees []dto.CalendarEmployee) int {
	count := 0
	for _, employee := range employees {
		if (employee.HireDate == nil || !employee.HireDate.After(day)) &&
			(employee.TerminationDate == nil || !employee.TerminationDate.Before(day)) {
			count++
		}
	}
	return count
}

// countAbsentAt считает отсутствующих сотрудников в день day.
func countAbsentAt(day time.Time, ranges []dto.UnavailableRange) int {
	count := 0
	for _, r := range ranges {
		if !r.StartDate.After(day) && !r.EndDate.Before(day) {
			count++
		}
	}
	return count
}
