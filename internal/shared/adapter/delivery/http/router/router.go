package router

import (
	"go-clean-arch/internal/shared/adapter/delivery/http/middleware"
	"go-clean-arch/pkg/auth"

	"github.com/gin-gonic/gin"
)

type RouteRegistrar interface {
	RegisterRoutes(api *gin.RouterGroup, authMW gin.HandlerFunc)
}

func SetupRouter(engine *gin.Engine, registrars []RouteRegistrar, tokenService auth.TokenService) {
	api := engine.Group("/api")
	authMW := middleware.AuthMiddleware(tokenService)

	for _, r := range registrars {
		r.RegisterRoutes(api, authMW)
	}
}
