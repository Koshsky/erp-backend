package postgres

import (
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/repository/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProvideMilestoneRepository builds the MilestoneRepository repository.
func ProvideMilestoneRepository(logger *slog.Logger, pool *pgxpool.Pool) *MilestoneRepository {
	return &MilestoneRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}
