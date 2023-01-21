package models

import "github.com/xybor-x/xypriv"

var AllSQL = []any{&UserCredential{}, &User{}, &Client{}}
var AllNoSQL = []any{RefreshToken{}}

func init() {
	setupRelation()
	setupUserCredentialTable()
	setupWelcomePage()
}

func setupRelation() {
	xypriv.AddRelation(nil, "self", xypriv.Self)
	xypriv.AddRelation(nil, "admin", xypriv.Admin)
	xypriv.AddRelation(nil, "mod", xypriv.Moderator)
	xypriv.AddRelation(nil, "member", xypriv.LowFamiliar)
	xypriv.AddRelation(nil, "anyone", xypriv.Anyone)
}

func setupUserCredentialTable() {
	userCredentialTable := xypriv.AbstractResource("user_credential_table")
	userCredentialTable.SetPermission(xypriv.HighSecret, "create", "admin")
	userCredentialTable.SetPermission(xypriv.HighConfidential, "create", "mod")
	userCredentialTable.SetPermission(xypriv.Public, "create", "user")
}

func setupWelcomePage() {
	welcomePage := xypriv.AbstractResource("welcome_page")
	welcomePage.SetPermission(xypriv.LowPrivate, "read")
}
