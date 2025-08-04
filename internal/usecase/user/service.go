package user

import (
	"context"
	"time"

	"go-clean-arch/internal/domain"
	"go-clean-arch/internal/usecase/iface"
	auth "go-clean-arch/pkg/auth"
)

type UserUseCase struct {
	userRepo     iface.UserRepository
	tokenService auth.TokenService
}

func NewUserUseCase(userRepo iface.UserRepository, tokenService auth.TokenService) iface.UserUseCase {
	return &UserUseCase{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (c *UserUseCase) SignUpUser(ctx context.Context, user domain.User) (domain.User, error) {
	users, err := c.userRepo.SignUpUser(ctx, user)
	return users, err
}

func (c *UserUseCase) FindByID(ctx context.Context, id uint) (domain.User, error) {
	users, err := c.userRepo.FindByID(ctx, id)
	return users, err
}

func (c *UserUseCase) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	users, err := c.userRepo.FindByEmail(ctx, email)
	return users, err
}

func (c *UserUseCase) GenerateAccessToken(
	ctx context.Context,
	tokenParams iface.GenerateTokenParams,
) (tokenString string,
	err error,
) {
	tokenRequest := auth.GenerateTokenRequest{
		UserID:   tokenParams.UserID,
		ExpireAt: time.Now().Add(time.Minute * 20),
	}
	tokenResponse, err := c.tokenService.GenerateToken(tokenRequest)
	if err != nil {
		return "", err
	}
	return tokenResponse.TokenString, err
}

func (c *UserUseCase) GenerateRefreshToken(
	ctx context.Context,
	tokenParams iface.GenerateTokenParams,
) (tokenString string,
	err error,
) {
	expiredAt := time.Now().Add(time.Hour * 24 * 7)
	tokenRequest := auth.GenerateTokenRequest{
		UserID:   tokenParams.UserID,
		ExpireAt: expiredAt,
	}
	tokenResponse, err := c.tokenService.GenerateToken(tokenRequest)
	if err != nil {
		return "", err
	}
	return tokenResponse.TokenString, err
}
