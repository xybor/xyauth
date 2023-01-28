package token

import (
	"github.com/mitchellh/mapstructure"
	"github.com/xybor-x/xypriv"
	"github.com/xybor/xyauth/internal/models"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID       uint
	Family   string
	FamilyID uint // The family id of token
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
	return models.User{Model: gorm.Model{ID: t.ID}}
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
