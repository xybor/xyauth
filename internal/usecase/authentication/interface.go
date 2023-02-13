package authentication

import "github.com/xybor/xyauth/internal/entity"

type CredentialRepo interface {
	Create(id entity.ID, password, role string) error
	Get(email string) (entity.Credential, error)
}

type AuthorizationCodeRepo interface {
	Create(id entity.ID, code string) error
}

type UserService interface {
	Create(email string) (entity.ID, error)
}

type UseCase interface {
	Register(email, password, role string) error
	Authenticate(email, password string) (code string, err error)
}
