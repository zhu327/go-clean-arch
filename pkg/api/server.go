package http

import (
	handlerInterface "go-wire/pkg/api/handler/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// gin-swagger middleware
// swagger embed files

type ServerHTTP struct {
	Engine *gin.Engine
}

// @title					Backend API
// @description				This is a sample server for Backend API.
//
// @contact.name				For API Support
// @contact.email				susiltiwari750@gmail.com
//
// @BasePath					/api
// @SecurityDefinitions.apikey	BearerAuth
// @Name						Authorization
// @In							header
// @Description				Add prefix of Bearer before  token Ex: "Bearer token"
// @Query.collection.format	multi
func NewServerHTTP(authHandler handlerInterface.AuthHandler, userHandler handlerInterface.UserHandler) *ServerHTTP {
	engine := gin.Default()

	engine.Use(gin.Logger())

	// swagger docs
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "invalid url path",
		})
	})

	return &ServerHTTP{
		Engine: engine,
	}
}

func (s *ServerHTTP) Start() error {
	return s.Engine.Run(":8080")
}
