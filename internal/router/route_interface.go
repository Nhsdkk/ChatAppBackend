package router

import (
	"chat_app_backend/internal/request_env"
	"github.com/gin-gonic/gin"
)

type HttpMethod = int

const (
	GET HttpMethod = iota
	POST
	PUT
	DELETE
	PATCH
)

type IRoute interface {
	getMethod() HttpMethod
	getPath() string
	getEndpointHandler(preferredResponseStatus int, env *request_env.RequestEnv) func(ctx *gin.Context)
}
