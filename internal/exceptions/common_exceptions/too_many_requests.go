package common_exceptions

import (
	"chat_app_backend/internal/exceptions"
	"net/http"
)

type TooManyRequestsException struct {
	exceptions.BaseRestException
}

func (TooManyRequestsException) GetResponse() exceptions.Response {
	return exceptions.Response{
		Message: "Too many requests",
	}
}

func (TooManyRequestsException) GetHttpStatusCode() int {
	return http.StatusTooManyRequests
}
