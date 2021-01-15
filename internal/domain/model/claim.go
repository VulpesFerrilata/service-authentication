package model

import (
	"github.com/VulpesFerrilata/library/pkg/model"
	"github.com/google/uuid"
)

type Claim struct {
	model.Model
	UserID uuid.UUID `gorm:"type:uuid;primaryKey" validate:"required"`
	Jti    uuid.UUID `gorm:"type:uuid" validate:"required"`
}
