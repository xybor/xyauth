package token

import (
	"github.com/mitchellh/mapstructure"
	"github.com/xybor-x/xypriv"
	"github.com/xybor/xyauth/internal/models"
	"gorm.io/gorm"
)

type AccessToken struct {
	ID        uint
	Email     string
	Username  string
	FirstName string
	LastName  string
	Role      string
	*xypriv.LeastPrivilegeToken
}

func (t *AccessToken) Unmarshal(payload any) error {
	if err := mapstructure.Decode(payload, t); err != nil {
		return TokenError.New(err)
	}

	// Currently, user delegates all privileges in the normal context to the
	// token.
	t.LeastPrivilegeToken = xypriv.NewToken()
	t.LeastPrivilegeToken.AllowScope(nil)

	return nil
}

func (t *AccessToken) GetUser() models.User {
	return models.User{
		Model: gorm.Model{
			ID: t.ID,
		},
		Email:     t.Email,
		Username:  t.Username,
		FirstName: t.FirstName,
		LastName:  t.LastName,
		Role:      t.Role,
	}
}
