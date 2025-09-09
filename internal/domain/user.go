package domain

import "time"

// User 是核心业务实体，不包含任何外部依赖
type User struct {
	ID        uint
	Username  string
	Email     string
	Password  string // 这应该是已加密的密码
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ChangePassword 是领域逻辑，例如
func (u *User) ChangePassword(newHashedPassword string) {
	u.Password = newHashedPassword
}
