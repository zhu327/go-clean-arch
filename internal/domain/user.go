package domain

import (
	"errors"
	"regexp"
	"time"
)

// 领域错误定义
var (
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrUsernameTooShort = errors.New("username is too short")
	ErrUsernameTooLong  = errors.New("username is too long")
	ErrEmptyPassword    = errors.New("password cannot be empty")
)

// 正则表达式可以作为包级别变量
var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

// User 是核心业务实体，不包含任何外部依赖
type User struct {
	ID        uint
	Username  string
	Email     string
	Password  string // 这应该是已加密的密码
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser 是创建 User 实体的工厂函数，确保了业务规则
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

// ChangePassword 是领域逻辑，例如
func (u *User) ChangePassword(newHashedPassword string) error {
	if newHashedPassword == "" {
		return ErrEmptyPassword
	}
	u.Password = newHashedPassword
	return nil
}
