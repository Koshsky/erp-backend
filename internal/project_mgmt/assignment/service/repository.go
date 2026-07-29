package service

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/domain"
)

type AssignmentRepository interface {
	CreateAssignment(ctx context.Context, Assignment domain.Assignment) (*domain.Assignment, error)
	GetAssignment(ctx context.Context, id int64) (*domain.Assignment, error)
	UpdateAssignment(ctx context.Context, new domain.Assignment) (*domain.Assignment, error)
	DeleteAssignment(ctx context.Context, id int64) error
	ListAssignments(ctx context.Context) ([]domain.Assignment, error)
}
