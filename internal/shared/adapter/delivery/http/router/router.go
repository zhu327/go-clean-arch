package router

import (
	userHandler "go-clean-arch/internal/user/adapter/delivery/http/handler"
	userRouter "go-clean-arch/internal/user/adapter/delivery/http/router"

	"github.com/gin-gonic/gin"
)

// RouterParams holds all handler parameters needed for route registration.
type RouterParams struct {
	UserHandler *userHandler.UserHandler
}

// SetupRouter registers all domain routes on the gin.Engine.
func SetupRouter(engine *gin.Engine, params RouterParams) {
	api := engine.Group("/api")

	userRouter.RegisterUserRoutes(api, params.UserHandler)
}
