package delivery

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Koshsky/erp/api/internal/assignment/dto"
	"github.com/gin-gonic/gin"
)

// TODO: move response to common DTO

type response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func writeBindError(c *gin.Context, logger *slog.Logger, err error) {
	logger.Warn("invalid request payload", "error", err)
	c.JSON(http.StatusBadRequest, response{Error: "invalid request payload"})
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

func (h *AssignmentHandler) ListAssignments(c *gin.Context) {
	assignments, err := h.service.ListAssignments(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list assignments", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: assignments})
}

func (h *AssignmentHandler) GetAssignment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid assignment id"})
		return
	}

	assignment, err := h.service.GetAssignment(c.Request.Context(), id)
	if err != nil {
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
		h.logger.Error("failed to update assignment", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: assignment})
}
