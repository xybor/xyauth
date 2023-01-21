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

	xycond.ExpectNil(service.Authenticate("foo@bar.com", "123456")).Test(t)
}

func TestAuthenticateNotFound(t *testing.T) {
	closer := test.StartMock(t, test.MockAuthenticate{ExpectedEmail: "foo@bar.com"})
	defer closer()

	xycond.ExpectError(service.Authenticate("foo@bar.com", "123456"), service.NotFoundError).Test(t)
}

func TestAuthenticateUnknownError(t *testing.T) {
	closer := test.StartMock(t, test.MockAuthenticate{ExpectedEmail: "foo@bar.com", ExpectedError: errors.New("")})
	defer closer()

	xycond.ExpectError(service.Authenticate("foo@bar.com", "123456"), service.ServiceError).Test(t)
}

func TestAuthenticateMismatchPassword(t *testing.T) {
	closer := test.StartMock(t, test.MockAuthenticate{ExpectedEmail: "foo@bar.com", ExpectedPassword: "123456"})
	defer closer()

	xycond.ExpectError(service.Authenticate("foo@bar.com", "invalid"), service.CredentialError).Test(t)
}

func TestCreateRefreshToken(t *testing.T) {
	closer := test.StartMock(t, test.MockInsertRefreshToken{})
	defer closer()

	_, err := service.CreateRefreshToken("foo@bar.com")
	xycond.ExpectNil(err).Test(t)
}

func TestCreateRefreshTokenError(t *testing.T) {
	closer := test.StartMock(t, test.MockInsertRefreshToken{ExpectedError: errors.New("")})
	defer closer()

	_, err := service.CreateRefreshToken("foo@bar.com")
	xycond.ExpectError(err, service.ServiceError).Test(t)
}
