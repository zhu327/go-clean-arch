package interfaces

import (
	"context"
	"go-wire/pkg/domain"
)

type UserUseCase interface {
	FindAll(ctx context.Context) ([]domain.User, error)
	FindByID(ctx context.Context, id uint) (domain.User, error)
}
