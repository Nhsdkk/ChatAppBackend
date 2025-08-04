package common_exceptions

import (
	"chat_app_backend/internal/exceptions"
	"net/http"
)

type ServerException struct {
	exceptions.BaseRestException
}

func (s ServerException) GetHttpStatusCode() int {
	return http.StatusInternalServerError
}

func (s ServerException) GetResponse() exceptions.Response {
	return exceptions.Response{
		Message: "Unexpected server exception occurred",
	}
}
