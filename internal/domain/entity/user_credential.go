package entity

import (
	"time"

	gorm_custom "github.com/VulpesFerrilata/library/pkg/gorm"
	"github.com/google/uuid"
)

type UserCredential struct {
	gorm_custom.Model
	HashPassword []byte
	UserID       uuid.UUID
	CreatedAt    time.Time
}
