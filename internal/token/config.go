package token

import (
	"time"

	"github.com/xybor/xyauth/internal/models"
)

type Config struct {
	payload    any
	expiration time.Duration
}

func NewAccessTokenConfig(u models.User) Config {
	return Config{
		payload:    AccessToken{User: u},
		expiration: AccessTokenExpiration,
	}
}

func NewRefreshTokenConfig(email, family string, familyID int) Config {
	return Config{
		payload:    RefreshToken{Email: email, Family: family, ID: familyID},
		expiration: RefreshTokenExpiration,
	}
}
