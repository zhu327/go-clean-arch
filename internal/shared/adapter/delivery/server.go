package delivery

import (
	"go-clean-arch/internal/shared/adapter/delivery/http/middleware"
	"go-clean-arch/internal/shared/adapter/delivery/http/router"
	userHandler "go-clean-arch/internal/user/adapter/delivery/http/handler"
	"go-clean-arch/pkg/auth"
	"go-clean-arch/pkg/config"
	"go-clean-arch/pkg/log"

	"github.com/gin-gonic/gin"
)

// Server is the HTTP server.
type Server struct {
	engine *gin.Engine
	port   string
}

// NewServer creates a new HTTP server, assembling middleware and routes.
func NewServer(cfg config.Config, userHandler *userHandler.UserHandler, tokenService auth.TokenService) *Server {
	engine := gin.New()

	engine.Use(gin.Recovery())
	engine.Use(middleware.ErrorHandler())
	engine.Use(gin.Logger())

	router.SetupRouter(engine, router.RouterParams{
		UserHandler:  userHandler,
		TokenService: tokenService,
	})

	return &Server{engine: engine, port: cfg.Port}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	addr := ":" + s.port
	log.Info("starting HTTP server", "addr", addr)
	return s.engine.Run(addr)
}
