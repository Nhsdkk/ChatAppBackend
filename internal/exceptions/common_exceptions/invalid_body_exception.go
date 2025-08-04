package common_exceptions

import (
	"chat_app_backend/internal/exceptions"
	"fmt"
	"net/http"
)

type InvalidBodyException struct {
	exceptions.BaseRestException
}

func (i InvalidBodyException) GetHttpStatusCode() int {
	return http.StatusBadRequest
}

func (i InvalidBodyException) GetResponse() exceptions.Response {
	return exceptions.Response{
		Message: fmt.Sprintf("Bad request: %s", i.Message),
	}
}
