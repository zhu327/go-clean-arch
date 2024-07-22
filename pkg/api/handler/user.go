package handler

import (
	"go-wire/pkg/domain"
	services "go-wire/pkg/usecase/interface"
	"go-wire/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type UserHandler struct {
	userUserCase services.UserUseCase

}

type Response struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewUserHandler(userUserCase services.UserUseCase) *UserHandler {
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
//	@Param			input	body utils.UserSignUp	true	"Input Fields"
//	@Router			/api/auth/signup [post]
//	@Success		200	{object}	utils.Response{}
//	@Failure		400	{object}	utils.Response{}
//	@Failure		409	{object}	utils.Response{}
//	@Failure		500	{object}	utils.Response{}
func (cr *UserHandler) UserSignUp(c *gin.Context) {
	var body utils.UserSignUp

	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to bind JSON", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var userBody domain.User
	copier.Copy(&userBody, &body)

	if _, err := cr.userUserCase.FindByEmail(c.Request.Context(), body.Email); err == nil {
		response := utils.ErrorResponse(400, "Error: Email already exist", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	body.Password, _ = utils.HashPassword(body.Password)
	_, err := cr.userUserCase.SignUpUser(c.Request.Context(), userBody)
	if err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to create user", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: User created", body)
	c.JSON(http.StatusOK, response)
}

// UserMe godoc
//
//	@Summary		Get User Profile
//	@Description	API for user to get their own profile
//	@Id				UserMe
//	@Tags			User
//	@Security		Bearer
//	@Router			/api/user/me [get]
//	@Success		200	{object}	Response{}
//	@Failure		400	{object}	utils.Response{}
//	@Failure		401	{object}	utils.Response{}
//	@Failure		500	{object}	utils.Response{}
func (cr *UserHandler) UserMe(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		response := utils.ErrorResponse(http.StatusBadRequest, "Couldn't get user id", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
	}
	userDetail, err := cr.userUserCase.FindByID(c.Request.Context(), uint(userID))
	if err != nil {
		response := utils.ErrorResponse(http.StatusBadRequest, "Couldn't get user detail", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := utils.SuccessResponse(http.StatusOK, "Successfully got user detail", userDetail)
	c.JSON(http.StatusOK, response)
}


// UserLogin godoc
//
//	@Summary		Login (User)
//	@Description	API for user to login
//	@Id				UserLogin
//	@Tags			Auth
//	@Param			input	body utils.LoginBody	true	"Input Fields"
//	@Router			/api/auth/login [post]
//	@Success		200	{object}	utils.Response{}
//	@Failure		400	{object}	utils.Response{}
//	@Failure		401	{object}	utils.Response{}
//	@Failure		500	{object}	utils.Response{}
func (cr *UserHandler) UserLogin(c *gin.Context) {
	var body utils.LoginBody

	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to bind JSON", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	user, err := cr.userUserCase.FindByEmail(c.Request.Context(), body.Email)
	if err != nil {
		response := utils.ErrorResponse(400, "Error: Email not found", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	err = utils.ValidatePassword(body.Password, user.Password)
	if err != nil {
		response := utils.ErrorResponse(400, "Error: Password not match", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
	}
	cr.GenerateTokens(c, user.ID)
}

// Common Function for Token generation 
func (cr *UserHandler) GenerateTokens(c *gin.Context, userID uint){

	params := services.GenerateTokenParams{
		UserID: userID,
	}
	accessToken, err := cr.userUserCase.GenerateAccessToken(c.Request.Context(), params)

	if err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to generate access token", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	refreshToken, err := cr.userUserCase.GenerateRefreshToken(c.Request.Context(), params)

	if err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to generate access token", err.Error(), nil)
		c.JSON(http.StatusBadGateway, response)
	}

	tokenResponse := utils.TokenResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}
	response := utils.SuccessResponse(http.StatusOK, "Success", tokenResponse)
	c.JSON(http.StatusOK, response)

}