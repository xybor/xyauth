package securitytoken

import "github.com/xybor-x/xyerror"

var (
	SecurityTokenError = xyerror.NewException("SecurityTokenError")
	AlgorithmError     = SecurityTokenError.NewException("AlgorithmError")
)
