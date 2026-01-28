package validation

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var V *validator.Validate

func init() {
	V = validator.New()

	// Make validator errors use json field names (email, discordUsername, etc.)
	V.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "-" {
			return ""
		}
		return strings.Split(name, ",")[0]
	})

	_ = V.RegisterValidation("no_whitespace", noWhitespace)
	_ = V.RegisterValidation("password_strong", passwordStrong)
}

func noWhitespace(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	if s == "" {
		return true // let "required" handle empty
	}
	for _, r := range s {
		if unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// Password rules:
// • At least 11 characters
// • 1 lowercase letter
// • 1 uppercase letter
// • 1 special character
func passwordStrong(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	if len([]rune(s)) < 11 {
		return false
	}

	var hasLower, hasUpper, hasSpecial bool
	for _, r := range s {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	return hasLower && hasUpper && hasSpecial
}
