package domain

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrUsernameTooShort = errors.New("username is too short")
	ErrUsernameTooLong  = errors.New("username is too long")
	ErrEmptyPassword    = errors.New("password cannot be empty")
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

// User is the core business entity.
type User struct {
	ID        uint
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser is the factory function for creating a User entity, enforcing business rules.
func NewUser(username, email, hashedPassword string) (*User, error) {
	if len(username) < 3 {
		return nil, ErrUsernameTooShort
	}
	if len(username) > 15 {
		return nil, ErrUsernameTooLong
	}
	if !emailRegex.MatchString(email) {
		return nil, ErrInvalidEmail
	}
	if hashedPassword == "" {
		return nil, ErrEmptyPassword
	}

	return &User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}, nil
}

// ChangePassword updates the user's password.
func (u *User) ChangePassword(newHashedPassword string) error {
	if newHashedPassword == "" {
		return ErrEmptyPassword
	}
	u.Password = newHashedPassword
	return nil
}
