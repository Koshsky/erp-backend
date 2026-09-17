package service

import (
	"fmt"

	"github.com/Koshsky/erp-backend/pkg/errors"
	"github.com/Koshsky/erp-backend/pkg/validator"
)

// maxCommentContentLen is the maximum comment text length (in characters).
const maxCommentContentLen = 4000

type CommentValidator struct {
	validator.Validator
}

func (v *CommentValidator) ValidateComment(taskID, authorID int64, content string) error {
	if err := v.ValidatePositiveID(taskID, "task_id"); err != nil {
		return err
	}
	if err := v.ValidatePositiveID(authorID, "author_id"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(content, "content"); err != nil {
		return err
	}
	if len([]rune(content)) > maxCommentContentLen {
		return errors.NewFieldError(
			"content",
			"max_length",
			fmt.Sprintf("content must not exceed %d characters", maxCommentContentLen),
		)
	}
	return nil
}
