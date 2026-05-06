package router

import (
	"context"
	"testing"

	"go-clean-arch/internal/user/adapter/delivery/http/handler"
	userRouter "go-clean-arch/internal/user/adapter/delivery/http/router"
	"go-clean-arch/internal/user/domain"
	ucdto "go-clean-arch/internal/user/usecase/dto"

	"github.com/gin-gonic/gin"
)

type stubUserUseCase struct{}

func (s *stubUserUseCase) SignUp(_ context.Context, _ ucdto.SignUpParams) (*domain.User, error) {
	return nil, nil
}

func (s *stubUserUseCase) Login(_ context.Context, _ ucdto.LoginParams) (*ucdto.AuthTokens, error) {
	return nil, nil
}

func (s *stubUserUseCase) FindByID(_ context.Context, _ uint) (*domain.User, error) {
	return nil, nil
}

func TestRouteRegistrarInterface_Exists(t *testing.T) {
	var _ RouteRegistrar = (*stubRegistrar)(nil)
}

type stubRegistrar struct{}

func (s *stubRegistrar) RegisterRoutes(_ *gin.RouterGroup, _ gin.HandlerFunc) {}

func TestUserRegistrar_ImplementsRouteRegistrar(t *testing.T) {
	h := handler.NewUserHandler(&stubUserUseCase{})
	var _ RouteRegistrar = userRouter.NewUserRegistrar(h)
}
