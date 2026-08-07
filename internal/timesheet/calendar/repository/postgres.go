//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/common/date"
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

// ListCapacity возвращает по-дневную мощность каждой категории ресурсов в диапазоне.
func (r *CalendarRepository) ListCapacity(ctx context.Context, start, end time.Time) ([]dto.ResourceDayCount, error) {
	rows, err := r.db.ListCalendarCapacity(ctx, sqlc.ListCalendarCapacityParams{
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		return nil, err
	}
	result := make([]dto.ResourceDayCount, 0, len(rows))
	for _, row := range rows {
		result = append(result, dto.ResourceDayCount{
			ResourceID: row.ResourceID,
			Date:       date.From(row.GDate),
			Count:      int(row.Capacity),
		})
	}
	return result, nil
}

// ListUnavailable возвращает по-дневное количество недоступных сотрудников по категориям.
func (r *CalendarRepository) ListUnavailable(
	ctx context.Context,
	start, end time.Time,
) ([]dto.ResourceDayCount, error) {
	rows, err := r.db.ListCalendarUnavailable(ctx, sqlc.ListCalendarUnavailableParams{
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		return nil, err
	}
	result := make([]dto.ResourceDayCount, 0, len(rows))
	for _, row := range rows {
		result = append(result, dto.ResourceDayCount{
			ResourceID: row.ResourceID,
			Date:       date.From(row.GDate),
			Count:      int(row.Unavailable),
		})
	}
	return result, nil
}
