package http

import (
	handlerInterface "go-wire/pkg/api/handler/interfaces"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// gin-swagger middleware
// swagger embed files

type ServerHTTP struct {
	Engine *gin.Engine
}

func NewServerHTTP(authHandler handlerInterface.AuthHandler, userHandler handlerInterface.UserHandler) *ServerHTTP {
	engine := gin.Default()

	engine.Use(gin.Logger())

	// swagger docs
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return &ServerHTTP{
		Engine: engine,
	}
}

func (s *ServerHTTP) Start() error {
	return s.Engine.Run(":8080")
}
