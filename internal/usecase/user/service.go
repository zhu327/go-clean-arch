package user

import (
	"context"
	"errors"
	"time"

	"go-clean-arch/internal/domain"
	"go-clean-arch/internal/usecase/iface"
	auth "go-clean-arch/pkg/auth"
	"go-clean-arch/pkg/crypto"
)

type UserUseCase struct {
	userRepo     iface.UserRepository
	tokenService auth.TokenService
}

func NewUserUseCase(userRepo iface.UserRepository, tokenService auth.TokenService) iface.UserUseCase {
	return &UserUseCase{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (uc *UserUseCase) SignUpUser(ctx context.Context, params iface.SignUpUserParams) (domain.User, error) {
	// 1. 检查用户是否已存在 (业务逻辑)
	_, err := uc.userRepo.FindByEmail(ctx, params.Email)
	if err == nil {
		// 应该返回一个明确的业务错误
		return domain.User{}, errors.New("email already exists")
	}

	// 2. 加密密码 (业务逻辑)
	hashedPassword, err := crypto.HashPassword(params.Password)
	if err != nil {
		return domain.User{}, err
	}

	// 3. 创建 Domain 实体
	user := domain.User{
		Username: params.Username,
		Email:    params.Email,
		Password: hashedPassword,
	}

	// 4. 调用 repository 进行持久化
	return uc.userRepo.SignUpUser(ctx, user)
}

func (c *UserUseCase) FindByID(ctx context.Context, id uint) (domain.User, error) {
	users, err := c.userRepo.FindByID(ctx, id)
	return users, err
}

func (c *UserUseCase) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	users, err := c.userRepo.FindByEmail(ctx, email)
	return users, err
}

func (uc *UserUseCase) Login(ctx context.Context, params iface.LoginParams) (iface.AuthTokens, error) {
	var tokens iface.AuthTokens

	// 1. 查找用户
	user, err := uc.userRepo.FindByEmail(ctx, params.Email)
	if err != nil {
		return tokens, errors.New("invalid credentials") // 不要暴露 "user not found"
	}

	// 2. 验证密码
	err = crypto.ValidatePassword(params.Password, user.Password)
	if err != nil {
		return tokens, errors.New("invalid credentials") // 返回通用错误讯息
	}

	// 3. 生成双 token (内部逻辑)
	tokenParams := auth.GenerateTokenRequest{
		UserID: user.ID,
	}

	// Access token
	tokenParams.ExpireAt = time.Now().Add(time.Minute * 20)
	accessTokenResp, err := uc.tokenService.GenerateToken(tokenParams)
	if err != nil {
		return tokens, err
	}

	// Refresh token
	tokenParams.ExpireAt = time.Now().Add(time.Hour * 24 * 7)
	refreshTokenResp, err := uc.tokenService.GenerateToken(tokenParams)
	if err != nil {
		return tokens, err
	}

	tokens.AccessToken = accessTokenResp.TokenString
	tokens.RefreshToken = refreshTokenResp.TokenString

	return tokens, nil
}
