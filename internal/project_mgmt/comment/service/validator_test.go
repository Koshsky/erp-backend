package service_test

import (
	"strings"
	"testing"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/comment/domain"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/comment/service"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

// maxCommentContentLen mirrors the service limit (see comment/service/validator.go).
const maxCommentContentLen = 4000

func TestCommentValidator(t *testing.T) {
	t.Parallel()
	v := &service.CommentValidator{}

	valid := domain.Comment{TaskID: 1, AuthorID: 2, Content: "Перенести сроки?"}

	cases := []struct {
		name    string
		comment domain.Comment
		wantErr bool
	}{
		{"valid root comment", valid, false},
		{"zero task id", domain.Comment{TaskID: 0, AuthorID: 2, Content: "x"}, true},
		{"zero author id", domain.Comment{TaskID: 1, AuthorID: 0, Content: "x"}, true},
		{"empty content", domain.Comment{TaskID: 1, AuthorID: 2, Content: "   "}, true},
		{
			"max length is allowed",
			domain.Comment{TaskID: 1, AuthorID: 2, Content: strings.Repeat("x", maxCommentContentLen)},
			false,
		},
		{
			"over max length",
			domain.Comment{TaskID: 1, AuthorID: 2, Content: strings.Repeat("x", maxCommentContentLen+1)},
			true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := v.ValidateComment(&tc.comment)
			if tc.wantErr != (err != nil) {
				t.Fatalf("ValidateComment() error = %v, wantErr %v", err, tc.wantErr)
			}
			if err != nil && !errors.IsValidationError(err) {
				t.Fatalf("ValidateComment() expected validation error, got %v", err)
			}
		})
	}
}
