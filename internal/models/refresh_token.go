package models

import "time"

type RefreshToken struct {
	Email      string
	Token      string
	Expiration time.Time
}
