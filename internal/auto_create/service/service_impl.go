package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	repo "github.com/Koshsky/erp-backend/internal/auto_create/repository"
	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"

	"github.com/Koshsky/erp-backend/internal/auto_create/dto"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

type AutoCreateService struct {
	logger     *slog.Logger
	tracer     *tracingpkg.Tracer
	repository AutoCreateRepository
}

// NewAutoCreateService builds the AutoCreateService service.
func NewAutoCreateService(
	logger *slog.Logger,
	tracer *tracingpkg.Tracer,
	r *repo.AutoCreateRepository,
) *AutoCreateService {
	return &AutoCreateService{
		logger:     logger,
		tracer:     tracer,
		repository: r,
	}
}

// GetConfig returns the current auto-create config (disabled+empty if unset).
func (s *AutoCreateService) GetConfig(ctx context.Context) (*dto.AutoCreateConfig, error) {
	ctx, end := s.tracer.Start(ctx, "autocreate.GetConfig")
	defer end(nil)

	return s.repository.GetConfig(ctx)
}

// SaveConfig replaces the auto-create config; validates shape and that all
// referenced owners/resources exist (otherwise the auto-created project would fail on an FK).
func (s *AutoCreateService) SaveConfig(ctx context.Context, cfg *dto.AutoCreateConfig) error {
	ctx, end := s.tracer.Start(ctx, "autocreate.SaveConfig")
	defer end(nil)

	if err := validateConfig(cfg); err != nil {
		return err
	}

	if resourceIDs := collectResourceIDs(cfg); len(resourceIDs) > 0 {
		existing, err := s.repository.ExistingResources(ctx, resourceIDs)
		if err != nil {
			return err
		}
		for _, id := range resourceIDs {
			if _, ok := existing[id]; !ok {
				return errors.NewValidationError(fmt.Sprintf("ресурс %d не найден", id))
			}
		}
	}

	if ownerIDs := collectOwnerIDs(cfg); len(ownerIDs) > 0 {
		existing, err := s.repository.ExistingUsers(ctx, ownerIDs)
		if err != nil {
			return err
		}
		for _, id := range ownerIDs {
			if _, ok := existing[id]; !ok {
				return errors.NewValidationError(fmt.Sprintf("владелец процесса %d не найден", id))
			}
		}
	}

	return s.repository.UpsertConfig(ctx, cfg)
}

// Limits protecting project creation from a pathological template: the DB
// trigger (V8) applies the whole template atomically inside the project INSERT
// transaction, so an unbounded template is a self-inflicted DoS on every new
// project. The frontend mirrors these limits in its own validation.
const (
	maxProcesses             = 20
	maxTasksPerProcess       = 50
	maxResourcesPerTask      = 10
	maxAssignmentsTotal      = 500
	maxQuantityPerAssignment = 99
)

func validateConfig(cfg *dto.AutoCreateConfig) error {
	if len(cfg.Processes) > maxProcesses {
		return errors.NewValidationError(fmt.Sprintf(
			"слишком много процессов: максимум %d", maxProcesses,
		))
	}
	totalAssignments := 0
	for pi, p := range cfg.Processes {
		if len(p.Tasks) > maxTasksPerProcess {
			return errors.NewValidationError(fmt.Sprintf(
				"процесс %d: слишком много задач: максимум %d", pi+1, maxTasksPerProcess,
			))
		}
		if err := validateProcess(p, pi); err != nil {
			return err
		}
		for _, t := range p.Tasks {
			totalAssignments += len(t.Resources)
		}
	}
	if totalAssignments > maxAssignmentsTotal {
		return errors.NewValidationError(fmt.Sprintf(
			"слишком много назначений ресурсов: максимум %d", maxAssignmentsTotal,
		))
	}
	return nil
}

func validateProcess(p dto.ProcessTemplate, pi int) error {
	if strings.TrimSpace(p.Title) == "" {
		return errors.NewValidationError(fmt.Sprintf("процесс %d: название не заполнено", pi+1))
	}
	for ti, t := range p.Tasks {
		if err := validateTask(p.Title, t, ti); err != nil {
			return err
		}
	}
	return nil
}

func validateTask(processTitle string, t dto.TaskTemplate, ti int) error {
	if strings.TrimSpace(t.Title) == "" {
		return errors.NewValidationError(
			fmt.Sprintf("процесс «%s», задача %d: название не заполнено", processTitle, ti+1),
		)
	}
	if len(t.Resources) > maxResourcesPerTask {
		return errors.NewValidationError(
			fmt.Sprintf("процесс «%s», задача %d: слишком много ресурсов: максимум %d",
				processTitle, ti+1, maxResourcesPerTask),
		)
	}
	seen := make(map[int64]struct{}, len(t.Resources))
	for ri, res := range t.Resources {
		if err := validateResource(t.Title, res, ri, seen); err != nil {
			return err
		}
	}
	return nil
}

func validateResource(taskTitle string, res dto.ResourceBinding, ri int, seen map[int64]struct{}) error {
	if res.ResourceID <= 0 {
		return errors.NewValidationError(
			fmt.Sprintf("задача «%s», ресурс %d: некорректный resource_id", taskTitle, ri+1),
		)
	}
	if res.Quantity <= 0 {
		return errors.NewValidationError(
			fmt.Sprintf("задача «%s», ресурс %d: количество должно быть больше 0", taskTitle, ri+1),
		)
	}
	if res.Quantity > maxQuantityPerAssignment {
		return errors.NewValidationError(
			fmt.Sprintf(
				"задача «%s», ресурс %d: количество не должно превышать %d",
				taskTitle, ri+1, maxQuantityPerAssignment,
			),
		)
	}
	if _, dup := seen[res.ResourceID]; dup {
		return errors.NewValidationError(
			fmt.Sprintf("задача «%s»: ресурс %d указан дважды", taskTitle, res.ResourceID),
		)
	}
	seen[res.ResourceID] = struct{}{}
	return nil
}

func collectResourceIDs(cfg *dto.AutoCreateConfig) []int64 {
	seen := make(map[int64]struct{})
	var out []int64
	for _, p := range cfg.Processes {
		for _, t := range p.Tasks {
			for _, res := range t.Resources {
				if _, ok := seen[res.ResourceID]; !ok {
					seen[res.ResourceID] = struct{}{}
					out = append(out, res.ResourceID)
				}
			}
		}
	}
	return out
}

func collectOwnerIDs(cfg *dto.AutoCreateConfig) []int64 {
	seen := make(map[int64]struct{})
	var out []int64
	for _, p := range cfg.Processes {
		if p.OwnerID != nil {
			if _, ok := seen[*p.OwnerID]; !ok {
				seen[*p.OwnerID] = struct{}{}
				out = append(out, *p.OwnerID)
			}
		}
	}
	return out
}
