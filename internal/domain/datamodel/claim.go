package datamodel

import (
	"gorm.io/gorm"
)

type Claim struct {
	*gorm.Model
	UserID uint   `gorm:"unique,index:user_id_jti"`
	Jti    string `gorm:"index:user_id_jti"`
}
