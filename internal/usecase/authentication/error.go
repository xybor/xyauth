package authentication

import "github.com/xybor-x/xyerror"

var (
	AuthenticationError = xyerror.NewException("AuthenticationError")

	RegistrationError = AuthenticationError.NewException("RegistrationError")
	FormatError       = RegistrationError.NewException("FormatError")
)
