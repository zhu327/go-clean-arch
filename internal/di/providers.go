package di

import (
	sharedRouter "go-clean-arch/internal/shared/adapter/delivery/http/router"
	userRouter "go-clean-arch/internal/user/adapter/delivery/http/router"
	userUsecase "go-clean-arch/internal/user/usecase"
	"go-clean-arch/pkg/config"
)

func provideTokenTTLs(cfg config.Config) userUsecase.TokenTTLs {
	return userUsecase.TokenTTLs{
		Access:  cfg.AccessTokenTTL,
		Refresh: cfg.RefreshTokenTTL,
	}
}

func provideRegistrars(userRegistrar *userRouter.UserRegistrar) []sharedRouter.RouteRegistrar {
	return []sharedRouter.RouteRegistrar{userRegistrar}
}
