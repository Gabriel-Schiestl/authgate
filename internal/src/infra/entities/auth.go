package entities

import (
	"time"
)

type Auth struct {
	ID                 string                `gorm:"primaryKey;type:uuid"`
	IdentifierType     IdentifierType `gorm:"type:int;not null"`
	IdentifierValue    string                `gorm:"not null"`
	Password           string                `gorm:"not null"`
	UserInfo           UserInfo       `gorm:"foreignKey:AuthID;references:ID"`
	EncryptToken       bool                  `gorm:"default:false"`
	LastLoginAt        *time.Time            `gorm:"default:null"`
	WrongAttempts      int                   `gorm:"not null"`
	MaxWrongAttempts   int                  `gorm:"not null"`
	RecoveryToken      *string                `gorm:"default:null"`
	MaxTokenAgeSeconds int                  `gorm:"not null"`
	CreatedAt          *time.Time            `gorm:"autoCreateTime"`
	UpdatedAt          *time.Time            `gorm:"autoUpdateTime"`
}