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
	var userModel UserModel
	if err := c.DB.WithContext(ctx).Where("email = ?", email).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	// 返回转换后的 domain 实体
	return *userModel.toDomain(), nil
}

// find user by id
func (c *UserRepository) FindByID(ctx context.Context, id uint) (domain.User, error) {
	var userModel UserModel
	if err := c.DB.WithContext(ctx).Where("id = ?", id).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	// 返回转换后的 domain 实体
	return *userModel.toDomain(), nil
}

// create new user
func (c *UserRepository) SignUpUser(ctx context.Context, user domain.User) (domain.User, error) {
	userModel := fromDomain(user)
	if err := c.DB.WithContext(ctx).Create(&userModel).Error; err != nil {
		return domain.User{}, err
	}
	// 返回创建后包含 ID 和时间戳的 domain 实体
	return *userModel.toDomain(), nil
}
