package service_test

import (
	"errors"
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor/xyauth/pkg/service"
	"github.com/xybor/xyauth/test"
)

func TestAuthenticate(t *testing.T) {
	closer := test.StartMock(t, test.MockAuthenticate{ExpectedEmail: "foo@bar.com", ExpectedPassword: "123456"})
	defer closer()

	_, err := service.Authenticate("foo@bar.com", "123456")
	xycond.ExpectNil(err).Test(t)
}

func TestAuthenticateNotFound(t *testing.T) {
	closer := test.StartMock(t, test.MockAuthenticate{ExpectedEmail: "foo@bar.com"})
	defer closer()

	_, err := service.Authenticate("foo@bar.com", "123456")
	xycond.ExpectError(err, service.NotFoundError).Test(t)
}

func TestAuthenticateUnknownError(t *testing.T) {
	closer := test.StartMock(t, test.MockAuthenticate{ExpectedEmail: "foo@bar.com", ExpectedError: errors.New("err")})
	defer closer()

	_, err := service.Authenticate("foo@bar.com", "123456")
	xycond.ExpectError(err, service.ServiceError).Test(t)
}

func TestAuthenticateMismatchPassword(t *testing.T) {
	closer := test.StartMock(t, test.MockAuthenticate{ExpectedEmail: "foo@bar.com", ExpectedPassword: "123456"})
	defer closer()

	_, err := service.Authenticate("foo@bar.com", "invalid")
	xycond.ExpectError(err, service.CredentialError).Test(t)
}

func TestCreateRefreshToken(t *testing.T) {
	closer := test.StartMock(t, test.MockInsertRefreshToken{})
	defer closer()

	_, err := service.CreateRefreshToken(0)
	xycond.ExpectNil(err).Test(t)
}

func TestCreateRefreshTokenError(t *testing.T) {
	closer := test.StartMock(t, test.MockInsertRefreshToken{ExpectedError: errors.New("err")})
	defer closer()

	_, err := service.CreateRefreshToken(0)
	xycond.ExpectError(err, service.ServiceError).Test(t)
}
