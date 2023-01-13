package service

import "github.com/xybor-x/xyerror"

var (
	ServiceError = xyerror.NewException("ServiceError")

	CredentialError = ServiceError.NewException("CredentialError")
	DuplicatedError = ServiceError.NewException("DuplicatedError")
	EncryptionError = ServiceError.NewException("EncryptionError")
	NotFoundError   = ServiceError.NewException("NotFoundError")

	ValueError  = ServiceError.NewException("ValueError")
	FormatError = ValueError.NewException("FormatError")
)
