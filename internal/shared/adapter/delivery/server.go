package delivery

import (
	"go-clean-arch/internal/shared/adapter/delivery/http/middleware"
	"go-clean-arch/internal/shared/adapter/delivery/http/router"
	"go-clean-arch/pkg/auth"
	"go-clean-arch/pkg/config"
	"go-clean-arch/pkg/log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	port   string
}

func NewServer(cfg config.Config, registrars []router.RouteRegistrar, tokenService auth.TokenService) *Server {
	engine := gin.New()

	engine.Use(gin.Recovery())
	engine.Use(middleware.ErrorHandler())
	engine.Use(gin.Logger())

	router.SetupRouter(engine, registrars, tokenService)

	return &Server{engine: engine, port: cfg.Port}
}

func (s *Server) Start() error {
	addr := ":" + s.port
	log.Info("starting HTTP server", "addr", addr)
	return s.engine.Run(addr)
}
