package router

import (
	"go-clean-arch/internal/user/adapter/delivery/http/handler"

	"github.com/gin-gonic/gin"
)

type UserRegistrar struct {
	handler *handler.UserHandler
}

func NewUserRegistrar(h *handler.UserHandler) *UserRegistrar {
	return &UserRegistrar{handler: h}
}

func (r *UserRegistrar) RegisterRoutes(api *gin.RouterGroup, authMW gin.HandlerFunc) {
	RegisterUserRoutes(api, r.handler, authMW)
}

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
