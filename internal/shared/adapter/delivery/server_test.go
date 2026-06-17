package delivery

import (
	"testing"

	"go-clean-arch/internal/shared/adapter/delivery/http/router"
	"go-clean-arch/pkg/auth"
	"go-clean-arch/pkg/config"

	"github.com/gin-gonic/gin"
)

type stubRegistrar struct {
	called bool
}

func (s *stubRegistrar) RegisterRoutes(_ *gin.RouterGroup, _ gin.HandlerFunc) {
	s.called = true
}

type stubTokenService struct{}

func (s *stubTokenService) GenerateToken(_ auth.GenerateTokenRequest) (auth.GenerateTokenResponse, error) {
	return auth.GenerateTokenResponse{}, nil
}

func (s *stubTokenService) ValidateToken(_ string) (*auth.TokenClaims, error) {
	return &auth.TokenClaims{}, nil
}

func TestNewServer_AcceptsRouteRegistrars(t *testing.T) {
	reg := &stubRegistrar{}
	cfg := config.Config{Port: "0"}
	svc := &stubTokenService{}

	srv := NewServer(cfg, []router.RouteRegistrar{reg}, svc)
	if srv == nil {
		t.Fatal("expected non-nil server")
	}
	if !reg.called {
		t.Fatal("expected RegisterRoutes to be called")
	}
}
