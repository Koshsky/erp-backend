package validator

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrValidation = errors.New("validation failed")

	ErrProjectNotFound    = fmt.Errorf("%w: project", ErrNotFound)
	ErrProcessNotFound    = fmt.Errorf("%w: process", ErrNotFound)
	ErrMilestoneNotFound  = fmt.Errorf("%w: milestone", ErrNotFound)
	ErrTaskNotFound       = fmt.Errorf("%w: task", ErrNotFound)
	ErrResourceNotFound   = fmt.Errorf("%w: resource", ErrNotFound)
	ErrAssignmentNotFound = fmt.Errorf("%w: assignment", ErrNotFound)
	ErrUserNotFound       = fmt.Errorf("%w: user", ErrNotFound)

	ErrInvalidAssignmentQuantity = NewFieldError("quantity", codeMinValue, msgAtLeast("quantity", 1))
)

// FieldError is a structured validation error carrying the offending field,
// a machine-readable code and a human-readable message. It wraps ErrValidation
// so errors.Is(err, ErrValidation) and IsValidationError keep working unchanged.
// Field/Code are not yet exposed over the API, but are kept ready for a future
// structured error response payload.
type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", ErrValidation.Error(), e.Message)
}

func (e *FieldError) Unwrap() error {
	return ErrValidation
}

func NewFieldError(field, code, message string) *FieldError {
	return &FieldError{Field: field, Code: code, Message: message}
}

// NewValidationError is a generic fallback constructor for validation errors
// that are not tied to a single field. Prefer NewFieldError with a message
// template from validation_messages.go for field-scoped errors.
func NewValidationError(message string) error {
	return fmt.Errorf("%w: %s", ErrValidation, message)
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound) ||
		errors.Is(err, pgx.ErrNoRows) ||
		errors.Is(err, sql.ErrNoRows)
}

func IsValidationError(err error) bool {
	if errors.Is(err, ErrValidation) {
		return true
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	// 23505 unique_violation, 23503 foreign_key_violation, 23514 check_violation
	return pgErr.Code == "23505" || pgErr.Code == "23503" || pgErr.Code == "23514"
}
