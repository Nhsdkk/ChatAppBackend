package common_exceptions

import (
	"chat_app_backend/internal/exceptions"
	"fmt"
	"net/http"
)

type ResourceNotFoundException struct {
	exceptions.BaseRestException
}

func (r ResourceNotFoundException) GetResponse() exceptions.Response {
	return exceptions.Response{
		Message: fmt.Sprintf("Resource not found: %s", r.Message),
	}
}

func (r ResourceNotFoundException) GetHttpStatusCode() int {
	return http.StatusNotFound
}
