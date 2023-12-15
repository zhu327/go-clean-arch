package handler

import (
	services "go-wire/pkg/usecase/interface"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

// FindAll godoc
// @summary Get all users
// @description Get all users
// @tags users
// @security ApiKeyAuth
// @id FindAll
// @produce json
// @Router /api/users [get]
// @response 200 {object} []Response "OK"
func (cr *UserHandler) FindAll(c *gin.Context) {
	users, err := cr.userUserCase.FindAll(c.Request.Context())

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, users)
	}
}

// FindByID godoc
// @summary Get user by id
// @description Get user by id
// @tags users
// @security ApiKeyAuth
// @id FindByID
// @produce json
// @param id path int true "User ID"
// @Router /api/users/{id} [get]
// @response 200 {object} Response "OK"
func (cr *UserHandler) FindByID(c *gin.Context) {
	paramid := c.Param("id")
	id, err := strconv.Atoi(paramid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot parse id",
		})
		return
	}

	user, err := cr.userUserCase.FindByID(c.Request.Context(), uint(id))

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, user)
	}
}
