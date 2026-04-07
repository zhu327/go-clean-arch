package router

import (
	"go-clean-arch/internal/user/adapter/delivery/http/handler"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes registers all user-related routes.
// The authMiddleware is applied to routes that require authentication.
func RegisterUserRoutes(api *gin.RouterGroup, userHandler *handler.UserHandler, authMiddleware gin.HandlerFunc) {
	auth := api.Group("/auth")
	{
		auth.POST("/login", userHandler.Login)
		auth.POST("/signup", userHandler.SignUp)
	}

	user := api.Group("/user")
	user.Use(authMiddleware)
	{
		user.GET("/me", userHandler.Me)
	}
}
