package service

import (
	"context"
	"errors"
	"time"

	"github.com/xybor/xyauth/internal/database"
	"github.com/xybor/xyauth/internal/models"
	"github.com/xybor/xyauth/internal/utils"
	"github.com/xybor/xyauth/pkg/token"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Authenticate(email, password string) error {
	var cred = models.UserCredential{}
	result := database.GetPostgresDB().Where("email=?", email).Take(&cred)

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
	result := database.GetPostgresDB().Where("email=?", email).Take(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", NotFoundError.Newf("could not find the email %s", email)
		}
		utils.GetLogger().Event("query-user-failed").
			Field("email", email).Field("error", result.Error).Warning()
		return "", ServiceError.New("failed to authenticate")
	}

	token, err := token.Create(token.NewAccessTokenConfig(user))
	if err != nil {
		return "", err
	}

	return token, nil
}

func CreateRefreshToken(email string) (string, error) {
	value, err := token.Create(token.NewRefreshTokenConfig(email))
	if err != nil {
		return "", err
	}

	_, err = database.GetMongoCollection(models.RefreshToken{}).InsertOne(
		context.Background(), models.RefreshToken{
			Email:      email,
			Token:      value,
			Expiration: time.Now().Add(token.RefreshTokenExpiration),
		},
	)

	if err != nil {
		utils.GetLogger().Event("whitelist-refresh-token-failed").
			Field("token", value).Field("error", err).Error()
		return "", ServiceError.New("can not insert refresh token")
	}

	return value, nil
}
