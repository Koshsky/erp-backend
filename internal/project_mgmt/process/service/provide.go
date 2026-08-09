package service

import (
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/project_mgmt/process/repository"
)

// ProvideProcessService builds the ProcessService service.
func ProvideProcessService(logger *slog.Logger, r *repo.ProcessRepository) *ProcessService {
	return &ProcessService{
		logger:     logger,
		repository: r,
		mapper:     NewProcessMapper(),
		validator:  &ProcessValidator{},
	}
}
