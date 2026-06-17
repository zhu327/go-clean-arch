package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-clean-arch/internal/user/domain"
	"go-clean-arch/internal/user/usecase/dto"
	"go-clean-arch/pkg/auth"
)

// TokenTTLs holds the TTL configuration for tokens.
type TokenTTLs struct {
	Access  time.Duration
	Refresh time.Duration
}

// UserManager handles user registration, login, and queries.
type UserManager struct {
	userRepo       UserRepository
	tokenService   auth.TokenService
	passwordHasher PasswordHasher
	ttls           TokenTTLs
}

// NewUserManager creates a new UserManager instance.
func NewUserManager(
	userRepo UserRepository,
	tokenService auth.TokenService,
	passwordHasher PasswordHasher,
	ttls TokenTTLs,
) *UserManager {
	return &UserManager{
		userRepo:       userRepo,
		tokenService:   tokenService,
		passwordHasher: passwordHasher,
		ttls:           ttls,
	}
}

// SignUp registers a new user.
func (m *UserManager) SignUp(ctx context.Context, params dto.SignUpParams) (*domain.User, error) {
	_, err := m.userRepo.FindByEmail(ctx, params.Email)
	if err == nil {
		return nil, domain.ErrUserAlreadyExists
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	hashedPassword, err := m.passwordHasher.Hash(params.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := domain.NewUser(params.Username, params.Email, hashedPassword)
	if err != nil {
		return nil, err
	}

	created, err := m.userRepo.Create(ctx, *user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &created, nil
}

// Login authenticates a user and returns tokens.
func (m *UserManager) Login(ctx context.Context, params dto.LoginParams) (*dto.AuthTokens, error) {
	user, err := m.userRepo.FindByEmail(ctx, params.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	if err := m.passwordHasher.Verify(params.Password, user.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	return m.issueTokens(user.ID)
}

// issueTokens generates an access token and a refresh token.
func (m *UserManager) issueTokens(userID uint) (*dto.AuthTokens, error) {
	req := auth.GenerateTokenRequest{UserID: userID}

	req.ExpireAt = time.Now().Add(m.ttls.Access)
	access, err := m.tokenService.GenerateToken(req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	req.ExpireAt = time.Now().Add(m.ttls.Refresh)
	refresh, err := m.tokenService.GenerateToken(req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &dto.AuthTokens{
		AccessToken:  access.TokenString,
		RefreshToken: refresh.TokenString,
	}, nil
}

// FindByID finds a user by ID.
func (m *UserManager) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	user, err := m.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, mapFindUserErr(err)
	}
	return &user, nil
}

// FindByEmail finds a user by email address.
func (m *UserManager) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := m.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, mapFindUserErr(err)
	}
	return &user, nil
}

func mapFindUserErr(err error) error {
	if errors.Is(err, domain.ErrUserNotFound) {
		return domain.ErrUserNotFound
	}
	return fmt.Errorf("failed to find user: %w", err)
}
