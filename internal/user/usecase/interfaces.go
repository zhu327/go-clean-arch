package usecase

import (
	"context"
	"errors"

	"go-clean-arch/internal/user/domain"
)

//go:generate mockgen -source=interfaces.go -destination=mock/interfaces.go -package=mock PasswordHasher,UserRepository

var ErrInvalidCredentials = errors.New("invalid credentials")

// PasswordHasher hashes and verifies passwords.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) error
}

// UserRepository defines the user data persistence port.
type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	FindByID(ctx context.Context, id uint) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}
