package entities

import (
	"github.com/lib/pq"
)

type UserInfo struct {
	UserID  string         `gorm:"primaryKey;column:user_id" json:"user_id"`
	Name    string         `gorm:"column:name" json:"name"`
	AuthID  string         `gorm:"column:auth_id" json:"auth_id"`
	Roles   pq.StringArray `gorm:"type:text[];column:roles" json:"roles"`
}