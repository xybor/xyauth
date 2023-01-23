package token

import (
	"github.com/mitchellh/mapstructure"
	"github.com/xybor-x/xypriv"
	"github.com/xybor/xyauth/internal/models"
)

type RefreshToken struct {
	Email  string
	Family string
	ID     int // The id of token in its family
}

func (t *RefreshToken) Unmarshal(payload any) error {
	if err := mapstructure.Decode(payload, t); err != nil {
		return TokenError.New("can not parse to refresh token")
	}
	return nil
}

func (t RefreshToken) Context() any {
	return nil
}

func (t RefreshToken) Owner() xypriv.Subject {
	return models.User{Email: t.Email}
}

func (t RefreshToken) Permission(action ...string) xypriv.AccessLevel {
	if len(action) != 1 {
		return xypriv.NotSupport
	}

	switch action[0] {
	case "delete":
		return xypriv.TopSecret
	}

	return xypriv.NotSupport
}
