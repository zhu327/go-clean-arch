package domain

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"not null;unique" binding:"required,min=3,max=15"`
	Email    string `json:"email" gorm:"unique;not null" binding:"required,email"`
	Password string `json:"password"`
}
