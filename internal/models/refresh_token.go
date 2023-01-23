package models

import "time"

// RefreshToken represents for a refresh token stored in nosql database.
type RefreshToken struct {
	Email      string
	Family     string
	Counter    int
	Expiration time.Time
}
