package router

import (
	"github.com/gin-gonic/gin"
	"go-clean-arch/internal/shared/adapter/delivery/http/middleware"
)

type RouteRegistrar interface {
	RegisterRoutes(api *gin.RouterGroup, authMW gin.HandlerFunc)
}

func SetupRouter(engine *gin.Engine, registrars []RouteRegistrar, verifier middleware.AccessTokenVerifier) {
	api := engine.Group("/api")
	authMW := middleware.AuthMiddleware(verifier)

	for _, r := range registrars {
		r.RegisterRoutes(api, authMW)
	}
}
