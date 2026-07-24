package delivery

import (
	"context"

	"github.com/Koshsky/erp/api/internal/dto"
)

type AssignmentService interface {
	ListAssignments(ctx context.Context) ([]dto.AssignmentResponse, error)
	GetAssignment(ctx context.Context, id int64) (*dto.AssignmentResponse, error)
	CreateAssignment(ctx context.Context, assignment dto.CreateAssignmentRequest) (*dto.AssignmentResponse, error)
	DeleteAssignment(ctx context.Context, id int64) error
	UpdateAssignment(ctx context.Context, id int64, assignment dto.UpdateAssignmentRequest) (*dto.AssignmentResponse, error)
}
