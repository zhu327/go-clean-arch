package iface

import "errors"

// 定义 Repository 层可能返回的错误
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// 定义 UseCase 层可能返回的错误
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)
