package models

import "time"

// RefreshToken represents for a refresh token stored in nosql database.
type RefreshToken struct {
	Email      string
	Token      string
	Expiration time.Time
}
