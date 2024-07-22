//go:build wireinject
// +build wireinject

package di

import (
	http "go-wire/pkg/api"
	"go-wire/pkg/api/handler"
	"go-wire/pkg/config"
	"go-wire/pkg/db"
	"go-wire/pkg/repository"
	"go-wire/pkg/service/token"
	"go-wire/pkg/usecase"

	"github.com/google/wire"
)

func InitailizeApi(config config.Config) (*http.ServerHTTP, error) {

	wire.Build(
		db.ConnectDatabase,
		token.NewTokenService,
		//repository
		repository.NewUserRepository,

		//usecases
		usecase.NewUserUseCase,

		//handlers
		handler.NewUserHandler,

		http.NewServerHTTP,
	)

	return &http.ServerHTTP{}, nil

}
