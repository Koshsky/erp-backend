package postgres

import (
	"log/slog"

	"github.com/Koshsky/erp/api/internal/repository/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	logger  *slog.Logger
	pool    *pgxpool.Pool
	queries *sqlc.Queries

	*ProjectRepository
	*ProcessRepository
	*MilestoneRepository
	*TaskRepository
	*ResourceRepository
	*AssignmentRepository
	*UserRepository
}

func NewRepository(logger *slog.Logger, pool *pgxpool.Pool) *Repository {
	queries := sqlc.New(pool)
	return &Repository{
		logger:               logger,
		pool:                 pool,
		queries:              queries,
		ProjectRepository:    NewProjectRepository(logger, queries),
		ProcessRepository:    NewProcessRepository(logger, queries),
		MilestoneRepository:  NewMilestoneRepository(logger, queries),
		TaskRepository:       NewTaskRepository(logger, queries),
		ResourceRepository:   NewResourceRepository(logger, queries),
		AssignmentRepository: NewAssignmentRepository(logger, queries),
		UserRepository:       NewUserRepository(logger, queries),
	}
}

func (r *Repository) Close() {
	r.pool.Close()
}
