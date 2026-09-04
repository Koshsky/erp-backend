package service

import (
	"context"
	"log/slog"
	"sort"

	repo "github.com/Koshsky/erp-backend/internal/planning/repository"
	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"

	"github.com/Koshsky/erp-backend/internal/planning/dto"
)

type PlanningService struct {
	logger     *slog.Logger
	tracer     *tracingpkg.Tracer
	repository PlanningRepository
}

// NewPlanningService builds the PlanningService service.
func NewPlanningService(logger *slog.Logger, tracer *tracingpkg.Tracer, r *repo.PlanningRepository) *PlanningService {
	return &PlanningService{
		logger:     logger,
		tracer:     tracer,
		repository: r,
	}
}

func (s *PlanningService) GetProjectPlanning(
	ctx context.Context,
	userID int64,
	viewScope string,
) (*dto.ProjectPlanning, error) {
	ctx, end := s.tracer.Start(ctx, "planning.GetProjectPlanning")
	defer end(nil)

	projects, err := s.repository.ListProjects(ctx, userID, viewScope)
	if err != nil {
		return nil, err
	}

	return &dto.ProjectPlanning{
		Projects: projects,
	}, nil
}

func (s *PlanningService) GetProcessPlanning(
	ctx context.Context,
	userID int64,
	viewScope string,
) (*dto.ProcessPlanning, error) {
	ctx, end := s.tracer.Start(ctx, "planning.GetProcessPlanning")
	defer end(nil)

	// Scoped by process.view (ListProcesses): a caller with the right sees its
	// processes even when it has no project.view (e.g. vp). The processes are
	// grouped under their parent projects, which are re-fetched by ids.
	processes, err := s.repository.ListProcesses(ctx, userID, viewScope)
	if err != nil {
		return nil, err
	}
	grouped := make(map[int64][]dto.Process)
	for _, p := range processes {
		grouped[p.ProjectID] = append(grouped[p.ProjectID], p)
	}
	projectIDs := make([]int64, 0, len(grouped))
	for projectID := range grouped {
		projectIDs = append(projectIDs, projectID)
	}
	projects, err := s.repository.ListProjectsByIDs(ctx, projectIDs)
	if err != nil {
		return nil, err
	}
	byID := make(map[int64]dto.Project, len(projects))
	for _, p := range projects {
		byID[p.ID] = p
	}

	planning := dto.ProcessPlanning{Projects: make([]dto.DetailedProject, 0, len(byID))}
	for projectID, ps := range grouped {
		parent, ok := byID[projectID]
		if !ok {
			// The parent project is not visible (deleted/omitted) — the process
			// is not shown either (it cannot be positioned without a project).
			continue
		}
		planning.Projects = append(planning.Projects, dto.DetailedProject{
			Project:   parent,
			Processes: ps,
		})
	}
	// Same ordering as the projects aggregate: by priority, then id.
	sort.Slice(planning.Projects, func(i, j int) bool {
		a, b := planning.Projects[i].Project, planning.Projects[j].Project
		if a.Priority != b.Priority {
			return a.Priority < b.Priority
		}
		return a.ID < b.ID
	})
	return &planning, nil
}

func (s *PlanningService) GetTaskPlanning(
	ctx context.Context,
	userID int64,
	viewScope string,
) (*dto.TaskPlanning, error) {
	ctx, end := s.tracer.Start(ctx, "planning.GetTaskPlanning")
	defer end(nil)

	processes, err := s.loadProcesses(ctx, userID, viewScope)
	if err != nil {
		return nil, err
	}

	if len(processes) == 0 {
		return &dto.TaskPlanning{
			Processes: []dto.DetailedProcess{},
		}, nil
	}

	milestones, tasks, assignments, resourcesMap, commentCounts, err := s.loadAllData(ctx, processes)
	if err != nil {
		return nil, err
	}

	return s.buildPlanning(processes, milestones, tasks, assignments, resourcesMap, commentCounts), nil
}
