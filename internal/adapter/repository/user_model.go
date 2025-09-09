package repository

import (
	"go-clean-arch/internal/domain"

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
