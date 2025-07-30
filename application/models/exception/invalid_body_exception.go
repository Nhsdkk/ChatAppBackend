package exception

import (
	"fmt"
	"net/http"
)

type InvalidBodyException struct {
	Err error
}

func (i InvalidBodyException) Error() string {
	return i.Err.Error()
}

func (i InvalidBodyException) GetHttpStatusCode() int {
	return http.StatusBadRequest
}

func (i InvalidBodyException) GetResponse() Response {
	return Response{
		Message: fmt.Sprintf("Bad request: %s", i.Err),
	}
}
