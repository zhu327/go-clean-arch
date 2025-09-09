package repository

import (
	"context"
	"errors"

	"go-clean-arch/internal/domain"
	"go-clean-arch/internal/usecase/iface"

	"gorm.io/gorm"
)

// UserModel 是 GORM 的数据库模型，与 domain.User 实体分离
type UserModel struct {
	gorm.Model
	Username string `gorm:"not null;unique"`
	Email    string `gorm:"unique;not null"`
	Password string
}

// TableName 指定 GORM 使用的数据表名称
func (UserModel) TableName() string {
	return "users"
}

// toDomain 将数据库模型转换为领域实体
func (m *UserModel) toDomain() *domain.User {
	return &domain.User{
		ID:        m.ID,
		Username:  m.Username,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// fromDomain 从领域实体创建数据库模型
func fromDomain(u domain.User) *UserModel {
	return &UserModel{
		Model: gorm.Model{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
}

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
			return domain.User{}, iface.ErrUserNotFound
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
			return domain.User{}, iface.ErrUserNotFound
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
