package models

import "github.com/xybor-x/xypriv"

var AllSQL = []any{&User{}, &UserCredential{}, &Client{}}
var AllNoSQL = []any{RefreshToken{}}

func init() {
	setupRelation()
	setupUserCredentialTable()
	setupWelcomePage()
	setupUserRole()
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
	userCredentialTable.SetPermission(xypriv.Public, "create", "member")
}

func setupWelcomePage() {
	welcomePage := xypriv.AbstractResource("welcome_page")
	welcomePage.SetPermission(xypriv.LowPrivate, "read")
}

func setupUserRole() {
	userRole := xypriv.AbstractResource("user_role")
	userRole.SetPermission(xypriv.LowPrivate, "read")
	userRole.SetPermission(xypriv.HighSecret, "update")
}
