package authentication

import (
	"net/mail"

	"github.com/xybor-x/xyerror"
	"github.com/xybor/xyauth/internal/entity"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type usecase struct {
	userCredentialRepo     CredentialRepo
	userService            UserService
	auhthorizationCodeRepo AuthorizationCodeRepo
}

func New(c CredentialRepo, u UserService, a AuthorizationCodeRepo) UseCase {
	return usecase{
		userCredentialRepo:     c,
		userService:            u,
		auhthorizationCodeRepo: a,
	}
}

func (u usecase) Register(email, password, role string) error {
	if err := checkEmail(email); err != nil {
		return err
	}
	if err := checkPassword(password); err != nil {
		return err
	}
	if err := checkRole(role); err != nil {
		return err
	}

	id, err := u.userService.Create(email)
	if err != nil {
		return err
	}

	hasedPwd, err := hashPassword(password)
	if err != nil {
		return err
	}

	return u.userCredentialRepo.Create(id, hasedPwd, role)
}

func (u usecase) Authenticate(email, password string) (string, error) {
	if err := checkEmail(email); err != nil {
		return "", err
	}
	if err := checkPassword(password); err != nil {
		return "", err
	}

	credential, err := u.userCredentialRepo.Get(email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(credential.Password), []byte(password)); err != nil {
		return "", AuthenticationError.New("password is incorrect")
	}

	code := generateAuthorizationCode()
	return code, u.auhthorizationCodeRepo.Create(credential.UserID, code)
}

func checkRole(role string) error {
	for i := range entity.Roles {
		if role == entity.Roles[i] {
			return nil
		}
	}
	return xyerror.ValueError.Newf("invalid role %s", role)
}

func checkEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return xyerror.ValueError.New("invalid email")
	}
	return nil
}

func checkPassword(pwd string) error {
	if pwdlen := len(pwd); pwdlen < 6 {
		return FormatError.Newf(
			"password is required at least 6 characters, but got %d characters", pwdlen)
	}

	return nil
}

func hashPassword(pwd string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		logger.Event("invalid-password-format").
			Field("password", pwd).
			Field("error", err).Debug()
		return "", RegistrationError.New("can not hash the password")
	}
	return string(hashedPwd), nil
}

func generateAuthorizationCode() string {
	return utils.GenerateRandomString(32)
}
