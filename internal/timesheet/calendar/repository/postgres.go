//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/timesheet/calendar/dto"
	"github.com/Koshsky/erp-backend/internal/timesheet/calendar/repository/sqlc"
)

type CalendarRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

func NewCalendarRepository(logger *slog.Logger, pool *pgxpool.Pool) *CalendarRepository {
	return &CalendarRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *CalendarRepository) ListResources(ctx context.Context) ([]dto.ResourceInfo, error) {
	rows, err := r.db.ListResources(ctx)
	if err != nil {
		return nil, err
	}
	resources := make([]dto.ResourceInfo, 0, len(rows))
	for _, row := range rows {
		resources = append(resources, dto.ResourceInfo{
			ID:    row.ID,
			Title: row.Title,
			Code:  row.Code,
		})
	}
	return resources, nil
}

// ListEmployeesForCalendar возвращает сотрудников, активных в окне, для расчёта календаря.
func (r *CalendarRepository) ListEmployeesForCalendar(
	ctx context.Context,
	start, end time.Time,
) ([]dto.CalendarEmployee, error) {
	rows, err := r.db.ListEmployeesForCalendar(ctx, sqlc.ListEmployeesForCalendarParams{
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		return nil, err
	}
	employees := make([]dto.CalendarEmployee, 0, len(rows))
	for _, row := range rows {
		employees = append(employees, dto.CalendarEmployee{
			EmployeeID:      row.ID,
			ResourceID:      row.ResourceID,
			HireDate:        fromDate(row.HireDate),
			TerminationDate: fromDate(row.TerminationDate),
		})
	}
	return employees, nil
}

// ListUnavailableRanges возвращает интервалы отсутствий, пересекающие окно.
func (r *CalendarRepository) ListUnavailableRanges(
	ctx context.Context,
	start, end time.Time,
) ([]dto.UnavailableRange, error) {
	rows, err := r.db.ListUnavailableRanges(ctx, sqlc.ListUnavailableRangesParams{
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		return nil, err
	}
	ranges := make([]dto.UnavailableRange, 0, len(rows))
	for _, row := range rows {
		ranges = append(ranges, dto.UnavailableRange{
			ResourceID: row.ResourceID,
			StartDate:  row.StartDate,
			EndDate:    row.EndDate,
		})
	}
	return ranges, nil
}

// fromDate разворачивает nullable-дату (pgtype.Date) в [time.Time].
func fromDate(v pgtype.Date) *time.Time {
	if !v.Valid {
		return nil
	}
	t := v.Time
	return &t
}
