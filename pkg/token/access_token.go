package token

import (
	"github.com/mitchellh/mapstructure"
	"github.com/xybor-x/xypriv"
	"github.com/xybor/xyauth/internal/models"
)

type AccessToken struct {
	User models.User
	*xypriv.LeastPrivilegeToken
}

func (t *AccessToken) Unmarshal(payload any) error {
	if err := mapstructure.Decode(payload, &t); err != nil {
		return TokenError.New("can not parse to access token")
	}

	// Currently, user delegates all privileges in the normal context to the
	// token.
	t.LeastPrivilegeToken = xypriv.NewToken()
	t.LeastPrivilegeToken.AllowScope(nil)

	return nil
}
