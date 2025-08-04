package router

import (
	"go-clean-arch/internal/adapter/delivery/http/handler"

	"github.com/gin-gonic/gin"
)

func UserRoutes(api *gin.RouterGroup, userHandler *handler.UserHandler) {
	// login route for user
	auth := api.Group("/auth")
	{
		auth.POST("/login", userHandler.UserLogin)
		auth.POST("/signup", userHandler.UserSignUp)
	}

	// signup routes for user
	user := api.Group("/user")
	{
		user.GET("/me", userHandler.UserMe)
	}
}
