package service

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/planning/dto"
)

type PlanningService struct {
	logger     *slog.Logger
	repository PlanningRepository
}

func NewPlanningService(logger *slog.Logger, repository PlanningRepository) *PlanningService {
	return &PlanningService{
		logger:     logger,
		repository: repository,
	}
}

func (s *PlanningService) GetProjectPlanning(
	ctx context.Context,
	userID int64,
	role string,
) (*dto.ProjectPlanning, error) {
	projects, err := s.repository.ListProjects(ctx, userID, role)
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
	role string,
) (*dto.ProcessPlanning, error) {
	projects, err := s.repository.ListProjects(ctx, userID, role)
	if err != nil {
		return nil, err
	}
	projectIDs := make([]int64, len(projects))
	for i, project := range projects {
		projectIDs[i] = project.ID
	}
	processes, err := s.repository.ListProcessesByProjectIDs(ctx, projectIDs)
	if err != nil {
		return nil, err
	}

	planning := dto.ProcessPlanning{
		Projects: make([]dto.DetailedProject, len(projects)),
	}

	for i, p := range projects {
		planning.Projects[i] = dto.DetailedProject{
			Project:   p,
			Processes: processes[p.ID],
		}
	}

	return &planning, nil
}

func (s *PlanningService) GetTaskPlanning(
	ctx context.Context,
	userID int64,
	role string,
) (*dto.TaskPlanning, error) {
	// 1. Получаем процессы
	processes, err := s.repository.ListProcesses(ctx, userID, role)
	if err != nil {
		return nil, err
	}

	if len(processes) == 0 {
		return &dto.TaskPlanning{
			Processes: []dto.DetailedProcess{},
		}, nil
	}

	// 2. Получаем ID процессов
	processIDs := make([]int64, len(processes))
	for i, p := range processes {
		processIDs[i] = p.ID
	}

	// 3. Получаем связанные данные
	milestones, err := s.repository.ListMilestonesByProcessIDs(ctx, processIDs)
	if err != nil {
		return nil, err
	}

	tasks, err := s.repository.ListTasksByProcessIDs(ctx, processIDs)
	if err != nil {
		return nil, err
	}

	// 4. Собираем ID всех задач
	taskIDs := make([]int64, 0)
	for _, taskList := range tasks {
		for _, task := range taskList {
			taskIDs = append(taskIDs, task.ID)
		}
	}

	// 5. Получаем назначения
	assignments, err := s.repository.ListAssignmentsByTaskIDs(ctx, taskIDs)
	if err != nil {
		return nil, err
	}

	// 6. Получаем все ресурсы и делаем мапу
	resourcesList, err := s.repository.ListResources(ctx)
	if err != nil {
		return nil, err
	}

	resourcesMap := make(map[int64]dto.Resource, len(resourcesList))
	for _, res := range resourcesList {
		resourcesMap[res.ID] = res
	}

	// 7. Собираем результат
	planning := &dto.TaskPlanning{
		Processes: make([]dto.DetailedProcess, 0, len(processes)),
	}

	for _, process := range processes {
		// Безопасно получаем задачи для процесса
		processTasks := tasks[process.ID] // если нет → nil
		if processTasks == nil {
			processTasks = []dto.Task{} // явно создаем пустой слайс
		}

		// Безопасно получаем вехи для процесса
		processMilestones := milestones[process.ID]
		if processMilestones == nil {
			processMilestones = []dto.Milestone{}
		}

		// Собираем детальные задачи
		detailedTasks := make([]dto.DetailedTask, 0, len(processTasks))
		for _, task := range processTasks {
			// Безопасно получаем назначения для задачи
			taskAssignments := assignments[task.ID]
			if taskAssignments == nil {
				taskAssignments = []dto.Assignment{}
			}

			// Собираем ресурсы
			resources := make([]dto.Resource, 0, len(taskAssignments))
			for _, assignment := range taskAssignments {
				// БЕЗОПАСНО получаем ресурс
				res, exists := resourcesMap[assignment.ResourceID]
				if !exists {
					s.logger.Warn("resource not found",
						"resource_id", assignment.ResourceID,
						"task_id", task.ID,
					)
					continue // пропускаем отсутствующий ресурс
				}

				resources = append(resources, dto.Resource{
					ID:       assignment.ResourceID,
					Title:    res.Title,
					Code:     res.Code,
					Quantity: assignment.Quantity,
				})
			}

			detailedTasks = append(detailedTasks, dto.DetailedTask{
				Task:      task,
				Resources: resources,
			})
		}

		planning.Processes = append(planning.Processes, dto.DetailedProcess{
			Process:    process,
			Milestones: processMilestones,
			Tasks:      detailedTasks,
		})
	}

	return planning, nil
}
