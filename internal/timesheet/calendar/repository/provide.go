package repository

import (
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/timesheet/calendar/repository/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProvideCalendarRepository builds the CalendarRepository repository.
func ProvideCalendarRepository(logger *slog.Logger, pool *pgxpool.Pool) *CalendarRepository {
	return &CalendarRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}
