package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/dto"
)

type TaskService interface {
	ListTasks(ctx context.Context) ([]dto.TaskResponse, error)
	FindTask(ctx context.Context, id int64) (*dto.TaskResponse, error)
	CreateTask(ctx context.Context, task dto.CreateTaskRequest) (*dto.TaskResponse, error)
	DeleteTask(ctx context.Context, id int64) error
	UpdateTask(ctx context.Context, id int64, task dto.UpdateTaskRequest) (*dto.TaskResponse, error)
}
