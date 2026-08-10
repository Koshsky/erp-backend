// Package errors is the single taxonomy of domain errors. Errors are
// self-describing: they carry their HTTP status (StatusCode), a machine
// readable code (ErrorCode) and stay compatible with sentinel matching via
// wrapped causes. DomainError doubles as the serialized error envelope
// ({ code, message, timestamp }).
package errors

import (
	"net/http"
	"time"

	stderrors "errors"
)

// DomainError is an error that knows its HTTP status, machine-readable code
// and may wrap a cause. It is also the API error body: only Code, Message and
// Timestamp are serialized.
type DomainError struct {
	Code      Code   `json:"code" example:"1"`
	Message   string `json:"message"   example:"Some error message"`
	Timestamp string `json:"timestamp" example:"2026-08-09T10:30:00Z"`

	Status int   `json:"-"`
	Cause  error `json:"-"`
}

func (e *DomainError) Error() string   { return e.Message }
func (e *DomainError) Unwrap() error   { return e.Cause }
func (e *DomainError) StatusCode() int { return e.Status }
func (e *DomainError) ErrorCode() Code { return e.Code }

// Base sentinels used as the Cause of constructed errors so that
// sentinel matching against ErrNotFound and friends keeps working.
var (
	ErrNotFound   = stderrors.New("not found")
	ErrForbidden  = stderrors.New("forbidden")
	ErrBadRequest = stderrors.New("bad request")
	ErrValidation = stderrors.New("validation failed")
)

// now returns the RFC3339 UTC timestamp used in error envelopes.
func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// NotFound returns a 404 error wrapping ErrNotFound.
func NotFound(msg string) error {
	return &DomainError{Status: http.StatusNotFound, Message: msg, Cause: ErrNotFound, Code: CodeNotFound, Timestamp: now()}
}

// Forbidden returns a 403 error wrapping ErrForbidden.
func Forbidden(msg string) error {
	return &DomainError{Status: http.StatusForbidden, Message: msg, Cause: ErrForbidden, Code: CodeForbidden, Timestamp: now()}
}

// BadRequest returns a 400 error wrapping ErrBadRequest.
func BadRequest(msg string) error {
	return &DomainError{Status: http.StatusBadRequest, Message: msg, Cause: ErrBadRequest, Code: CodeBadRequest, Timestamp: now()}
}

// Entity-specific not-found errors for stable sentinel matching in services.
var (
	ErrProjectNotFound    = NotFound("project not found")
	ErrProcessNotFound    = NotFound("process not found")
	ErrMilestoneNotFound  = NotFound("milestone not found")
	ErrTaskNotFound       = NotFound("task not found")
	ErrResourceNotFound   = NotFound("resource not found")
	ErrAssignmentNotFound = NotFound("assignment not found")
	ErrUserNotFound       = NotFound("user not found")
	ErrStateNotFound      = NotFound("state not found")
	ErrEmployeeNotFound   = NotFound("employee not found")
)
