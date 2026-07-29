package response

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response wraps API responses into { "data": ..., "error": ... } format.
// Fields use pointer/omitempty so that when Data is nil it's omitted,
// and when Error is empty it's omitted.
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// OK sends a 200 response with the provided data.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Data: data})
}

// Created sends a 201 response with the provided data.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{Data: data})
}

// NoContent sends a 204 response with no body.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error sends a JSON error response with the given status code and message.
func Error(c *gin.Context, status int, msg string) {
	c.JSON(status, Response{Error: msg})
}

// BadRequest sends a 400 error response.
func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Response{Error: msg})
}

// InternalError sends a 500 error response and logs the internal error.
func InternalError(c *gin.Context, logger *slog.Logger, msg string, err error) {
	logger.Error(msg, "error", err)
	c.JSON(http.StatusInternalServerError, Response{Error: msg})
}

// HandleBindError is a convenience helper for JSON binding errors.
// It logs the original error at warn level and responds with 400.
func HandleBindError(c *gin.Context, logger *slog.Logger, err error) {
	logger.Warn("invalid request payload", "error", err)
	BadRequest(c, "invalid request payload")
}

// Forbidden sends a 403 error response.
func Forbidden(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, Response{Error: msg})
}

// NotFound sends a 404 error response.
func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, Response{Error: msg})
}

