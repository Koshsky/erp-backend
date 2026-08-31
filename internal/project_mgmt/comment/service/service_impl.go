package service

import (
	"context"
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/project_mgmt/comment/repository"
	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/comment/domain"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/comment/dto"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

type CommentService struct {
	logger     *slog.Logger
	tracer     *tracingpkg.Tracer
	repository CommentRepository
	mapper     *CommentMapper
	validator  *CommentValidator
}

// NewCommentService builds the CommentService service.
func NewCommentService(
	logger *slog.Logger,
	tracer *tracingpkg.Tracer,
	r *repo.CommentRepository,
) *CommentService {
	return &CommentService{
		logger:     logger,
		tracer:     tracer,
		repository: r,
		mapper:     NewCommentMapper(),
		validator:  &CommentValidator{},
	}
}

// CreateComment создаёт комментарий задачи: автор (authorID) приходит из
// контекста авторизации, parent_id — ответ на комментарий той же задачи.
func (s *CommentService) CreateComment(
	ctx context.Context,
	taskID int64,
	req dto.CreateCommentRequest,
	authorID int64,
) (*dto.CommentResponse, error) {
	ctx, end := s.tracer.Start(ctx, "comment.CreateComment")
	defer end(nil)

	comment := s.mapper.ToDomainFromCreate(taskID, req, authorID)
	if err := s.validator.ValidateComment(&comment); err != nil {
		return nil, err
	}

	if comment.ParentID != nil {
		if _, err := s.findValidParent(ctx, taskID, *comment.ParentID); err != nil {
			return nil, err
		}
	}

	created, err := s.repository.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(created), nil
}

// findValidParent проверяет, что комментарий-родитель существует и принадлежит
// той же задаче (ответы не могут ссылаться на комментарии чужих задач).
func (s *CommentService) findValidParent(ctx context.Context, taskID int64, parentID int64) (*domain.Comment, error) {
	parent, err := s.repository.FindComment(ctx, parentID)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil, errors.BadRequest("parent comment not found")
		}
		return nil, err
	}
	if parent == nil {
		return nil, errors.BadRequest("parent comment not found")
	}
	if parent.TaskID != taskID {
		return nil, errors.BadRequest("parent comment belongs to another task")
	}
	return parent, nil
}

func (s *CommentService) ListComments(ctx context.Context, taskID int64) ([]dto.CommentResponse, error) {
	ctx, end := s.tracer.Start(ctx, "comment.ListComments")
	defer end(nil)

	rows, err := s.repository.ListComments(ctx, taskID)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(rows), nil
}

func (s *CommentService) DeleteComment(ctx context.Context, commentID int64) error {
	ctx, end := s.tracer.Start(ctx, "comment.DeleteComment")
	defer end(nil)

	comment, err := s.repository.FindComment(ctx, commentID)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil // идемпотентный delete: уже удалено — не ошибка
		}
		return err
	}
	if comment == nil {
		return nil // идемпотентный delete
	}

	return s.repository.DeleteComment(ctx, commentID)
}
