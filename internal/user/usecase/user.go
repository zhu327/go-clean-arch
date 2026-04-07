package usecase

import (
	"context"
	"errors"
	"time"

	"go-clean-arch/internal/user/domain"
	"go-clean-arch/internal/user/usecase/dto"
	"go-clean-arch/pkg/auth"
	"go-clean-arch/pkg/crypto"
	"go-clean-arch/pkg/utils"
)

// UserManager handles user registration, login, and queries.
type UserManager struct {
	userRepo     UserRepository
	tokenService auth.TokenService
}

// NewUserManager creates a new UserManager instance.
func NewUserManager(userRepo UserRepository, tokenService auth.TokenService) *UserManager {
	return &UserManager{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// SignUp registers a new user.
func (m *UserManager) SignUp(ctx context.Context, params dto.SignUpParams) (*domain.User, error) {
	_, err := m.userRepo.FindByEmail(ctx, params.Email)
	if err == nil {
		return nil, utils.ConflictError(ErrUserAlreadyExists.Error())
	}
	if !errors.Is(err, ErrUserNotFound) {
		return nil, utils.WrapError(err, "failed to check existing user")
	}

	hashedPassword, err := crypto.HashPassword(params.Password)
	if err != nil {
		return nil, utils.InternalError("failed to hash password")
	}

	user, err := domain.NewUser(params.Username, params.Email, hashedPassword)
	if err != nil {
		return nil, utils.BadRequestError(err.Error())
	}

	created, err := m.userRepo.Create(ctx, *user)
	if err != nil {
		return nil, utils.WrapError(err, "failed to create user")
	}
	return &created, nil
}

// Login authenticates a user and returns tokens.
func (m *UserManager) Login(ctx context.Context, params dto.LoginParams) (*dto.AuthTokens, error) {
	user, err := m.userRepo.FindByEmail(ctx, params.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, utils.UnauthorizedError(ErrInvalidCredentials.Error())
		}
		return nil, utils.WrapError(err, "failed to query user")
	}

	if err := crypto.ValidatePassword(params.Password, user.Password); err != nil {
		return nil, utils.UnauthorizedError(ErrInvalidCredentials.Error())
	}

	return m.issueTokens(user.ID)
}

// issueTokens generates an access token and a refresh token.
func (m *UserManager) issueTokens(userID uint) (*dto.AuthTokens, error) {
	req := auth.GenerateTokenRequest{UserID: userID}

	req.ExpireAt = time.Now().Add(20 * time.Minute)
	access, err := m.tokenService.GenerateToken(req)
	if err != nil {
		return nil, utils.InternalError("failed to generate access token")
	}

	req.ExpireAt = time.Now().Add(7 * 24 * time.Hour)
	refresh, err := m.tokenService.GenerateToken(req)
	if err != nil {
		return nil, utils.InternalError("failed to generate refresh token")
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
	if errors.Is(err, ErrUserNotFound) {
		return utils.NotFoundError(ErrUserNotFound.Error())
	}
	return utils.WrapError(err, "failed to find user")
}
