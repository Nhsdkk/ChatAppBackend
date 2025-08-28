package user_validators

import (
	"context"
	"regexp"
)

var specialSymbolRegexp = regexp.MustCompile("[\\_\\-\\*\\&\\#\\@\\%]")
var numbersRegexp = regexp.MustCompile("\\d")
var lowercaseLettersRegexp = regexp.MustCompile("[a-z]")
var uppercaseLettersRegexp = regexp.MustCompile("[A-Z]")

type PasswordValidator struct{}

func (p PasswordValidator) Validate(password *string, _ context.Context) bool {
	if !numbersRegexp.MatchString(*password) ||
		!lowercaseLettersRegexp.MatchString(*password) ||
		!uppercaseLettersRegexp.MatchString(*password) ||
		!specialSymbolRegexp.MatchString(*password) {
		return false
	}

	return true
}
