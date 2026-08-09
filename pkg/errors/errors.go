// Package errors is the single taxonomy of domain errors. Errors are
// self-describing: they carry their HTTP status (StatusCode) and stay
// compatible with sentinel matching via wrapped causes.
package errors

import (
	"net/http"

	stderrors "errors"
)

// DomainError is an error that knows its HTTP status and may wrap a cause.
type DomainError struct {
	Status int
	Msg    string
	Cause  error
}

func (e *DomainError) Error() string   { return e.Msg }
func (e *DomainError) Unwrap() error   { return e.Cause }
func (e *DomainError) StatusCode() int { return e.Status }

// Base sentinels used as the Cause of constructed errors so that
// sentinel matching against ErrNotFound and friends keeps working.
var (
	ErrNotFound   = stderrors.New("not found")
	ErrForbidden  = stderrors.New("forbidden")
	ErrBadRequest = stderrors.New("bad request")
	ErrValidation = stderrors.New("validation failed")
)

// NotFound returns a 404 error wrapping ErrNotFound.
func NotFound(msg string) error {
	return &DomainError{Status: http.StatusNotFound, Msg: msg, Cause: ErrNotFound}
}

// Forbidden returns a 403 error wrapping ErrForbidden.
func Forbidden(msg string) error {
	return &DomainError{Status: http.StatusForbidden, Msg: msg, Cause: ErrForbidden}
}

// BadRequest returns a 400 error wrapping ErrBadRequest.
func BadRequest(msg string) error {
	return &DomainError{Status: http.StatusBadRequest, Msg: msg, Cause: ErrBadRequest}
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
