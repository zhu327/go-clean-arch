//go:build wireinject
// +build wireinject

package di

import (
	delivery "go-clean-arch/internal/shared/adapter/delivery"
	userHandler "go-clean-arch/internal/user/adapter/delivery/http/handler"
	userRepo "go-clean-arch/internal/user/adapter/repository"
	userUsecase "go-clean-arch/internal/user/usecase"
	"go-clean-arch/pkg/auth"
	"go-clean-arch/pkg/config"
	"go-clean-arch/pkg/db"

	"github.com/google/wire"
)

// InitializeAPI builds and injects all dependencies, returning the HTTP server.
func InitializeAPI(cfg config.Config) (*delivery.Server, error) {
	wire.Build(
		db.ConnectDatabase,
		auth.NewTokenService,
		userRepo.NewUserRepository,
		wire.Bind(new(userUsecase.UserRepository), new(*userRepo.UserRepository)),
		userUsecase.NewUserManager,
		userHandler.NewUserHandler,
		delivery.NewServer,
	)
	return &delivery.Server{}, nil
}
