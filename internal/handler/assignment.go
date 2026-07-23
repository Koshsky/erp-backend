package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Koshsky/erp/api/internal/dto"
	"github.com/gin-gonic/gin"
)

type AssignmentService interface {
	GetAssignment(ctx context.Context, id int64) (*dto.AssignmentResponse, error)
	CreateAssignment(ctx context.Context, assignment dto.CreateAssignmentRequest) (*dto.AssignmentResponse, error)
	DeleteAssignment(ctx context.Context, id int64) error
	UpdateAssignment(ctx context.Context, id int64, assignment dto.UpdateAssignmentRequest) (*dto.AssignmentResponse, error)
}

type AssignmentHandler struct {
	logger  *slog.Logger
	service AssignmentService
}

func NewAssignmentHandler(logger *slog.Logger, service AssignmentService) *AssignmentHandler {
	return &AssignmentHandler{
		logger:  logger,
		service: service,
	}
}

func (h *AssignmentHandler) GetAssignment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid assignment id"})
		return
	}

	assignment, err := h.service.GetAssignment(c.Request.Context(), id)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, response{Error: "assignment not found"})
			return
		}
		h.logger.Error("failed to get assignment", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: assignment})
}

func (h *AssignmentHandler) CreateAssignment(c *gin.Context) {
	var assignment dto.CreateAssignmentRequest
	if err := c.ShouldBindJSON(&assignment); err != nil {
		writeBindError(c, h.logger, err)
		return
	}

	created, err := h.service.CreateAssignment(c.Request.Context(), assignment)
	if err != nil {
		if isValidationError(err) {
			c.JSON(http.StatusBadRequest, response{Error: err.Error()})
			return
		}
		h.logger.Error("failed to create assignment", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response{Data: created})
}

func (h *AssignmentHandler) DeleteAssignment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid assignment id"})
		return
	}

	if err := h.service.DeleteAssignment(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete assignment", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *AssignmentHandler) UpdateAssignment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid assignment id"})
		return
	}

	body := dto.UpdateAssignmentRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		writeBindError(c, h.logger, err)
		return
	}

	assignment, err := h.service.UpdateAssignment(c.Request.Context(), id, body)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, response{Error: "assignment not found"})
			return
		}
		if isValidationError(err) {
			c.JSON(http.StatusBadRequest, response{Error: err.Error()})
			return
		}
		h.logger.Error("failed to update assignment", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: assignment})
}
