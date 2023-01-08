package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`

	Email string `gorm:"unique;index:,unique"`
	Name  string
	Role  string

	Client []Client `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
