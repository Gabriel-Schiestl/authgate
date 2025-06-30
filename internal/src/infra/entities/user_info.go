package entities

import "github.com/Gabriel-Schiestl/authgate/internal/src/utils"

type UserInfo struct {
	UserID string  `gorm:"primaryKey;type:uuid"`
	Name   string  `gorm:"not null"`
	Roles  utils.StringArray `gorm:"type:text[];default:null"`
	AuthID string            `gorm:"not null;type:uuid;constraint:OnDelete:CASCADE"`
}