package service

import (
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/project_mgmt/task/repository"
)

// ProvideTaskService builds the TaskService service.
func ProvideTaskService(logger *slog.Logger, r *repo.TaskRepository) *TaskService {
	return &TaskService{
		logger:     logger,
		repository: r,
		mapper:     &TaskMapper{},
		validator:  &TaskValidator{},
	}
}
