package di

import (
	http "go-wire/pkg/api"
	"go-wire/pkg/config"
	"go-wire/pkg/db"

	"github.com/google/wire"
)

func InitailizeApi(config config.Config) (*http.ServerHTTP, error) {

	wire.Build(
		db.ConnectDatabase,

		http.NewServerHTTP,
	)

	return &http.ServerHTTP{}, nil

}
