package common_exceptions

import (
	"chat_app_backend/internal/exceptions"
	"net/http"
)

type UnauthorizedException struct {
	exceptions.BaseRestException
}

func (u UnauthorizedException) GetHttpStatusCode() int {
	return http.StatusUnauthorized
}

func (u UnauthorizedException) GetResponse() exceptions.Response {
	return exceptions.Response{
		Message: "Unauthorized",
	}
}
