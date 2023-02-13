package entity

import (
	"gorm.io/gorm"
)

type ID uint

var Roles = [...]string{}

type User struct {
	gorm.Model

	Email       string `gorm:"unique;not null"`
	Username    string `gorm:"unique;not null;index,unique"`
	FirstName   string
	LastName    string
	PhoneNumber string
	Address     string
	Role        string

	Credential        Credential        `gorm:"foreignKey:UserID;references:ID"`
	AuthorizationCode AuthorizationCode `gorm:"foreignKey:UserID;references:ID"`
	RefreshToken      RefreshToken      `gorm:"foreignKey:UserID;references:ID"`
}
