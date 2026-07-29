package delivery

import (
	"log/slog"
	"strconv"

	"github.com/Koshsky/erp-backend/internal/common/response"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/dto"
	"github.com/gin-gonic/gin"
)

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

// @Tags Assignments
// @Summary List assignments
// @Description Get a list of all assignments
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=[]dto.AssignmentResponse}
// @Failure 500 {object} response.Response
// @Router /assignments [get]
func (h *AssignmentHandler) ListAssignments(c *gin.Context) {
	assignments, err := h.service.ListAssignments(c.Request.Context())
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, assignments)
}

// @Tags Assignments
// @Summary Get assignment
// @Description Get assignment by id
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Assignment ID"
// @Success 200 {object} response.Response{data=dto.AssignmentResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /assignments/{id} [get]
func (h *AssignmentHandler) FindAssignment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid assignment id")
		return
	}

	assignment, err := h.service.FindAssignment(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, assignment)
}

// @Tags Assignments
// @Summary Create assignment
// @Description Create a new assignment
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateAssignmentRequest true "Assignment data"
// @Success 201 {object} response.Response{data=dto.AssignmentResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /assignments [post]
func (h *AssignmentHandler) CreateAssignment(c *gin.Context) {
	var assignment dto.CreateAssignmentRequest
	if err := c.ShouldBindJSON(&assignment); err != nil {
		response.HandleBindError(c, h.logger, err)
		return
	}

	created, err := h.service.CreateAssignment(c.Request.Context(), assignment)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.Created(c, created)
}

// @Tags Assignments
// @Summary Delete an assignment
// @Description Delete an assignment by ID
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Assignment ID"
// @Success 204
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /assignments/{id} [delete]
func (h *AssignmentHandler) DeleteAssignment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid assignment id")
		return
	}

	if err := h.service.DeleteAssignment(c.Request.Context(), id); err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.NoContent(c)
}

// @Tags Assignments
// @Summary Update an assignment
// @Description Update an assignment by ID
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Assignment ID"
// @Param body body dto.UpdateAssignmentRequest true "Assignment data"
// @Success 200 {object} response.Response{data=dto.AssignmentResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /assignments/{id} [put]
func (h *AssignmentHandler) UpdateAssignment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid assignment id")
		return
	}

	body := dto.UpdateAssignmentRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.HandleBindError(c, h.logger, err)
		return
	}

	assignment, err := h.service.UpdateAssignment(c.Request.Context(), id, body)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, assignment)
}
