package service

import (
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/timesheet/state/repository"
)

// ProvideStateService builds the StateService service.
func ProvideStateService(logger *slog.Logger, r *repo.StateRepository) *StateService {
	return &StateService{
		logger:     logger,
		repository: r,
		mapper:     NewStateMapper(),
		validator:  &StateValidator{},
	}
}
