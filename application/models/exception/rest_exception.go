package exception

type IRestException interface {
	error
	GetHttpStatusCode() int
	GetResponse() Response
}

type Response struct {
	Message string `json:"message"`
}
