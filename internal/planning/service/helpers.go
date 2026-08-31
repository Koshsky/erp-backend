package service

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/planning/dto"
)

// getSlice returns the slice associated with the given key in the map.
func getSlice[T any](m map[int64][]T, key int64) []T {
	if val, exists := m[key]; exists {
		return val
	}
	return []T{}
}

// loadProcesses loads processes for the given user ID and role.
func (s *PlanningService) loadProcesses(ctx context.Context, userID int64, role string) ([]dto.Process, error) {
	processes, err := s.repository.ListProcesses(ctx, userID, role)
	if err != nil {
		return nil, err
	}
	return processes, nil
}

// loadAllData load milestones, tasks, assignments, resources and comment counts
// for the given processes.
func (s *PlanningService) loadAllData(
	ctx context.Context,
	processes []dto.Process,
) (map[int64][]dto.Milestone, map[int64][]dto.Task, map[int64][]dto.Assignment, map[int64]dto.Resource, map[int64]int64, error) {
	processIDs := make([]int64, len(processes))
	for i, p := range processes {
		processIDs[i] = p.ID
	}

	milestones, err := s.repository.ListMilestonesByProcessIDs(ctx, processIDs)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	tasks, err := s.repository.ListTasksByProcessIDs(ctx, processIDs)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	taskIDs := s.collectTaskIDs(tasks)

	assignments, err := s.repository.ListAssignmentsByTaskIDs(ctx, taskIDs)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	commentCounts, err := s.repository.ListTaskCommentCountsByTaskIDs(ctx, taskIDs)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	resourcesList, err := s.repository.ListResources(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	resourcesMap := s.buildResourceMap(resourcesList)

	return milestones, tasks, assignments, resourcesMap, commentCounts, nil
}

// collectTaskIDs collects all task IDs from the tasks map.
func (s *PlanningService) collectTaskIDs(tasks map[int64][]dto.Task) []int64 {
	taskIDs := make([]int64, 0)
	for _, taskList := range tasks {
		for _, task := range taskList {
			taskIDs = append(taskIDs, task.ID)
		}
	}
	return taskIDs
}

// buildResourceMap creates a map of resources by their IDs.
func (s *PlanningService) buildResourceMap(resources []dto.Resource) map[int64]dto.Resource {
	resourceMap := make(map[int64]dto.Resource, len(resources))
	for _, res := range resources {
		resourceMap[res.ID] = res
	}
	return resourceMap
}

// buildPlanning constructs the final task planning structure.
func (s *PlanningService) buildPlanning(
	processes []dto.Process,
	milestones map[int64][]dto.Milestone,
	tasks map[int64][]dto.Task,
	assignments map[int64][]dto.Assignment,
	resourcesMap map[int64]dto.Resource,
	commentCounts map[int64]int64,
) *dto.TaskPlanning {
	planning := &dto.TaskPlanning{
		Processes: make([]dto.DetailedProcess, 0, len(processes)),
	}

	for _, process := range processes {
		detailedProcess := s.buildDetailedProcess(
			process,
			milestones,
			tasks,
			assignments,
			resourcesMap,
			commentCounts,
		)
		planning.Processes = append(planning.Processes, detailedProcess)
	}

	return planning
}

// buildDetailedProcess constructs a detailed process with its tasks and milestones.
func (s *PlanningService) buildDetailedProcess(
	process dto.Process,
	milestones map[int64][]dto.Milestone,
	tasks map[int64][]dto.Task,
	assignments map[int64][]dto.Assignment,
	resourcesMap map[int64]dto.Resource,
	commentCounts map[int64]int64,
) dto.DetailedProcess {
	// Use the generic getSlice function
	processTasks := getSlice(tasks, process.ID)
	processMilestones := getSlice(milestones, process.ID)
	detailedTasks := s.buildDetailedTasks(processTasks, assignments, resourcesMap, commentCounts)

	return dto.DetailedProcess{
		Process:    process,
		Milestones: processMilestones,
		Tasks:      detailedTasks,
	}
}

// buildDetailedTasks constructs detailed tasks with their resources.
func (s *PlanningService) buildDetailedTasks(
	tasks []dto.Task,
	assignments map[int64][]dto.Assignment,
	resourcesMap map[int64]dto.Resource,
	commentCounts map[int64]int64,
) []dto.DetailedTask {
	detailedTasks := make([]dto.DetailedTask, 0, len(tasks))

	for _, task := range tasks {
		taskAssignments := getSlice(assignments, task.ID)
		resources := s.buildTaskResources(taskAssignments, resourcesMap)

		detailedTasks = append(detailedTasks, dto.DetailedTask{
			Task:          task,
			Resources:     resources,
			CommentsCount: commentCounts[task.ID],
		})
	}

	return detailedTasks
}

// buildTaskResources constructs resources for a task.
func (s *PlanningService) buildTaskResources(
	assignments []dto.Assignment,
	resourcesMap map[int64]dto.Resource,
) []dto.Resource {
	resources := make([]dto.Resource, 0, len(assignments))

	for _, assignment := range assignments {
		res, exists := resourcesMap[assignment.ResourceID]
		if !exists {
			continue
		}

		resources = append(resources, dto.Resource{
			ID:           assignment.ResourceID,
			Title:        res.Title,
			Code:         res.Code,
			Color:        res.Color,
			Quantity:     assignment.Quantity,
			AssignmentID: assignment.ID,
		})
	}

	return resources
}
