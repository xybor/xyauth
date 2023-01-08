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

func NewRefreshTokenConfig(email string) Config {
	return Config{
		payload:    RefreshToken{Email: email},
		expiration: RefreshTokenExpiration,
	}
}
