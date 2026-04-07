package router

import (
	"go-clean-arch/internal/shared/adapter/delivery/http/middleware"
	userHandler "go-clean-arch/internal/user/adapter/delivery/http/handler"
	userRouter "go-clean-arch/internal/user/adapter/delivery/http/router"
	"go-clean-arch/pkg/auth"

	"github.com/gin-gonic/gin"
)

// RouterParams holds all handler parameters needed for route registration.
type RouterParams struct {
	UserHandler  *userHandler.UserHandler
	TokenService auth.TokenService
}

// SetupRouter registers all domain routes on the gin.Engine.
func SetupRouter(engine *gin.Engine, params RouterParams) {
	api := engine.Group("/api")
	authMW := middleware.AuthMiddleware(params.TokenService)

	userRouter.RegisterUserRoutes(api, params.UserHandler, authMW)
}
