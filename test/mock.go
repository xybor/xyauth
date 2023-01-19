package test

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	gomonkey "github.com/agiledragon/gomonkey/v2"
	"github.com/xybor-x/xycond"
	"github.com/xybor/xyauth/internal/database"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/models"
	"github.com/xybor/xyauth/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type AnyData struct{}

func (a AnyData) Match(v driver.Value) bool {
	return true
}

var Any AnyData

type Mocker interface {
	Mock(sqlmock.Sqlmock, *gomonkey.Patches)
}

func StartMock(t *testing.T, m ...Mocker) func() {
	db, mock, err := sqlmock.New()
	xycond.ExpectNil(err).Test(t)

	patches := gomonkey.NewPatches()
	patches.ApplyFunc(database.GetMongoCollection, getMockCollection)

	for i := range m {
		m[i].Mock(mock, patches)
	}

	database.InitPostgresDB(db)

	f := func() {
		xycond.ExpectNil(mock.ExpectationsWereMet()).Test(t)
		db.Close()
		patches.Reset()
	}

	return f
}

type MockRegister struct {
	ExpectedEmail string
	ExpectedName  string
	ExpectedRole  string
	ExpectedError error
}

func (m MockRegister) Mock(mock sqlmock.Sqlmock, patches *gomonkey.Patches) {
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "user_credentials"`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	expectedArgs := []driver.Value{Any, Any, Any}
	if m.ExpectedEmail != "" {
		expectedArgs = append(expectedArgs, m.ExpectedEmail)
	} else {
		expectedArgs = append(expectedArgs, Any)
	}

	if m.ExpectedName != "" {
		expectedArgs = append(expectedArgs, m.ExpectedName)
	} else {
		expectedArgs = append(expectedArgs, Any)
	}

	if m.ExpectedRole != "" {
		expectedArgs = append(expectedArgs, m.ExpectedRole)
	} else {
		expectedArgs = append(expectedArgs, Any)
	}

	userQuery := mock.ExpectQuery(`INSERT INTO "users"`).WithArgs(expectedArgs...)
	if m.ExpectedError != nil {
		userQuery.WillReturnError(m.ExpectedError)
		mock.ExpectRollback()
	} else {
		userQuery.WillReturnRows(&sqlmock.Rows{})
		mock.ExpectCommit()
	}
}

type MockAuthenticate struct {
	ExpectedEmail    string
	ExpectedPassword string
	ExpectedError    error
}

func (m MockAuthenticate) Mock(mock sqlmock.Sqlmock, patches *gomonkey.Patches) {
	query := mock.ExpectQuery(`SELECT "password" FROM "user_credentials"`)

	if m.ExpectedEmail != "" {
		query.WithArgs(m.ExpectedEmail)
	}

	if m.ExpectedError != nil {
		query.WillReturnError(m.ExpectedError)
	} else if m.ExpectedPassword != "" {
		pwd, err := bcrypt.GenerateFromPassword([]byte(m.ExpectedPassword), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		query.WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow(pwd))
	} else {
		query.WillReturnRows(sqlmock.NewRows([]string{"password"}))
	}
}

var mockCollections = map[string]*mongo.Collection{}

func getMockCollection(a any) *mongo.Collection {
	c, err := utils.GetSnakeCase(a)
	if err != nil || c == "" {
		logger.Event("extract-snake-case-failed").
			Field("error", err).Field("struct", a).Panic()
	}

	if _, ok := mockCollections[c]; !ok {
		mockCollections[c] = &mongo.Collection{}
	}

	return mockCollections[c]
}

type MockInsertRefreshToken struct {
	ExpectedError error
}

func (m MockInsertRefreshToken) Mock(mock sqlmock.Sqlmock, patches *gomonkey.Patches) {
	patches.ApplyMethodFunc(
		database.GetMongoCollection(models.RefreshToken{}),
		"InsertOne",
		func(context.Context, any, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
			return &mongo.InsertOneResult{}, m.ExpectedError
		},
	)
}
