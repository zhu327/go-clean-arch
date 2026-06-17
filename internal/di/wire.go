//go:build wireinject
// +build wireinject

package di

import (
	delivery "go-clean-arch/internal/shared/adapter/delivery"
	userHandler "go-clean-arch/internal/user/adapter/delivery/http/handler"
	userRouter "go-clean-arch/internal/user/adapter/delivery/http/router"
	userRepo "go-clean-arch/internal/user/adapter/repository"
	userUsecase "go-clean-arch/internal/user/usecase"
	"go-clean-arch/pkg/auth"
	"go-clean-arch/pkg/config"
	"go-clean-arch/pkg/crypto"
	"go-clean-arch/pkg/db"

	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*delivery.Server, error) {
	wire.Build(
		db.ConnectDatabase,
		auth.NewTokenService,
		crypto.NewBcryptHasher,
		wire.Bind(new(userUsecase.PasswordHasher), new(*crypto.BcryptHasher)),
		userRepo.NewUserRepository,
		wire.Bind(new(userUsecase.UserRepository), new(*userRepo.UserRepository)),
		provideTokenTTLs,
		userUsecase.NewUserManager,
		wire.Bind(new(userHandler.UserUseCase), new(*userUsecase.UserManager)),
		userHandler.NewUserHandler,
		userRouter.NewUserRegistrar,
		provideRegistrars,
		delivery.NewServer,
	)
	return &delivery.Server{}, nil
}
