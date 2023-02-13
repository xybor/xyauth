package securitytoken

import "github.com/xybor/xyauth/internal/entity"

type RefreshTokenRepo interface {
	Create(id entity.ID, family string) error
	GetAndUpdate(id entity.ID, family string) (entity.RefreshToken, error)
}

type AuthorizationCodeRepo interface {
	Check(code string) (entity.AuthorizationCode, error)
}

type UserService interface {
	Get(id entity.ID) (entity.User, error)
}

type UseCase interface {
	Create(code string) (accessToken string, refreshToken string, err error)
	Refresh(refreshToken string) (string, error)
	Revoke(refreshToken string) error
	Validate(accessToken string) (entity.User, error)
}
