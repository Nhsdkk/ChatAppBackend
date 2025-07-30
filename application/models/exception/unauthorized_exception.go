package exception

import "net/http"

type UnauthorizedException struct {
	Err error
}

func (u UnauthorizedException) Error() string {
	return u.Err.Error()
}

func (u UnauthorizedException) GetHttpStatusCode() int {
	return http.StatusUnauthorized
}

func (u UnauthorizedException) GetResponse() Response {
	return Response{
		Message: "Unauthorized",
	}
}
