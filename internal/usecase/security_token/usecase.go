package securitytoken

import (
	"time"

	"github.com/xybor/xyauth/internal/entity"
	"github.com/xybor/xyauth/internal/utils"
)

type usecase struct {
	codeRepo               AuthorizationCodeRepo
	userService            UserService
	refreshTokenRepo       RefreshTokenRepo
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
	engine                 Engine
}

func New(c AuthorizationCodeRepo, u UserService, r RefreshTokenRepo, e Engine) UseCase {
	return usecase{
		codeRepo:         c,
		refreshTokenRepo: r,
		userService:      u,
		engine:           e,
	}
}

func (u usecase) Create(code string) (string, string, error) {
	info, err := u.codeRepo.Check(code)
	if err != nil {
		return "", "", err
	}

	user, err := u.userService.Get(info.UserID)
	if err != nil {
		return "", "", err
	}

	atoken, err := Create(AccessToken{
		UserID:    user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}, accessTokenExpiration)
	if err != nil {
		return "", "", err
	}

	family := utils.GenerateRandomString(32)
	if err := u.refreshTokenRepo.Create(info.UserID, family); err != nil {
		return "", "", err
	}

	tToken, err := Create(RefreshToken{
		UserID:   user.ID,
		Family:   family,
		FamilyID: 0,
	}, refreshTokenExpiration)
	if err != nil {
		return "", "", err
	}

	return atoken, tToken, nil
}

func (u usecase) Refresh(refreshToken string) (string, error) {
	return "", nil
}

func (u usecase) Revoke(refreshToken string) error {
	return nil
}

func (u usecase) Validate(accessToken string) (entity.User, error) {
	user := entity.User{}

	return user, nil
}
