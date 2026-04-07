package handler

import (
	"net/http"
	"strconv"

	"go-clean-arch/internal/user/adapter/delivery/http/dto"
	"go-clean-arch/internal/user/usecase"
	ucdto "go-clean-arch/internal/user/usecase/dto"
	"go-clean-arch/pkg/utils"

	"github.com/gin-gonic/gin"
)

// UserHandler is the HTTP handler for user operations.
type UserHandler struct {
	userManager *usecase.UserManager
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(userManager *usecase.UserManager) *UserHandler {
	return &UserHandler{userManager: userManager}
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
//	@Failure		400		{object}	utils.AppError
//	@Failure		409		{object}	utils.AppError
//	@Failure		500		{object}	utils.AppError
func (h *UserHandler) SignUp(c *gin.Context) {
	var req dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.BadRequestError(err.Error()))
		return
	}

	user, err := h.userManager.SignUp(c.Request.Context(), ucdto.SignUpParams{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.Error(err)
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
//	@Failure		401		{object}	utils.AppError
//	@Failure		500		{object}	utils.AppError
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.BadRequestError(err.Error()))
		return
	}

	tokens, err := h.userManager.Login(c.Request.Context(), ucdto.LoginParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// Me godoc
//
//	@Summary		Get User Profile
//	@Description	API for user to get their own profile
//	@Id				UserMe
//	@Tags			User
//	@Security		Bearer
//	@Router			/api/user/me [get]
//	@Success		200		{object}	dto.UserResponse
//	@Failure		400		{object}	utils.AppError
//	@Failure		404		{object}	utils.AppError
func (h *UserHandler) Me(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.Error(utils.BadRequestError("invalid user id"))
		return
	}

	user, err := h.userManager.FindByID(c.Request.Context(), uint(userID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}
