package policies

import (
	"strconv"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

//nolint:gochecknoglobals // rule registry
var commentPolicies = []rbac.Policy{
	// Чтение и создание комментариев — всем, кто видит задачу (та же матрица,
	// что и task.view: недоступная задача не раскрывается — 404).
	{Name: "task.comment.list", Check: EntityCheck(rbac.ResourceTask, ActionView)},
	{Name: "task.comment.create", Check: EntityCheck(rbac.ResourceTask, ActionView)},
	// Удаление: автор комментария — всегда; остальным — по праву обновления
	// задачи (admin, vp своего процесса).
	{Name: "task.comment.delete", Check: CommentDeleteCheck()},
}

// CommentDeleteCheck разрешает удаление автору комментария; остальным — по
// праву обновления задачи (admin, vp своего процесса).
func CommentDeleteCheck() func(*rbac.CheckCtx) error {
	return func(rc *rbac.CheckCtx) error {
		// В URL две сущности: /task/:id/comments/:comment_id — id комментария в :comment_id.
		commentID, err := strconv.ParseInt(rc.C.Param("comment_id"), 10, 64)
		if err != nil {
			return errors.BadRequest("invalid comment id")
		}

		owners, err := rc.Owners(rbac.ResourceComment, commentID)
		if err != nil {
			return err
		}
		if owners.Owner != 0 && owners.Owner == rc.User.ID {
			return nil // автор удаляет своё
		}
		if !Authorize(rc.User.Role, rbac.ResourceTask, ActionUpdate, owners, rc.User.ID) {
			return errors.ErrForbidden
		}
		return nil
	}
}
