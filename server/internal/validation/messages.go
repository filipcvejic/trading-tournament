package validation

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"reflect"
)

func FirstMessage(err error) string {
	var verrs validator.ValidationErrors
	ok := errors.As(err, &verrs)
	if !ok || len(verrs) == 0 {
		return "Invalid request"
	}

	e := verrs[0]
	field := e.Field() // already json name

	switch e.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return "email is invalid"
	case "min":
		if e.Kind() == reflect.Int || e.Kind() == reflect.Int64 {
			return field + " must be at least " + e.Param()
		}
		return field + " must be at least " + e.Param() + " characters"
	case "max":
		return field + " must be at most " + e.Param() + " characters"
	case "no_whitespace":
		return field + " must not contain whitespace"
	case "password_strong":
		return "password must be at least 11 characters and include 1 lowercase, 1 uppercase, and 1 special character"
	default:
		return field + " is invalid"
	}
}
