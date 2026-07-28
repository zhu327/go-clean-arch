package di

import (
	"fmt"

	sharedRouter "go-clean-arch/internal/shared/adapter/delivery/http/router"
	userRouter "go-clean-arch/internal/user/adapter/delivery/http/router"
	userUsecase "go-clean-arch/internal/user/usecase"
	"go-clean-arch/pkg/config"

	"gorm.io/gorm"
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

func provideDatabaseCleanup(db *gorm.DB) func() error {
	return func() error {
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("get database connection: %w", err)
		}
		return sqlDB.Close()
	}
}
