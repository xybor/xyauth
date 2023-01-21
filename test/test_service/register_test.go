package service_test

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/xybor-x/xycond"
	"github.com/xybor/xyauth/pkg/service"
	"github.com/xybor/xyauth/test"
)

func TestRegister(t *testing.T) {
	closer := test.StartMock(t, test.MockRegister{})
	defer closer()

	xycond.ExpectNil(service.Register("foo@bar.com", "123456", "user")).Test(t)
}

func TestRegisterEmptyRole(t *testing.T) {
	closer := test.StartMock(t, test.MockRegister{ExpectedRole: "user"})
	defer closer()

	xycond.ExpectNil(service.Register("foo@bar.com", "123456", "")).Test(t)
}

func TestRegisterInvalidRole(t *testing.T) {
	xycond.ExpectError(service.Register("foo@bar.com", "123456", "foo"), service.ValueError).Test(t)
}

func TestRegisterInvalidEmail(t *testing.T) {
	xycond.ExpectError(service.Register("foo-bar", "123456", ""), service.ValueError).Test(t)
}

func TestRegisterInvalidPassword(t *testing.T) {
	xycond.ExpectError(service.Register("foo@bar.com", "123", ""), service.ValueError).Test(t)
}

func TestRegisterDuplicated(t *testing.T) {
	closer := test.StartMock(t, test.MockRegister{ExpectedError: &pgconn.PgError{Code: "23505"}})
	defer closer()

	xycond.ExpectError(service.Register("foo@bar.com", "123456", ""), service.DuplicatedError).Test(t)
}

func TestRegisterUnknownError(t *testing.T) {
	closer := test.StartMock(t, test.MockRegister{ExpectedError: errors.New("")})
	defer closer()

	xycond.ExpectError(service.Register("foo@bar.com", "123456", ""), service.ServiceError).Test(t)
}
