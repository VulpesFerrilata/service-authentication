package datamodel

import (
	"time"

	"gorm.io/gorm"
)

type Claim struct {
	UserID    uint   `gorm:"primarykey" validate:"required"`
	Jti       string `validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
