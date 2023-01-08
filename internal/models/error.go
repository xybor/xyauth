package models

import "github.com/xybor-x/xyerror"

var (
	DatabaseError   = xyerror.NewException("DatabaseError")
	ValueError      = DatabaseError.NewException("ValueError")
	EncryptionError = DatabaseError.NewException("EncryptionError")
	FormatError     = DatabaseError.NewException("FormatError")
	CredentialError = DatabaseError.NewException("CredentialError")
	DuplicatedError = DatabaseError.NewException("DuplicatedError")
)
