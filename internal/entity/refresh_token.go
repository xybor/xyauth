package entity

import "time"

// RefreshToken represents for a refresh token stored in database.
type RefreshToken struct {
	UserID     ID
	Family     string
	Counter    uint
	Expiration time.Time
}
