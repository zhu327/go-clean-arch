package handler

import (
	"errors"
	"net/http"

	"go-clean-arch/internal/shared/adapter/delivery/http/middleware"
	"go-clean-arch/internal/user/adapter/delivery/http/dto"
	ucdto "go-clean-arch/internal/user/usecase/dto"
	"go-clean-arch/pkg/utils"

	"github.com/gin-gonic/gin"
)

// UserHandler is the HTTP handler for user operations.
type UserHandler struct {
	userUseCase UserUseCase
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(userUseCase UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

// SignUp godoc
//
//	@Summary		Signup (User)
//	@Description	API for user to register a new account
//	@Id				UserSignUp
//	@Tags			Auth
//	@Param			input	body		dto.SignUpRequest	true	"Input Fields"
//	@Router			/api/auth/signup [post]
//	@Success		201		{object}	dto.UserResponse
//	@Failure		400		{object}	middleware.ErrorResponse
//	@Failure		409		{object}	middleware.ErrorResponse
//	@Failure		413		{object}	middleware.ErrorResponse
//	@Failure		429		{object}	middleware.ErrorResponse
//	@Failure		500		{object}	middleware.ErrorResponse
func (h *UserHandler) SignUp(c *gin.Context) {
	var req dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(requestBodyError(err))
		return
	}

	user, err := h.userUseCase.SignUp(c.Request.Context(), ucdto.SignUpParams{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}

// Login godoc
//
//	@Summary		Login (User)
//	@Description	API for user to login
//	@Id				UserLogin
//	@Tags			Auth
//	@Param			input	body		dto.LoginRequest	true	"Input Fields"
//	@Router			/api/auth/login [post]
//	@Success		200		{object}	dto.TokenResponse
//	@Failure		400		{object}	middleware.ErrorResponse
//	@Failure		401		{object}	middleware.ErrorResponse
//	@Failure		413		{object}	middleware.ErrorResponse
//	@Failure		429		{object}	middleware.ErrorResponse
//	@Failure		500		{object}	middleware.ErrorResponse
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(requestBodyError(err))
		return
	}

	tokens, err := h.userUseCase.Login(c.Request.Context(), ucdto.LoginParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func requestBodyError(err error) error {
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		return utils.PayloadTooLargeError("request body too large")
	}
	return utils.BadRequestError(err.Error())
}

// Me godoc
//
//	@Summary		Get User Profile
//	@Description	API for authenticated user to get their own profile
//	@Id				UserMe
//	@Tags			User
//	@Security		Bearer
//	@Router			/api/user/me [get]
//	@Success		200		{object}	dto.UserResponse
//	@Failure		401		{object}	middleware.ErrorResponse
//	@Failure		404		{object}	middleware.ErrorResponse
//	@Failure		500		{object}	middleware.ErrorResponse
func (h *UserHandler) Me(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		_ = c.Error(utils.UnauthorizedError("authentication required"))
		return
	}

	user, err := h.userUseCase.FindByID(c.Request.Context(), userID.(uint))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}
