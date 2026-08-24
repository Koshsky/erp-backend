//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/comment/domain"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/comment/repository/sqlc"
	nullable "github.com/Koshsky/erp-backend/pkg/database"
)

type CommentRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

// NewCommentRepository builds the CommentRepository repository.
func NewCommentRepository(logger *slog.Logger, pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *CommentRepository) CreateComment(ctx context.Context, comment domain.Comment) (*domain.Comment, error) {
	row, err := r.db.CreateComment(ctx, sqlc.CreateCommentParams{
		TaskID:   comment.TaskID,
		AuthorID: comment.AuthorID,
		ParentID: nullable.ToInt8(comment.ParentID),
		Content:  comment.Content,
	})
	if err != nil {
		return nil, err
	}

	mapped := mapComment(row)
	return &mapped, nil
}

func (r *CommentRepository) FindComment(ctx context.Context, id int64) (*domain.Comment, error) {
	row, err := r.db.FindComment(ctx, id)
	if err != nil {
		return nil, err
	}

	mapped := mapComment(row)
	return &mapped, nil
}

func (r *CommentRepository) DeleteComment(ctx context.Context, id int64) error {
	return r.db.DeleteComment(ctx, id)
}

func (r *CommentRepository) ListComments(ctx context.Context, taskID int64) ([]domain.Comment, error) {
	rows, err := r.db.ListComments(ctx, taskID)
	if err != nil {
		return nil, err
	}
	comments := make([]domain.Comment, 0, len(rows))
	for _, row := range rows {
		comments = append(comments, mapComment(row))
	}
	return comments, nil
}

func mapComment(row sqlc.TaskComment) domain.Comment {
	return domain.Comment{
		ID:        row.ID,
		TaskID:    row.TaskID,
		AuthorID:  row.AuthorID,
		ParentID:  nullable.Int64Ptr(row.ParentID),
		Content:   row.Content,
		CreatedAt: row.CreatedAt,
	}
}

// OwnerChain returns the owner chain of a comment (task owners + the author).
func (r *CommentRepository) OwnerChain(ctx context.Context, id int64) (rbac.Owners, error) {
	row, err := r.db.OwnerChain(ctx, id)
	if err != nil {
		return rbac.Owners{}, err
	}
	return rbac.Owners{
		ProjectOwner: row.ProjectOwner,
		ProcessOwner: row.ProcessOwner,
		Owner:        row.Author,
	}, nil
}
