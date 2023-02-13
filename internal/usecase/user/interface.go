package user

import "github.com/xybor/xyauth/internal/entity"

type UserRepo interface {
	GetByUsername(username string) (entity.User, error)
	UpdateByUsername(username string, user entity.User) error
}

type SecurityTokenService interface {
	Validate(accessToken string) (entity.User, error)
}

type UseCase interface {
	GetByUsername(username string) (entity.User, error)
	UpdateByUsername(username string, user entity.User) error
}
