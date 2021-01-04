package model

import "github.com/VulpesFerrilata/library/pkg/model"

type Claim struct {
	model.Model
	UserID int    `gorm:"primaryKey" validate:"required"`
	Jti    string `validate:"required"`
}
