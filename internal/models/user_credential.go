package models

import (
	"time"

	"gorm.io/gorm"
)

// UserCredential represents for credential info of a User.
type UserCredential struct {
	Email     string         `gorm:"primaryKey;not null;index:,unique"`
	Password  string         `json:"-"`
	User      User           `gorm:"foreignKey:Email;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
