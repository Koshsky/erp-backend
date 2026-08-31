package service

import (
	"github.com/Koshsky/erp-backend/internal/project_mgmt/comment/domain"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/comment/dto"
)

type CommentMapper struct{}

func NewCommentMapper() *CommentMapper {
	return &CommentMapper{}
}

func (m *CommentMapper) ToDTO(comment *domain.Comment) *dto.CommentResponse {
	if comment == nil {
		return nil
	}
	return &dto.CommentResponse{
		ID:        comment.ID,
		TaskID:    comment.TaskID,
		AuthorID:  comment.AuthorID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}
}

func (m *CommentMapper) ToDTOs(comments []domain.Comment) []dto.CommentResponse {
	if comments == nil {
		return []dto.CommentResponse{}
	}

	responses := make([]dto.CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = *m.ToDTO(&comment)
	}
	return responses
}

// ToDomainFromCreate builds a comment from the request: the author always
// comes from the authorization context (authorID), not from the request body.
func (m *CommentMapper) ToDomainFromCreate(taskID int64, req dto.CreateCommentRequest, authorID int64) domain.Comment {
	return domain.Comment{
		TaskID:   taskID,
		AuthorID: authorID,
		ParentID: req.ParentID,
		Content:  req.Content,
	}
}
