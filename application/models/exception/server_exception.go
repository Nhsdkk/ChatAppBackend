package exception

import "net/http"

type ServerException struct {
	Err error
}

func (s ServerException) Error() string {
	return s.Err.Error()
}

func (s ServerException) GetHttpStatusCode() int {
	return http.StatusInternalServerError
}

func (s ServerException) GetResponse() Response {
	return Response{
		Message: "Unexpected server exception occurred",
	}
}
