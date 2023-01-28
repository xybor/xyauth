package token

import (
	"time"
)

type Config struct {
	payload    any
	expiration time.Duration
}

func NewAccessTokenConfig(a AccessToken) Config {
	return Config{
		payload:    a,
		expiration: AccessTokenExpiration,
	}
}

func NewRefreshTokenConfig(r RefreshToken) Config {
	return Config{
		payload:    r,
		expiration: RefreshTokenExpiration,
	}
}
