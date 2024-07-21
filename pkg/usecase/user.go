package usecase

import (
	"context"
	"go-wire/pkg/domain"
	interfaces "go-wire/pkg/repository/interface"
	"go-wire/pkg/service/token"
	services "go-wire/pkg/usecase/interface"
	"time"
)

type UserUseCase struct {
	userRepo interfaces.UserRepository
	tokenService token.TokenService
}

func NewUserUseCase(userRepo interfaces.UserRepository, tokenService token.TokenService) services.UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
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

func (c *UserUseCase) GenerateAccessToken(ctx context.Context, tokenParams services.GenerateTokenParams) (tokenString string, err error) {
	tokenRequest := token.GenerateTokenRequest{
		UserID: tokenParams.UserID,
		ExpireAt: time.Now().Add(time.Minute * 20),
	}
	tokenResponse, err := c.tokenService.GenerateToken(tokenRequest)
	return tokenResponse.TokenString, err
}