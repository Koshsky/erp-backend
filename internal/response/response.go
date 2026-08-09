package response

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/pkg/errors"
)

// Response wraps API responses into { "data": ..., "error": ... } format.
// Fields use pointer/omitempty so that when Data is nil it's omitted,
// and when Error is empty it's omitted.
type Response struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty" example:"error message"`
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
func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Response{Error: msg})
}

// Unauthorized sends a 401 error response.
func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, Response{Error: msg})
}

// InternalError sends a 500 error response and logs the internal error.
func InternalError(c *gin.Context, logger *slog.Logger, msg string, err error) {
	logger.Error(msg, "error", err)
	c.JSON(http.StatusInternalServerError, Response{Error: msg})
}

// Forbidden sends a 403 error response.
func Forbidden(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, Response{Error: msg})
}

// NotFound sends a 404 error response.
func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, Response{Error: msg})
}

// TooManyRequests sends a 429 error response.
func TooManyRequests(c *gin.Context, msg string) {
	c.JSON(http.StatusTooManyRequests, Response{Error: msg})
}

// Error maps a domain error to an HTTP response using its status: errors that
// carry StatusCode (forbidden → 403, not found → 404, validation → 400) are
// honored, otherwise the error is classified by cause. Internal errors (500)
// are logged.
func Error(c *gin.Context, logger *slog.Logger, err error) {
	status := errors.StatusCode(err)
	if status == http.StatusInternalServerError {
		logger.Error("internal error", "error", err)
	}
	c.JSON(status, Response{Error: err.Error()})
}
