package usecase

import (
	"context"
	"errors"

	"go-clean-arch/internal/user/domain"
)

//go:generate mockgen -source=interfaces.go -destination=mock/interfaces.go -package=mock

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// UserRepository defines the user data persistence port.
type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	FindByID(ctx context.Context, id uint) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}
