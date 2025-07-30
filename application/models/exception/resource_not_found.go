package exception

import (
	"fmt"
	"net/http"
)

type ResourceNotFoundException struct {
	Err error
}

func (r ResourceNotFoundException) GetResponse() Response {
	return Response{
		Message: fmt.Sprintf("Resource not found: %s", r.Err),
	}
}

func (r ResourceNotFoundException) Error() string {
	return r.Err.Error()
}

func (r ResourceNotFoundException) GetHttpStatusCode() int {
	return http.StatusNotFound
}
