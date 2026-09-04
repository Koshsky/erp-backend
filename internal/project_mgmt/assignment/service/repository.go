package service

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/domain"
)

type AssignmentRepository interface {
	CreateAssignment(ctx context.Context, Assignment domain.Assignment) (*domain.Assignment, error)
	FindAssignment(ctx context.Context, id int64) (*domain.Assignment, error)
	UpdateAssignment(ctx context.Context, assignment domain.Assignment) (*domain.Assignment, error)
	DeleteAssignment(ctx context.Context, id int64) error
	ListAssignments(
		ctx context.Context,
		userID int64,
		viewScope string,
		ownerID int64,
		limit, offset int,
	) ([]domain.Assignment, error)
	CountAssignments(ctx context.Context, userID int64, viewScope string, ownerID int64) (int64, error)
}
