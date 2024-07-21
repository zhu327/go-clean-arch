package routes

import (
	"go-wire/pkg/api/handler"

	"github.com/gin-gonic/gin"
)

func UserRoutes(api *gin.RouterGroup, userHandler *handler.UserHandler) {
	// login route for user
	api.POST("/login", userHandler.UserLogin)
	// signup routes for user
	user := api.Group("/user")
	{
		user.POST("/signup", userHandler.UserSignUp)
		user.GET("/me", userHandler.UserMe)
	}

}
