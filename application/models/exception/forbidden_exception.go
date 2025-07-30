package exception

import (
	"fmt"
	"net/http"
)

type ForbiddenException struct {
	Err error
}

func (f ForbiddenException) Error() string {
	return f.Err.Error()
}

func (f ForbiddenException) GetHttpStatusCode() int {
	return http.StatusForbidden
}

func (f ForbiddenException) GetResponse() Response {
	return Response{
		Message: fmt.Sprintf("Forbidden: %s", f.Err),
	}
}
