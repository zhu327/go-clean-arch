package handler

import (
	"context"
	"testing"

	"go-clean-arch/internal/user/domain"
	"go-clean-arch/internal/user/usecase"
	ucdto "go-clean-arch/internal/user/usecase/dto"
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

func TestNewUserHandler_AcceptsUserUseCaseInterface(t *testing.T) {
	var uc UserUseCase = &stubUserUseCase{}
	h := NewUserHandler(uc)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestUserManager_SatisfiesUserUseCaseInterface(t *testing.T) {
	var _ UserUseCase = (*usecase.UserManager)(nil)
}
