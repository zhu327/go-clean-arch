package router

import (
	"go-clean-arch/internal/user/adapter/delivery/http/handler"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes registers all user-related routes.
func RegisterUserRoutes(api *gin.RouterGroup, userHandler *handler.UserHandler) {
	auth := api.Group("/auth")
	{
		auth.POST("/login", userHandler.Login)
		auth.POST("/signup", userHandler.SignUp)
	}

	user := api.Group("/user")
	{
		user.GET("/me", userHandler.Me)
	}
}
