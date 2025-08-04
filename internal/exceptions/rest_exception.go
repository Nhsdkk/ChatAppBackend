package exceptions

import (
	"net/http"
)

type IRestException interface {
	ITrackableException
	GetHttpStatusCode() int
	GetResponse() Response
}

type BaseRestException struct {
	ITrackableException
	Message string
}

func (b BaseRestException) GetHttpStatusCode() int {
	return http.StatusInternalServerError
}

func (b BaseRestException) GetResponse() Response {
	return Response{
		Message: b.Message,
	}
}

type Response struct {
	Message string `json:"message"`
}
