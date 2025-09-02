package user_validators

import (
	"chat_app_backend/internal/request_env"
	"context"
	"regexp"
)

var specialSymbolRegexp = regexp.MustCompile("[\\_\\-\\*\\&\\#\\@\\%]")
var numbersRegexp = regexp.MustCompile("\\d")
var lowercaseLettersRegexp = regexp.MustCompile("[a-z]")
var uppercaseLettersRegexp = regexp.MustCompile("[A-Z]")

type PasswordValidator struct{}

func (p PasswordValidator) Validate(password *string, _ context.Context, _ request_env.RequestEnv) bool {
	if !numbersRegexp.MatchString(*password) ||
		!lowercaseLettersRegexp.MatchString(*password) ||
		!uppercaseLettersRegexp.MatchString(*password) ||
		!specialSymbolRegexp.MatchString(*password) {
		return false
	}

	return true
}
