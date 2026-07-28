package delivery

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Koshsky/erp-backend/internal/user/dto"
	"github.com/gin-gonic/gin"
)

// Response structures
type response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func writeBindError(c *gin.Context, logger *slog.Logger, err error) {
	logger.Warn("invalid request payload", "error", err)
	c.JSON(http.StatusBadRequest, response{Error: "invalid request payload"})
}

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
// @Success 200 {object} response{data=[]dto.UserResponse}
// @Failure 500 {object} response
// @Router /user [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.service.ListUsers(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list users", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: users})
}

// GetUser
// @Summary User information
// @Description Returns information about a specific user
// @Tags Users
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response{data=dto.UserResponse}
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /user/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid user id"})
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get user", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: user})
}

// @Tags Users
// @Summary Create a new user
// @Description Create a new user with the provided data
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body dto.CreateUserRequest true "User data"
// @Success 201 {object} response{data=dto.UserResponse}
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /user [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var body dto.CreateUserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		writeBindError(c, h.logger, err)
		return
	}
	created, err := h.service.CreateUser(c.Request.Context(), body)
	if err != nil {
		h.logger.Error("failed to create user", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response{Data: created})
}

// @Tags Users
// @Summary Delete a user
// @Description Delete a user by ID
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /user/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid user id"})
		return
	}

	if err := h.service.DeleteUser(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete user", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// @Tags Users
// @Summary Update a user
// @Description Update a user by ID
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body dto.UpdateUserRequest true "User data"
// @Success 200 {object} response{data=dto.UserResponse}
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /user/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid user id"})
		return
	}

	body := dto.UpdateUserRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		writeBindError(c, h.logger, err)
		return
	}

	updated, err := h.service.UpdateUser(c.Request.Context(), id, body)
	if err != nil {
		h.logger.Error("failed to update user", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: updated})
}
