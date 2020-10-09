package model

import (
	"gorm.io/gorm"
)

type Token struct {
	*gorm.Model
	UserID uint   `gorm:"uniqueIndex"`
	Jti    string `gorm:"uniqueIndex"`
}
