package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func writeBindError(c *gin.Context, logger *slog.Logger, err error) {
	logger.Warn("invalid request payload", "error", err)
	c.JSON(http.StatusBadRequest, response{Error: "invalid request payload"})
}
