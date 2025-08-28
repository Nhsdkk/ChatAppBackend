package user_validators

import (
	"context"
	"regexp"
)

var emailRegexp = regexp.MustCompile("^\\S+@\\w+\\.\\w{2,4}$")

type EmailFormatValidator struct{}

func (e EmailFormatValidator) Validate(email *string, _ context.Context) bool {
	if !emailRegexp.MatchString(*email) {
		return false
	}

	return true
}
