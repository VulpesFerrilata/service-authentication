package entity

import (
	gorm_custom "github.com/VulpesFerrilata/library/pkg/gorm"
)

type UserCredential struct {
	gorm_custom.Model
	Username     string
	HashPassword []byte
}
