package http

import (
	"go-clean-arch/internal/adapter/delivery/http/handler"
	"go-clean-arch/internal/adapter/delivery/http/router"

	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(userHandler *handler.UserHandler) *ServerHTTP {
	engine := gin.New()

	engine.Use(gin.Logger())

	// swagger
	// engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// user routes
	router.UserRoutes(engine.Group("/api"), userHandler)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() error {
	return sh.engine.Run(":8000")
}
