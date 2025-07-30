package handler

import (
	"chat_app_backend/internal/service_wrapper"
	"github.com/gin-gonic/gin"
)

type IHandler[TRequest interface{}, TResponse interface{}, TEnv interface{}] interface {
	Handle(
		request *TRequest,
		service service_wrapper.IServiceWrapper,
		ctx *gin.Context,
		requestEnvironment *TEnv,
	) (*TResponse, error)
}

type HFunc[TRequest interface{}, TResponse interface{}, TEnv interface{}] = func(
	request *TRequest,
	service service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	requestEnvironment *TEnv,
) (*TResponse, error)
