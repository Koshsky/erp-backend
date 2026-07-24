package domain

import (
	"context"
)

type RepositoryInterface interface {
	CreateAssignment(ctx context.Context, Assignment Assignment) (*Assignment, error)
	GetAssignment(ctx context.Context, id int64) (*Assignment, error)
	UpdateAssignment(ctx context.Context, new Assignment) (*Assignment, error)
	DeleteAssignment(ctx context.Context, id int64) error
	ListAssignments(ctx context.Context) ([]Assignment, error)
}
