package usecase

import (
	"context"
	"errors"
	"time"

	"go-clean-arch/internal/user/domain"
)

//go:generate mockgen -source=interfaces.go -destination=mock/interfaces.go -package=mock PasswordHasher,TokenIssuer,UserRepository

var ErrInvalidCredentials = errors.New("invalid credentials")

// TokenIssuer issues tokens with an explicit, use-case-owned purpose.
type TokenIssuer interface {
	IssueAccessToken(userID uint, expireAt time.Time) (string, error)
	IssueRefreshToken(userID uint, expireAt time.Time) (string, error)
}

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
