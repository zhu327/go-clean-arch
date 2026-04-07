package repository

import (
	"context"
	"errors"

	"go-clean-arch/internal/user/domain"
	"go-clean-arch/internal/user/usecase"

	"gorm.io/gorm"
)

// UserRepository implements user data persistence.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create persists a new user record.
func (r *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	userModel := fromDomain(user)
	if err := r.db.WithContext(ctx).Create(userModel).Error; err != nil {
		return domain.User{}, err
	}
	return *userModel.toDomain(), nil
}

// FindByID queries a user by ID.
func (r *UserRepository) FindByID(ctx context.Context, id uint) (domain.User, error) {
	return r.findOne(ctx, "id = ?", id)
}

// FindByEmail queries a user by email address.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	return r.findOne(ctx, "email = ?", email)
}

func (r *UserRepository) findOne(ctx context.Context, where string, arg any) (domain.User, error) {
	var userModel UserModel
	if err := r.db.WithContext(ctx).Where(where, arg).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, usecase.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return *userModel.toDomain(), nil
}
