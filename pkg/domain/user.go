package domain

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"not null;unique" binding:"required,min=3,max=15"`
	Password string `json:"password"`
}
