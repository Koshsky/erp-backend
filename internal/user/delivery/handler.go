package delivery

import (
	"log/slog"
	"strconv"

	"github.com/Koshsky/erp-backend/internal/common/ctx"
	"github.com/Koshsky/erp-backend/internal/common/response"
	"github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/internal/user/dto"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	logger  *slog.Logger
	service UserService
}

func NewUserHandler(logger *slog.Logger, service UserService) *UserHandler {
	return &UserHandler{
		logger:  logger,
		service: service,
	}
}

// ListUsers
// @Summary List all users
// @Description Returns all users in the system
// @Tags Users
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=[]dto.UserResponse}
// @Failure 500 {object} response.Response
// @Router /user [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.service.ListUsers(c.Request.Context())
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, users)
}

// GetUser
// @Summary User information
// @Description Returns information about a specific user
// @Tags Users
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, user)
}

// @Tags Users
// @Summary Create a new user
// @Description Create a new user with the provided data
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body dto.CreateUserRequest true "User data"
// @Success 201 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	if role := ctx.GetRole(c); role != domain.ProjectDirector {
		response.Forbidden(c, "only ДП can create users")
		return
	}

	var body dto.CreateUserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		response.HandleBindError(c, h.logger, err)
		return
	}
	created, err := h.service.CreateUser(c.Request.Context(), body)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.Created(c, created)
}

// @Tags Users
// @Summary Delete a user
// @Description Delete a user by ID
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	role := ctx.GetRole(c)
	userID := ctx.GetUserID(c)
	if role != domain.ProjectDirector && userID != id {
		response.BadRequest(c, "you can only delete your own account")
		return
	}

	if err := h.service.DeleteUser(c.Request.Context(), id); err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.NoContent(c)
}

// @Tags Users
// @Summary Update a user
// @Description Update a user by ID
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body dto.UpdateUserRequest true "User data"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	role := ctx.GetRole(c)
	userID := ctx.GetUserID(c)
	if role != domain.ProjectDirector && userID != id {
		response.BadRequest(c, "you can only update your own account")
		return
	}

	body := dto.UpdateUserRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.HandleBindError(c, h.logger, err)
		return
	}

	updated, err := h.service.UpdateUser(c.Request.Context(), id, body)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, updated)
}
