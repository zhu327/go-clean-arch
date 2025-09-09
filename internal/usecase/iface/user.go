package iface

import (
	"context"

	"go-clean-arch/internal/domain"
)

// SignUpUserParams 是注册 use case 的输入参数
type SignUpUserParams struct {
	Username string
	Email    string
	Password string // 原始密码
}

// LoginParams 是登录 use case 的输入参数
type LoginParams struct {
	Email    string
	Password string
}

// AuthTokens 是登录成功后返回的结构
type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

type UserUseCase interface {
	// 直接接收参数，而不是 domain.User
	SignUpUser(ctx context.Context, params SignUpUserParams) (domain.User, error)
	// Login 负责完整的登录流程，并返回 token
	Login(ctx context.Context, params LoginParams) (AuthTokens, error)
	FindByID(ctx context.Context, id uint) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}

type GenerateTokenParams struct {
	UserID uint
}
