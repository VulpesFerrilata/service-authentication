package datamodel

import (
	"gorm.io/gorm"
)

type Claim struct {
	*gorm.Model
	UserID uint   `gorm:"uniqueIndex"`
	Jti    string `gorm:"uniqueIndex"`
}
