package datamodel

import (
	"gorm.io/gorm"
)

type Claim struct {
	gorm.Model
	UserID uint   `gorm:"uniqueIndex" validate:"required"`
	Jti    string `validate:"required"`
}
