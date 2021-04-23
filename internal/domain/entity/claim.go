package entity

import (
	gorm_custom "github.com/VulpesFerrilata/library/pkg/gorm"
	"github.com/google/uuid"
)

type Claim struct {
	UserID uuid.UUID `gorm:"primaryKey"`
	Jti    uuid.UUID
	gorm_custom.Version
}
