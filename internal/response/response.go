package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/pkg/errors"
)

// Response wraps API responses into { "data": ..., "error": ... } format.
// On success error is an empty object; on failure data is an empty object and
// error carries the domain envelope { code, message, timestamp }.
type Response struct {
	Data  any                 `json:"data"`
	Error *errors.DomainError `json:"error"`
}

// SuccessResponse documents the 2xx envelope for swagger: data carries the
// payload, error is always the empty object {}. Handlers annotate
// response.SuccessResponse{data=...,error=nil} so swag generates error as an
// empty object instead of an untyped field rendered as a plain string.
type SuccessResponse struct {
	Data  any `json:"data"`
	Error any `json:"error"`
}

// ErrorResponse is the swagger model for 4xx/5xx responses: data is empty,
// error carries the domain envelope { code, message, timestamp }.
type ErrorResponse struct {
	Data  any                 `json:"data"`
	Error *errors.DomainError `json:"error"`
}

// MarshalJSON renders empty Data/Error as {} instead of null. The Error field
// keeps its DomainError type so swagger documents the envelope; the runtime
// payload uses an empty object when there is no error.
func (r Response) MarshalJSON() ([]byte, error) {
	data := r.Data
	if data == nil {
		data = map[string]any{}
	}
	var errBody any
	if r.Error != nil {
		errBody = r.Error
	} else {
		errBody = map[string]any{}
	}
	return json.Marshal(struct {
		Data  any `json:"data"`
		Error any `json:"error"`
	}{Data: data, Error: errBody})
}

// internalErrorMessage is the stable client-facing message for 5xx responses.
// Internal details are never exposed to the client; they go to the logs only.
const internalErrorMessage = "internal server error"

// errorBody builds an error envelope with the given code, message and now.
func errorBody(code errors.Code, msg string) *errors.DomainError {
	return &errors.DomainError{
		Code:      code,
		Message:   msg,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// OK sends a 200 response with the provided data.
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{Data: data})
}

// Created sends a 201 response with the provided data.
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{Data: data})
}

// NoContent sends a 204 response with no body.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequest sends a 400 error response.
func BadRequest(c *gin.Context, code errors.Code, msg string) {
	c.JSON(http.StatusBadRequest, Response{Error: errorBody(code, msg)})
}

// Unauthorized sends a 401 error response.
func Unauthorized(c *gin.Context, code errors.Code, msg string) {
	c.JSON(http.StatusUnauthorized, Response{Error: errorBody(code, msg)})
}

// InternalError sends a 500 error response and logs the internal error. The
// client always receives the generic internalErrorMessage; msg and err (which
// may contain internal details) are only written to the logs.
func InternalError(c *gin.Context, logger *slog.Logger, msg string, err error) {
	logger.Error(msg, "error", err)
	c.JSON(http.StatusInternalServerError, Response{Error: errorBody(errors.CodeInternal, internalErrorMessage)})
}

// Forbidden sends a 403 error response.
func Forbidden(c *gin.Context, code errors.Code, msg string) {
	c.JSON(http.StatusForbidden, Response{Error: errorBody(code, msg)})
}

// NotFound sends a 404 error response.
func NotFound(c *gin.Context, code errors.Code, msg string) {
	c.JSON(http.StatusNotFound, Response{Error: errorBody(code, msg)})
}

// TooManyRequests sends a 429 error response.
func TooManyRequests(c *gin.Context, msg string) {
	c.JSON(http.StatusTooManyRequests, Response{Error: errorBody(errors.CodeTooManyRequests, msg)})
}

// Error maps a domain error to an HTTP response using its status and code:
// errors that carry StatusCode/ErrorCode (forbidden → 403/FORBIDDEN,
// not found → 404/NOT_FOUND, validation → 400/VALIDATION_ERROR) are honored,
// otherwise the error is classified by cause. Internal errors (500) are logged.
func Error(c *gin.Context, logger *slog.Logger, err error) {
	status := errors.StatusCode(err)
	code := errors.CodeOf(err, status)
	if status == http.StatusInternalServerError {
		logger.Error("internal error", "error", err)
		c.JSON(status, Response{Error: errorBody(code, internalErrorMessage)})
		return
	}
	c.JSON(status, Response{Error: errorBody(code, err.Error())})
}
