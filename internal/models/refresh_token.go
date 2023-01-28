package models

import "time"

// RefreshToken represents for a refresh token stored in nosql database.
type RefreshToken struct {
	ID         uint
	Family     string
	Counter    uint
	Expiration time.Time
}
