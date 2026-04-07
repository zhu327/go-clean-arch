package delivery

import (
	"go-clean-arch/internal/shared/adapter/delivery/http/middleware"
	"go-clean-arch/internal/shared/adapter/delivery/http/router"
	userHandler "go-clean-arch/internal/user/adapter/delivery/http/handler"
	"go-clean-arch/pkg/log"

	"github.com/gin-gonic/gin"
)

// Server is the HTTP server.
type Server struct {
	engine *gin.Engine
}

// NewServer creates a new HTTP server, assembling middleware and routes.
func NewServer(userHandler *userHandler.UserHandler) *Server {
	engine := gin.New()

	engine.Use(gin.Recovery())
	engine.Use(middleware.ErrorHandler())
	engine.Use(gin.Logger())

	router.SetupRouter(engine, router.RouterParams{
		UserHandler: userHandler,
	})

	return &Server{engine: engine}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	log.Info("starting HTTP server on :8000")
	return s.engine.Run(":8000")
}
