package service

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/comment/domain"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment domain.Comment) (*domain.Comment, error)
	FindComment(ctx context.Context, id int64) (*domain.Comment, error)
	DeleteComment(ctx context.Context, id int64) error
	ListComments(ctx context.Context, taskID int64) ([]domain.Comment, error)
}
