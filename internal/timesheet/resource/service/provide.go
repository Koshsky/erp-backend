package service

import (
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/timesheet/resource/repository"
)

// ProvideResourceService builds the ResourceService service.
func ProvideResourceService(logger *slog.Logger, r *repo.ResourceRepository) *ResourceService {
	return &ResourceService{
		logger:     logger,
		repository: r,
		mapper:     NewResourceMapper(),
		validator:  &ResourceValidator{},
	}
}
