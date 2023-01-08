package models

import (
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model `json:"-"`
	UserID     string
	Secret     string
}
