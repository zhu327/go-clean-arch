package delivery

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go-clean-arch/internal/shared/adapter/delivery/http/middleware"
	"go-clean-arch/internal/shared/adapter/delivery/http/router"
	"go-clean-arch/pkg/auth"
	"go-clean-arch/pkg/config"
	"go-clean-arch/pkg/log"

	"github.com/gin-gonic/gin"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	maxHeaderBytes    = 1 << 20
)

type Server struct {
	httpServer  *http.Server
	cleanup     func() error
	cleanupOnce sync.Once
	cleanupErr  error
}

func NewServer(
	cfg config.Config,
	registrars []router.RouteRegistrar,
	tokenService auth.TypedTokenService,
	cleanup func() error,
) *Server {
	engine := gin.New()
	if err := engine.SetTrustedProxies(nil); err != nil {
		panic(fmt.Sprintf("disable trusted proxies: %v", err))
	}
	engine.Use(middleware.Recovery())
	engine.Use(middleware.ErrorHandler())
	engine.Use(gin.Logger())
	engine.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.SetupRouter(engine, registrars, middleware.NewAccessTokenVerifier(tokenService))
	return &Server{cleanup: cleanup, httpServer: &http.Server{
		Addr: ":" + cfg.Port, Handler: engine, ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout: readTimeout, WriteTimeout: writeTimeout, IdleTimeout: idleTimeout, MaxHeaderBytes: maxHeaderBytes,
	}}
}

func (s *Server) Start() error {
	log.Info("starting HTTP server", "addr", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	s.cleanupOnce.Do(func() {
		if s.cleanup != nil {
			s.cleanupErr = s.cleanup()
		}
	})
	return s.cleanupErr
}
