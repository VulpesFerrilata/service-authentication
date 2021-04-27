package entity

import (
	gorm_custom "github.com/VulpesFerrilata/library/pkg/gorm"
	"github.com/google/uuid"
)

type Claim struct {
	gorm_custom.Model
	Jti uuid.UUID
}
