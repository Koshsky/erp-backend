//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/timesheet/calendar/repository/sqlc"
)

type CalendarRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

// NewCalendarRepository builds the CalendarRepository repository.
func NewCalendarRepository(logger *slog.Logger, pool *pgxpool.Pool) *CalendarRepository {
	return &CalendarRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *CalendarRepository) ListResources(ctx context.Context) ([]sqlc.ListResourcesRow, error) {
	return r.db.ListResources(ctx)
}

// ListEmployeesForCalendar returns resource members active within the window for the calendar.
func (r *CalendarRepository) ListEmployeesForCalendar(
	ctx context.Context,
	start, end time.Time,
) ([]sqlc.ListEmployeesForCalendarRow, error) {
	return r.db.ListEmployeesForCalendar(ctx, sqlc.ListEmployeesForCalendarParams{
		StartDate: start,
		EndDate:   end,
	})
}

// ListUnavailableRanges returns absence intervals overlapping the window.
func (r *CalendarRepository) ListUnavailableRanges(
	ctx context.Context,
	start, end time.Time,
) ([]sqlc.ListUnavailableRangesRow, error) {
	return r.db.ListUnavailableRanges(ctx, sqlc.ListUnavailableRangesParams{
		StartDate: start,
		EndDate:   end,
	})
}
