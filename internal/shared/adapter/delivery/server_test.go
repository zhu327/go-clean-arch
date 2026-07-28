package delivery

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func (s *stubTokenService) ValidateTokenOfType(_ string, _ auth.TokenType) (*auth.TokenClaims, error) {
	return &auth.TokenClaims{}, nil
}
func (s *stubTokenService) IssueAccessToken(uint, time.Time) (string, error)  { return "", nil }
func (s *stubTokenService) IssueRefreshToken(uint, time.Time) (string, error) { return "", nil }

func TestNewServerConfiguresHTTPServerAndHealthz(t *testing.T) {
	reg := &stubRegistrar{}
	srv := NewServer(config.Config{Port: "0"}, []router.RouteRegistrar{reg}, &stubTokenService{}, nil)
	if srv == nil || srv.httpServer.ReadHeaderTimeout == 0 || srv.httpServer.MaxHeaderBytes == 0 {
		t.Fatal("expected configured explicit HTTP server")
	}
	if !reg.called {
		t.Fatal("expected RegisterRoutes to be called")
	}

	w := httptest.NewRecorder()
	srv.httpServer.Handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("healthz status = %d, want 200", w.Code)
	}
}

func TestServerShutdown(t *testing.T) {
	cleaned := 0
	srv := NewServer(config.Config{Port: "0"}, nil, &stubTokenService{}, func() error { cleaned++; return nil })
	if err := srv.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown returned error: %v", err)
	}
	if err := srv.Shutdown(context.Background()); err != nil {
		t.Fatalf("second Shutdown returned error: %v", err)
	}
	if cleaned != 1 {
		t.Fatalf("cleanup calls = %d, want 1", cleaned)
	}
}
