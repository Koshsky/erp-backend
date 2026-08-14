package delivery

import (
	"log/slog"

	authservice "github.com/Koshsky/erp-backend/internal/auth/service"
	"github.com/Koshsky/erp-backend/pkg/errors"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/auth/dto"
	"github.com/Koshsky/erp-backend/internal/response"
)

type AuthHandler struct {
	service AuthService
	logger  *slog.Logger
}

// NewAuthHandler builds the auth handler.
func NewAuthHandler(logger *slog.Logger, svc *authservice.AuthService) *AuthHandler {
	return &AuthHandler{
		logger:  logger,
		service: svc,
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
//	@Success		200		{object}	response.SuccessResponse{data=dto.AuthResponse,error=nil}
//	@Failure		400		{object}	response.ErrorResponse{data=nil}
//	@Failure		401		{object}	response.ErrorResponse{data=nil}
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid request")
		return
	}

	resp, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Unauthorized(c, errors.CodeInvalidCredentials, "invalid credentials")
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
//	@Success		201		{object}	response.SuccessResponse{data=dto.AuthResponse,error=nil}
//	@Failure		400		{object}	response.ErrorResponse{data=nil}
//	@Failure		500		{object}	response.ErrorResponse{data=nil}
//	@Router			/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid request")
		return
	}

	resp, err := h.service.Register(
		c.Request.Context(),
		req.LastName,
		req.FirstName,
		req.MiddleName,
		req.Username,
		req.Password,
	)
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
//	@Success		200		{object}	response.SuccessResponse{data=dto.RefreshResponse,error=nil}
//	@Failure		400		{object}	response.ErrorResponse{data=nil}
//	@Failure		401		{object}	response.ErrorResponse{data=nil}
//	@Router			/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "refresh_token required")
		return
	}

	resp, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, errors.CodeInvalidToken, "invalid refresh token")
		return
	}

	response.OK(c, resp)
}
