package repository

import (
	"context"
	"errors"

	"go-clean-arch/internal/domain"
	"go-clean-arch/internal/usecase/iface"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) iface.UserRepository {
	return &UserRepository{
		DB: db,
	}
}

// find whether email is already registered
func (c *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	_ = c.DB.Where(("Email = ?"), email).First(&user)
	if user.ID == 0 {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}

// find user by id
func (c *UserRepository) FindByID(ctx context.Context, id uint) (domain.User, error) {
	var user domain.User
	_ = c.DB.Where(("ID = ?"), id).First(&user)
	if user.ID == 0 {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}

// create new user
func (c *UserRepository) SignUpUser(ctx context.Context, user domain.User) (domain.User, error) {
	err := c.DB.Create(&user).Error
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
