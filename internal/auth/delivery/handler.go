package delivery

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/auth/dto"
	"github.com/Koshsky/erp-backend/internal/common/response"
)

type AuthHandler struct {
	service AuthService
	logger  *slog.Logger
}

func NewAuthHandler(logger *slog.Logger, service AuthService) *AuthHandler {
	return &AuthHandler{
		logger:  logger,
		service: service,
	}
}

// Login handles the login request.
//
//	@Tags			Auth
//	@Summary		Login
//	@Description	Authenticate user and return JWT token
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"Login credentials"
//	@Success		200		{object}	response.Response{data=dto.AuthResponse}
//	@Failure		400		{object}	response.Response{data=nil}
//	@Failure		401		{object}	response.Response{data=nil}
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	resp, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Unauthorized(c, "invalid credentials")
		return
	}

	response.OK(c, resp)
}

// Register handles the register request.
//
//	@Tags			Auth
//	@Summary		Register
//	@Description	Create user and return JWT token
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RegisterRequest	true	"Register credentials"
//	@Success		201		{object}	response.Response{data=dto.AuthResponse}
//	@Failure		400		{object}	response.Response{data=nil}
//	@Failure		500		{object}	response.Response{data=nil}
//	@Router			/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	resp, err := h.service.Register(c.Request.Context(), req.Name, req.Username, req.Password)
	if err != nil {
		response.InternalError(c, h.logger, "failed to register user", err)
		return
	}

	response.Created(c, resp)
}

// RefreshToken handles the request to refresh access token.
//
//	@Tags			Auth
//	@Summary		Refresh Token
//	@Description	Refresh access token using refresh token, returns new pair
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RefreshTokenRequest	true	"Refresh token"
//	@Success		200		{object}	response.Response{data=dto.RefreshResponse}
//	@Failure		400		{object}	response.Response{data=nil}
//	@Failure		401		{object}	response.Response{data=nil}
//	@Router			/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "refresh_token required")
		return
	}

	resp, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "invalid refresh token")
		return
	}

	response.OK(c, resp)
}
