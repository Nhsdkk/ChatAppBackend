package common_exceptions

import (
	"chat_app_backend/internal/exceptions"
	"fmt"
	"net/http"
)

type ForbiddenException struct {
	exceptions.BaseRestException
}

func (f ForbiddenException) GetHttpStatusCode() int {
	return http.StatusForbidden
}

func (f ForbiddenException) GetResponse() exceptions.Response {
	return exceptions.Response{
		Message: fmt.Sprintf("Forbidden: %s", f.Message),
	}
}
