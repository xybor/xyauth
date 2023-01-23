package token

import "github.com/xybor-x/xyerror"

var (
	TokenError     = xyerror.NewException("TokenError")
	CertKeyError   = TokenError.NewException("CertKeyError")
	AlgorithmError = TokenError.NewException("AlgorithmError")
)
