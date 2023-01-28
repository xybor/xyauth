package service

import (
	"errors"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/xybor/xyauth/internal/database"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/models"
	"github.com/xybor/xyauth/internal/utils"
	"gorm.io/gorm"
)

func Register(email, password, role string) error {
	if role == "" {
		role = "member"
	}

	if err := utils.CheckRole(role); err != nil {
		return FormatError.New(err)
	}

	if err := utils.CheckEmail(email); err != nil {
		return FormatError.New(err)
	}

	hashedPassword, err := utils.CheckAndHashPassword(password)
	if err != nil {
		return FormatError.New(err)
	}

	username, err := generateUsername(email)
	if err != nil {
		return err
	}

	err = database.GetPostgresDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&models.User{Email: email, Username: username, Role: role}).Error; err != nil {
			return err
		}
		if err := tx.Create(&models.UserCredential{Email: email, Password: hashedPassword}).Error; err != nil {
			return err
		}
		return nil
	})

	var pgerr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgerr) && pgerr.Code == "23505" {
			return DuplicatedError.Newf("duplicated user %s", email)
		}

		logger.Event("register-failed").
			Field("email", email).Field("role", role).Field("error", err).Warning()
		return ServiceError.New("failed to register")
	}

	return nil
}

func generateUsername(email string) (string, error) {
	name, domain, found := strings.Cut(email, "@")
	if !found {
		return "", FormatError.Newf("invalid email %s", email)
	}

	if !UsernameExists(name) {
		return name, nil
	}

	username := name + "." + domain
	if !UsernameExists(username) {
		return username, nil
	}

	for i := 1; ; i++ {
		username := name + strconv.Itoa(i)
		if !UsernameExists(username) {
			return username, nil
		}
	}
}
