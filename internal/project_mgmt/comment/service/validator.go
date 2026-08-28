package service

import (
	"fmt"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/comment/domain"
	"github.com/Koshsky/erp-backend/pkg/errors"
	"github.com/Koshsky/erp-backend/pkg/validator"
)

// maxCommentContentLen is the maximum comment text length (in characters).
const maxCommentContentLen = 4000

type CommentValidator struct {
	validator.Validator
}

func (v *CommentValidator) ValidateComment(c *domain.Comment) error {
	if err := v.ValidatePositiveID(c.TaskID, "task_id"); err != nil {
		return err
	}
	if err := v.ValidatePositiveID(c.AuthorID, "author_id"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(c.Content, "content"); err != nil {
		return err
	}
	if len([]rune(c.Content)) > maxCommentContentLen {
		return errors.NewFieldError(
			"content",
			"max_length",
			fmt.Sprintf("content must not exceed %d characters", maxCommentContentLen),
		)
	}
	return nil
}
