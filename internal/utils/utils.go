package utils

import (
	"reflect"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/xybor-x/xyerror"
	"github.com/xybor-x/xypriv"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/pkg/token"
)

// Check returns a privilege Checker. If the access token is available, the
// User will delegate its privilege to access token. Otherwise, the Checker will
// check with the nil subject.
func Check(ctx *gin.Context) *xypriv.Checker {
	if accessToken, ok := GetAccessToken(ctx); ok {
		return xypriv.Check(accessToken.User).Delegate(accessToken)
	}
	return xypriv.Check(nil)
}

// GetAccessToken returns the AccessToken in context. If the AccessToken does
// not exist, return (empty token, false)
func GetAccessToken(ctx *gin.Context) (token.AccessToken, bool) {
	if val, ok := ctx.Get("accessToken"); ok {
		if accessToken, ok := val.(token.AccessToken); ok {
			return accessToken, true
		}
		logger.Event("invalid-access-token").Field("token", val).Warning()
	}
	return token.AccessToken{}, false
}

// IsAuthenticated returns true if the context has access token or refresh token.
func IsAuthenticated(ctx *gin.Context) bool {
	if val, err := ctx.Cookie("access_token"); err == nil && val != "" {
		return true
	}

	if val, err := ctx.Cookie("refresh_token"); err == nil && val != "" {
		return true
	}

	return false
}

// GetSnakeCase returns the name of struct, pointer of struct, or string under
// snake case format.
func GetSnakeCase(a any) (string, error) {
	t := reflect.TypeOf(a)
	name := ""

	switch t.Kind() {
	case reflect.Pointer:
		name = t.Elem().Name()
	case reflect.Struct:
		name = t.Name()
	case reflect.String:
		name = a.(string)
	default:
		return "", xyerror.TypeError.Newf("expected input as string, struct, or pointer, but got %s", t.Name())
	}

	result := make([]rune, 0, len(name))
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}

	return string(result), nil
}
