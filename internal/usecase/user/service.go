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
	// 这里需要判断错误类型，如果不是 "not found" 错误，说明可能发生了其他DB错误
	if err == nil {
		return domain.User{}, iface.ErrUserAlreadyExists
	}
	// 确保只有 "user not found" 错误才能继续
	if !errors.Is(err, iface.ErrUserNotFound) {
		return domain.User{}, err // 返回其他未知错误
	}

	// 2. 加密密码 (业务逻辑)
	hashedPassword, err := crypto.HashPassword(params.Password)
	if err != nil {
		return domain.User{}, err
	}

	// 3. 使用工厂函数创建实体，业务规则得到保证
	user, err := domain.NewUser(params.Username, params.Email, hashedPassword)
	if err != nil {
		return domain.User{}, err // 返回领域层的验证错误
	}

	// 4. 调用 repository 进行持久化
	return uc.userRepo.SignUpUser(ctx, *user)
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
		// 不管是 ErrUserNotFound 还是其他错误，都返回统一的凭证错误
		return tokens, iface.ErrInvalidCredentials
	}

	// 2. 验证密码
	err = crypto.ValidatePassword(params.Password, user.Password)
	if err != nil {
		return tokens, iface.ErrInvalidCredentials
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
