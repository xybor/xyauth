package utils

import (
	"reflect"
	"unicode"

	"github.com/xybor-x/xyerror"
)

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
