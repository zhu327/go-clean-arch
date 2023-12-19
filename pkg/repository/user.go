package repository

import (
	"context"
	"errors"
	"go-wire/pkg/domain"
	interfaces "go-wire/pkg/repository/interface"

	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &userDatabase{
		DB: db,
	}
}

// find whether email is already registered
func (c *userDatabase) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	_ = c.DB.Where(("Email = ?"), email).First(&user)
	if user.ID == 0 {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}

// find user by id
func (c *userDatabase) FindByID(ctx context.Context, id uint) (domain.User, error) {
	var user domain.User
	_ = c.DB.Where(("ID = ?"), id).First(&user)
	if user.ID == 0 {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}

// create new user
func (c *userDatabase) SignUpUser(ctx context.Context, user domain.User) (domain.User, error) {
	err := c.DB.Create(&user).Error
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
