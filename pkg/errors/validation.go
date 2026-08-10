package errors

import (
	"fmt"
	"net/http"
)

// FieldError is a structured validation error carrying the offending field, a
// machine-readable code and a human-readable message. It maps to HTTP 400 and
// wraps ErrValidation so sentinel matching keeps working.
type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewFieldError builds a structured validation error.
func NewFieldError(field, code, message string) *FieldError {
	return &FieldError{Field: field, Code: code, Message: message}
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", ErrValidation.Error(), e.Message)
}

func (e *FieldError) Unwrap() error {
	return ErrValidation
}

func (e *FieldError) StatusCode() int {
	return http.StatusBadRequest
}

func (e *FieldError) ErrorCode() Code {
	return CodeValidation
}

// NewValidationError is a generic fallback for validation errors that are not
// tied to a single field. Prefer NewFieldError with a field-scoped message.
func NewValidationError(message string) error {
	return fmt.Errorf("%w: %s", ErrValidation, message)
}
