package user_validators

import (
	"chat_app_backend/internal/request_env"
	"context"
	"regexp"
)

var emailRegexp = regexp.MustCompile("^\\S+@\\w+\\.\\w{2,4}$")

type EmailFormatValidator struct{}

func (e EmailFormatValidator) Validate(email *string, _ context.Context, _ request_env.RequestEnv) bool {
	if !emailRegexp.MatchString(*email) {
		return false
	}

	return true
}
