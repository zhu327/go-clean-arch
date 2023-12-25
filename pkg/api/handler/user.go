package handler

import (
	"go-wire/pkg/domain"
	services "go-wire/pkg/usecase/interface"
	"go-wire/pkg/utils"
	"net/http"

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
//	@Tags			User
//	@Param			input	body utils.UserSignUp	true	"Input Fields"
//	@Router			/api/user/signup [post]
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
