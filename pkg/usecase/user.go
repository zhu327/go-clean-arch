package usecase

import (
	"context"
	"go-wire/pkg/domain"
	interfaces "go-wire/pkg/repository/interface"
	services "go-wire/pkg/usecase/interface"
)

type UserUseCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(userRepo interfaces.UserRepository) services.UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
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
