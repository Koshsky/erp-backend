package service

import (
	"fmt"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/task_dependency/domain"
	"github.com/Koshsky/erp-backend/pkg/date"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

// boundDate returns the date of the successor that the link constrains: its
// start for fs/ss, its end for ff/sf.
func boundDate(typ string, taskStart, taskEnd date.Date) date.Date {
	if domain.BoundIsStart(typ) {
		return taskStart
	}
	return taskEnd
}

// anchorDate returns the date of the predecessor that anchors the constraint:
// its start for ss/sf, its end for fs/ff.
func anchorDate(typ string, predStart, predEnd date.Date) date.Date {
	if domain.AnchorIsStart(typ) {
		return predStart
	}
	return predEnd
}

// violationMessage builds the user-facing message for a violated link.
func violationMessage(typ, taskTitle, predTitle string) string {
	switch typ {
	case domain.TypeFinishToStart:
		return fmt.Sprintf(
			"задача «%s» должна начинаться не раньше окончания задачи «%s» (связь «окончание → начало»)",
			taskTitle, predTitle,
		)
	case domain.TypeStartToStart:
		return fmt.Sprintf(
			"задача «%s» должна начинаться не раньше начала задачи «%s» (связь «начало → начало»)",
			taskTitle, predTitle,
		)
	case domain.TypeFinishToFinish:
		return fmt.Sprintf(
			"задача «%s» должна заканчиваться не раньше окончания задачи «%s» (связь «окончание → окончание»)",
			taskTitle, predTitle,
		)
	case domain.TypeStartToFinish:
		return fmt.Sprintf(
			"задача «%s» должна заканчиваться не раньше начала задачи «%s» (связь «начало → окончание»)",
			taskTitle, predTitle,
		)
	default:
		return "недопустимый тип связи"
	}
}

// checkConstraint validates one link: the successor's bound date must not
// precede the predecessor's anchor date. A violation is returned as a
// validation error (400) so the conflicting pair is visible to the user.
func checkConstraint(
	taskTitle string,
	taskStart, taskEnd date.Date,
	typ, predTitle string,
	predStart, predEnd date.Date,
) error {
	if boundDate(typ, taskStart, taskEnd) < anchorDate(typ, predStart, predEnd) {
		return errors.NewValidationError(violationMessage(typ, taskTitle, predTitle))
	}
	return nil
}
