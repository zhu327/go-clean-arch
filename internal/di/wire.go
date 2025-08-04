//go:build wireinject
// +build wireinject

package di

import (
	"go-clean-arch/internal/adapter/delivery/http"
	"go-clean-arch/internal/adapter/delivery/http/handler"
	"go-clean-arch/internal/adapter/repository"
	"go-clean-arch/internal/usecase/user"
	"go-clean-arch/pkg/auth"
	"go-clean-arch/pkg/config"
	"go-clean-arch/pkg/db"

	"github.com/google/wire"
)

func InitailizeApi(config config.Config) (*http.ServerHTTP, error) {
	wire.Build(
		db.ConnectDatabase,
		auth.NewTokenService,
		// repository
		repository.NewUserRepository,

		// usecases
		user.NewUserUseCase,

		// handlers
		handler.NewUserHandler,

		// http server
		http.NewServerHTTP,
	)

	return &http.ServerHTTP{}, nil
}
