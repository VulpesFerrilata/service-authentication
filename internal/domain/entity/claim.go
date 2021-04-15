package entity

import (
	"time"

	gorm_custom "github.com/VulpesFerrilata/library/pkg/gorm"
	"github.com/google/uuid"
)

type Claim struct {
	UserID    uuid.UUID `gorm:"primaryKey"`
	Jti       uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   gorm_custom.Version
}
