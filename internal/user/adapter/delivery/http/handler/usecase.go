package handler

import (
	"context"

	"go-clean-arch/internal/user/domain"
	ucdto "go-clean-arch/internal/user/usecase/dto"
)

//go:generate mockgen -source=usecase.go -destination=mock/usecase.go -package=mock

type UserUseCase interface {
	SignUp(ctx context.Context, params ucdto.SignUpParams) (*domain.User, error)
	Login(ctx context.Context, params ucdto.LoginParams) (*ucdto.AuthTokens, error)
	FindByID(ctx context.Context, id uint) (*domain.User, error)
}
