package service_test

// func TestRegister(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	mock.ExpectBegin()
// 	mock.ExpectExec(`INSERT INTO "user_credentials"`).WillReturnResult(sqlmock.NewResult(1, 1))
// 	mock.ExpectQuery(`INSERT INTO "users"`).WillReturnRows(&sqlmock.Rows{})
// 	mock.ExpectCommit()

// 	database.InitPostgresDB(db)

// 	xycond.ExpectNil(service.Register("foo@bar.com", "123456", "user")).Test(t)
// 	xycond.ExpectNil(mock.ExpectationsWereMet()).Test(t)
// }
