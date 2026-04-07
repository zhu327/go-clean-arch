package repository

import (
	"go-clean-arch/internal/user/domain"

	"gorm.io/gorm"
)

// UserModel is the GORM database model, separated from domain.User.
type UserModel struct {
	gorm.Model
	Username string `gorm:"not null;unique"`
	Email    string `gorm:"unique;not null"`
	Password string
}

// TableName specifies the GORM table name.
func (UserModel) TableName() string {
	return "users"
}

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
