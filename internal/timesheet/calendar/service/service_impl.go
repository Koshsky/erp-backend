package service

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	repo "github.com/Koshsky/erp-backend/internal/timesheet/calendar/repository"

	"github.com/Koshsky/erp-backend/internal/timesheet/calendar/dto"
	"github.com/Koshsky/erp-backend/pkg/date"
)

// maxCalendarRange is the maximum calendar range width per request (in days).
const (
	maxCalendarRange = 730
	hoursPerDay      = 24
)

type CalendarService struct {
	logger     *slog.Logger
	repository CalendarRepository
}

// NewCalendarService builds the CalendarService service.
func NewCalendarService(logger *slog.Logger, r *repo.CalendarRepository) *CalendarService {
	return &CalendarService{
		logger:     logger,
		repository: r,
	}
}

// GetCalendar returns resource availability as ranges (constant-availability
// segments): capacity, unavailable and available. Complexity O((E+S) log(E+S))
// depends on employees and state intervals, not the number of days.
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

// availabilityEvent is a change in the number of active/unavailable on a day.
type availabilityEvent struct {
	active int
	absent int
}

// buildPeriods computes [start, end] constant-availability segments via an event
// sweep: segment bounds are hire/termination dates and absence interval edges.
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

// countActiveAt counts employees active on day.
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

// countAbsentAt counts employees absent on day.
func countAbsentAt(day time.Time, ranges []dto.UnavailableRange) int {
	count := 0
	for _, r := range ranges {
		if !r.StartDate.After(day) && !r.EndDate.Before(day) {
			count++
		}
	}
	return count
}
