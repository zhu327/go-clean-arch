package handler

import (
	"net/http"
	"strconv"

	"go-clean-arch/internal/adapter/delivery/http/dto"
	"go-clean-arch/internal/usecase/iface"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUserCase iface.UserUseCase
}

type Response struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewUserHandler(userUserCase iface.UserUseCase) *UserHandler {
	return &UserHandler{
		userUserCase: userUserCase,
	}
}

// UserSignUp godoc
//
//	@Summary		Signup (User)
//	@Description	API for user to register a new account
//	@Id				UserSignUp
//	@Tags			Auth
//	@Param			input	body		dto.UserSignUp	true	"Input Fields"
//	@Router			/api/auth/signup [post]
//	@Success		200		{object}	dto.Response{}
//	@Failure		400		{object}	dto.Response{}
//	@Failure		409		{object}	dto.Response{}
//	@Failure		500		{object}	dto.Response{}
func (h *UserHandler) UserSignUp(c *gin.Context) {
	var body dto.UserSignUp
	if err := c.BindJSON(&body); err != nil {
		response := dto.ErrorResponse(http.StatusBadRequest, "Error: Failed to bind JSON", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 将 handler DTO 转换为 use case 参数
	params := iface.SignUpUserParams{
		Username: body.Username,
		Email:    body.Email,
		Password: body.Password,
	}

	// 只调用一个方法
	user, err := h.userUserCase.SignUpUser(c.Request.Context(), params)
	if err != nil {
		// 根据 use case 返回的错误类型决定 HTTP 状态码
		if err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, dto.ErrorResponse(http.StatusConflict, err.Error(), err.Error(), nil))
			return
		}
		c.JSON(
			http.StatusBadRequest,
			dto.ErrorResponse(http.StatusBadRequest, "Failed to create user", err.Error(), nil),
		)
		return
	}

	// 成功响应，注意不要返回密码
	userResponse := Response{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	c.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "User created", userResponse))
}

// UserMe godoc
//
//	@Summary		Get User Profile
//	@Description	API for user to get their own profile
//	@Id				UserMe
//	@Tags			User
//	@Security		Bearer
//	@Router			/api/user/me [get]
//	@Success		200		{object}	Response{}
//	@Failure		400		{object}	dto.Response{}
//	@Failure		401		{object}	dto.Response{}
//	@Failure		500		{object}	dto.Response{}
func (cr *UserHandler) UserMe(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		response := dto.ErrorResponse(http.StatusBadRequest, "Couldn't get user id", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
	}
	userDetail, err := cr.userUserCase.FindByID(c.Request.Context(), uint(userID))
	if err != nil {
		response := dto.ErrorResponse(http.StatusBadRequest, "Couldn't get user detail", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := dto.SuccessResponse(http.StatusOK, "Successfully got user detail", userDetail)
	c.JSON(http.StatusOK, response)
}

// UserLogin godoc
//
//	@Summary		Login (User)
//	@Description	API for user to login
//	@Id				UserLogin
//	@Tags			Auth
//	@Param			input	body		dto.LoginBody	true	"Input Fields"
//	@Router			/api/auth/login [post]
//	@Success		200		{object}	dto.Response{}
//	@Failure		400		{object}	dto.Response{}
//	@Failure		401		{object}	dto.Response{}
//	@Failure		500		{object}	dto.Response{}
func (h *UserHandler) UserLogin(c *gin.Context) {
	var body dto.LoginBody
	if err := c.BindJSON(&body); err != nil {
		response := dto.ErrorResponse(http.StatusBadRequest, "Error: Failed to bind JSON", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	params := iface.LoginParams{
		Email:    body.Email,
		Password: body.Password,
	}

	tokens, err := h.userUserCase.Login(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse(http.StatusUnauthorized, "Login failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Login successful", tokens))
}
