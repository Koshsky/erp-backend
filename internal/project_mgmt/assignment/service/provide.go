package service

import (
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/repository"
)

// ProvideAssignmentService builds the AssignmentService service.
func ProvideAssignmentService(logger *slog.Logger, r *repo.AssignmentRepository) *AssignmentService {
	return &AssignmentService{
		logger:     logger,
		repository: r,
		mapper:     NewAssignmentMapper(),
		validator:  &AssignmentValidator{},
	}
}
