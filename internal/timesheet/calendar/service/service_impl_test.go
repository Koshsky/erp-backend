//nolint:testpackage // tests the unexported buildPeriods directly
package service

import (
	"testing"
	"time"

	"github.com/Koshsky/erp-backend/internal/timesheet/calendar/dto"
	"github.com/Koshsky/erp-backend/pkg/date"
)

func dt(s string) time.Time {
	d, err := date.Parse(s)
	if err != nil {
		panic(err)
	}
	return d.Time()
}

func timePtr(s string) *time.Time {
	t := dt(s)
	return &t
}

func period(start, end string, capacity, unavailable, available int) dto.AvailabilityPeriod {
	return dto.AvailabilityPeriod{
		StartDate:   date.From(dt(start)),
		EndDate:     date.From(dt(end)),
		Capacity:    capacity,
		Unavailable: unavailable,
		Available:   available,
	}
}

func TestBuildPeriodsEmpty(t *testing.T) {
	t.Parallel()
	got := buildPeriods(dt("2026-01-01"), dt("2026-01-31"), nil, nil)
	want := []dto.AvailabilityPeriod{period("2026-01-01", "2026-01-31", 0, 0, 0)}
	if len(got) != 1 || got[0] != want[0] {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestBuildPeriodsHireAndVacation(t *testing.T) {
	t.Parallel()
	employees := []dto.CalendarEmployee{
		{EmployeeID: 1, ResourceID: 1, HireDate: timePtr("2025-01-01")},
		{EmployeeID: 2, ResourceID: 1, HireDate: timePtr("2026-06-01")},
	}
	ranges := []dto.UnavailableRange{
		{ResourceID: 1, StartDate: dt("2026-07-20"), EndDate: dt("2026-08-02")},
	}
	got := buildPeriods(dt("2026-01-01"), dt("2026-08-31"), employees, ranges)
	want := []dto.AvailabilityPeriod{
		period("2026-01-01", "2026-05-31", 1, 0, 1),
		period("2026-06-01", "2026-07-19", 2, 0, 2),
		period("2026-07-20", "2026-08-02", 2, 1, 1),
		period("2026-08-03", "2026-08-31", 2, 0, 2),
	}
	if len(got) != len(want) {
		t.Fatalf("got %d periods, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("period %d: got %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestBuildPeriodsTerminationAndMerge(t *testing.T) {
	t.Parallel()
	employees := []dto.CalendarEmployee{
		{EmployeeID: 1, ResourceID: 1, HireDate: timePtr("2024-01-01"), TerminationDate: timePtr("2026-03-15")},
		{EmployeeID: 2, ResourceID: 1, HireDate: timePtr("2024-01-01")},
	}
	got := buildPeriods(dt("2026-01-01"), dt("2026-12-31"), employees, nil)
	want := []dto.AvailabilityPeriod{
		period("2026-01-01", "2026-03-15", 2, 0, 2),
		period("2026-03-16", "2026-12-31", 1, 0, 1),
	}
	if len(got) != len(want) {
		t.Fatalf("got %d periods, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("period %d: got %+v, want %+v", i, got[i], want[i])
		}
	}
}
