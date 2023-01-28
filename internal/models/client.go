package models

import (
	"gorm.io/gorm"
)

// Client represents for an OAuth client.
type Client struct {
	gorm.Model
	UserID string
	Secret string
}
