package service

import (
	"errors"

	"github.com/xybor/xyauth/internal/database"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/models"
	"github.com/xybor/xyauth/internal/token"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Authenticate(email, password string) error {
	var cred = models.UserCredential{}
	result := database.GetPostgresDB().
		Select("password").
		Where("email=?", email).
		Take(&cred)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return NotFoundError.Newf("could not find the email %s", email)
		}
		return ServiceError.New("failed to authenticate")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(cred.Password), []byte(password)); err != nil {
		return CredentialError.New("failed to authenticate")
	}

	return nil
}

func CreateAccessToken(email string) (string, error) {
	var user = models.User{}
	result := database.GetPostgresDB().
		Select("email", "name", "role").
		Where("email=?", email).
		Take(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", NotFoundError.Newf("could not find the email %s", email)
		}
		logger.Event("query-user-failed").
			Field("email", email).Field("error", result.Error).Warning()
		return "", ServiceError.New("failed to authenticate")
	}

	token, err := token.Create(token.NewAccessTokenConfig(user))
	if err != nil {
		return "", err
	}

	return token, nil
}
