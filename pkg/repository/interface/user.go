package interfaces

import (
	"context"
	"go-wire/pkg/domain"
)

type UserRepository interface {
	SignUpUser(ctx context.Context, user domain.User) (domain.User, error)
	FindByID(ctx context.Context, id uint) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}
