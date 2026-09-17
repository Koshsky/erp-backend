package service

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/planning/dto"
	"github.com/Koshsky/erp-backend/internal/planning/repository/sqlc"
)

// getSlice returns the slice associated with the given key in the map.
func getSlice[T any](m map[int64][]T, key int64) []T {
	if val, exists := m[key]; exists {
		return val
	}
	return []T{}
}

// loadProcesses loads processes for the given user ID and view scope code
// using the TASK scope semantics (parent = "in my processes"): the caller is
// the task-planning aggregate, and the same process list is what a process
// owner (vp) must see on the task diagram. The process aggregate uses the
// process scope directly (parent = "in my projects").
func (s *PlanningService) loadProcesses(ctx context.Context, userID int64, viewScope string) ([]dto.Process, error) {
	rows, err := s.repository.ListProcessesByTaskScope(ctx, userID, viewScope)
	if err != nil {
		return nil, err
	}
	processes := make([]dto.Process, len(rows))
	for i, row := range rows {
		processes[i] = toTaskScopeProcess(row)
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

	milestoneRows, err := s.repository.ListMilestonesByProcessIDs(ctx, processIDs)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	milestones := groupByKey(milestoneRows, func(m sqlc.Milestone) int64 { return m.ProcessID })
	milestoneDTOs := make(map[int64][]dto.Milestone, len(milestones))
	for processID, rows := range milestones {
		items := make([]dto.Milestone, len(rows))
		for i, row := range rows {
			items[i] = toMilestone(row)
		}
		milestoneDTOs[processID] = items
	}

	taskRows, err := s.repository.ListTasksByProcessIDs(ctx, processIDs)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	tasks := groupByKey(taskRows, func(t sqlc.Task) int64 { return t.ProcessID })
	taskDTOs := make(map[int64][]dto.Task, len(tasks))
	for processID, rows := range tasks {
		items := make([]dto.Task, len(rows))
		for i, row := range rows {
			items[i] = toTask(row)
		}
		taskDTOs[processID] = items
	}

	taskIDs := s.collectTaskIDs(taskDTOs)

	assignmentRows, err := s.repository.ListAssignmentsByTaskIDs(ctx, taskIDs)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	assignments := groupByKey(assignmentRows, func(a sqlc.Assignment) int64 { return a.TaskID })
	assignmentDTOs := make(map[int64][]dto.Assignment, len(assignments))
	for taskID, rows := range assignments {
		items := make([]dto.Assignment, len(rows))
		for i, row := range rows {
			items[i] = toAssignment(row)
		}
		assignmentDTOs[taskID] = items
	}

	commentCounts, err := s.repository.ListTaskCommentCountsByTaskIDs(ctx, taskIDs)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	resourceRows, err := s.repository.ListResources(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	resources := make([]dto.Resource, len(resourceRows))
	for i, row := range resourceRows {
		resources[i] = toResource(row)
	}
	resourcesMap := s.buildResourceMap(resources)

	return milestoneDTOs, taskDTOs, assignmentDTOs, resourcesMap, commentCounts, nil
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

// buildDetailedTasks constructs detailed tasks with their resources, nesting
// subtasks (operations) under their parent task: only top-level tasks appear
// in the returned slice, each carrying its subtasks in display order.
func (s *PlanningService) buildDetailedTasks(
	tasks []dto.Task,
	assignments map[int64][]dto.Assignment,
	resourcesMap map[int64]dto.Resource,
	commentCounts map[int64]int64,
) []dto.DetailedTask {
	// Build every detailed task first, then group subtasks by parent.
	detailed := make([]dto.DetailedTask, 0, len(tasks))
	byParent := make(map[int64][]dto.DetailedTask, len(tasks))
	for _, task := range tasks {
		taskAssignments := getSlice(assignments, task.ID)
		resources := s.buildTaskResources(taskAssignments, resourcesMap)

		item := dto.DetailedTask{
			Task:          task,
			Resources:     resources,
			CommentsCount: commentCounts[task.ID],
		}
		if task.ParentID != nil {
			byParent[*task.ParentID] = append(byParent[*task.ParentID], item)
			continue
		}
		detailed = append(detailed, item)
	}

	// Attach children to their parents (input is already sorted by the
	// repository, so appended children keep their display order).
	for i := range detailed {
		if children := byParent[detailed[i].ID]; len(children) > 0 {
			detailed[i].Subtasks = children
		}
	}

	return detailed
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
