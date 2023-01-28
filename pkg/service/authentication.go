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

func Authenticate(email, password string) (ID uint, err error) {
	record := struct {
		ID       uint
		Password string
	}{}
	result := database.GetPostgresDB().
		Model(&models.UserCredential{}).
		Select("users.id", "user_credentials.password").
		Where("user_credentials.email=?", email).
		Joins("join users on user_credentials.email=users.email").
		Take(&record)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, NotFoundError.Newf("could not find the email %s", email)
		}
		logger.Event("authenticate-failed").Field("error", result.Error).Warning()
		return 0, ServiceError.New("failed to authenticate")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(record.Password), []byte(password)); err != nil {
		return 0, CredentialError.New("failed to authenticate")
	}

	return record.ID, nil
}

func CreateAccessToken(id uint) (string, error) {
	user := models.User{}

	result := database.GetPostgresDB().
		Select("id", "email", "username", "first_name", "last_name", "role").
		Where("id=?", id).
		Take(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", NotFoundError.Newf("could not find the email %d", id)
		}
		logger.Event("query-user-failed").
			Field("user_id", id).Field("error", result.Error).Warning()
		return "", ServiceError.New("failed to authenticate")
	}

	token, err := token.Create(token.NewAccessTokenConfig(token.AccessToken{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}))
	if err != nil {
		return "", err
	}

	return token, nil
}
