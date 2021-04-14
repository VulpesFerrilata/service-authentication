package entity

import (
	gorm_custom "github.com/VulpesFerrilata/library/pkg/gorm"
	"github.com/google/uuid"
)

type Claim struct {
	gorm_custom.Model
	UserID uuid.UUID `gorm:"type:uuid;primaryKey" validate:"required"`
	Jti    uuid.UUID `gorm:"type:uuid" validate:"required"`
}
