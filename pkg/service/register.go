package service

import (
	"errors"
	"net/mail"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/xybor/xyauth/internal/database"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(email, password, role string) error {
	if role == "" {
		role = "user"
	}

	if err := checkRole(role); err != nil {
		return err
	}

	if err := checkEmail(email); err != nil {
		return err
	}

	hashedPassword, err := checkAndHashPassword(password)
	if err != nil {
		return err
	}

	err = database.GetPostgresDB().Create(&models.UserCredential{
		Email: email,
		User: models.User{
			Email: email,
			Role:  role,
		},
		Password: hashedPassword,
	}).Error

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

var roles = []string{"admin", "mod", "user"}

func checkRole(role string) error {
	for i := range roles {
		if role == roles[i] {
			return nil
		}
	}
	return ValueError.Newf("invalid role %s", role)
}

func checkEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return FormatError.New("invalid email")
	}

	return nil
}

func checkAndHashPassword(pwd string) (string, error) {
	if pwdlen := len(pwd); pwdlen < 6 {
		return "", FormatError.Newf(
			"password is required at least 6 characters, but got %d characters", pwdlen)
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		logger.Event("invalid-password-format").
			Field("password", pwd).
			Field("error", err).Debug()
		return "", EncryptionError.New("password is invalid format")
	}
	return string(hashedPwd), nil
}
