package router

import (
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/exceptions/common_exceptions"
	"chat_app_backend/internal/handler"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/validator"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouteType int

const (
	Authorized RouteType = iota
	Base
)

type BaseRoute[TRequest interface{}, TResponse interface{}] struct {
	validator validator.IValidator[TRequest]
	handler   handler.HFunc[TRequest, TResponse, request_env.RequestEnv]
	wrapper   service_wrapper.IServiceWrapper
	path      string
	method    HttpMethod
}

func (r *BaseRoute[TRequest, TResponse]) getMethod() HttpMethod {
	return r.method
}

func (r *BaseRoute[TRequest, TResponse]) getPath() string {
	return r.path
}

func (r *BaseRoute[TRequest, TResponse]) getEndpointHandler(preferredResponseStatus int, env *request_env.RequestEnv) func(ctx *gin.Context) {
	return r.wrapper.WrapRoute(
		func(serviceWrapper service_wrapper.IServiceWrapper, ctx *gin.Context) {
			var requestDto TRequest

			if err := ctx.ShouldBind(&requestDto); err != nil {
				_ = ctx.Error(
					common_exceptions.InvalidBodyException{
						BaseRestException: exceptions.BaseRestException{
							ITrackableException: exceptions.WrapErrorWithTrackableException(err),
							Message:             "can't deserialize request",
						},
					},
				)
				return
			}

			if err := ctx.ShouldBindQuery(&requestDto); err != nil {
				_ = ctx.Error(
					common_exceptions.InvalidBodyException{
						BaseRestException: exceptions.BaseRestException{
							ITrackableException: exceptions.WrapErrorWithTrackableException(err),
							Message:             "can't deserialize request",
						},
					},
				)
				return
			}

			if err := ctx.ShouldBindUri(&requestDto); err != nil {
				_ = ctx.Error(
					common_exceptions.InvalidBodyException{
						BaseRestException: exceptions.BaseRestException{
							ITrackableException: exceptions.WrapErrorWithTrackableException(err),
							Message:             "can't deserialize request",
						},
					},
				)
				return
			}

			if validationError := r.validator.Validate(&requestDto); validationError != nil {
				var restException exceptions.IRestException

				switch {
				case errors.As(validationError, &restException):
				default:
					restException = common_exceptions.InvalidBodyException{
						BaseRestException: exceptions.BaseRestException{
							ITrackableException: exceptions.WrapErrorWithTrackableException(validationError),
							Message:             validationError.Error(),
						},
					}
				}

				_ = ctx.Error(restException)
				return
			}

			response, handlerError := r.handler(&requestDto, serviceWrapper, ctx, env)
			if handlerError != nil {
				_ = ctx.Error(handlerError)
				return
			}

			if response == nil {
				ctx.Status(preferredResponseStatus)
				return
			}

			ctx.JSON(preferredResponseStatus, response)
		},
	)
}

func RegisterRoute(router *gin.RouterGroup, route IRoute) {
	switch route.getMethod() {
	case POST:
		router.POST(route.getPath(), route.getEndpointHandler(http.StatusCreated, &request_env.RequestEnv{}))
		break
	case GET:
		router.GET(route.getPath(), route.getEndpointHandler(http.StatusOK, &request_env.RequestEnv{}))
		break
	case PUT:
		router.PUT(route.getPath(), route.getEndpointHandler(http.StatusOK, &request_env.RequestEnv{}))
		break
	case PATCH:
		router.PATCH(route.getPath(), route.getEndpointHandler(http.StatusOK, &request_env.RequestEnv{}))
		break
	case DELETE:
		router.DELETE(route.getPath(), route.getEndpointHandler(http.StatusOK, &request_env.RequestEnv{}))
		break
	}
}

func CreateBaseRoute[TRequest interface{}, TResponse interface{}](
	wrapper service_wrapper.IServiceWrapper,
	path string,
	handlerFunc handler.HFunc[TRequest, TResponse, request_env.RequestEnv],
	validator validator.IValidator[TRequest],
	method HttpMethod,
) IRoute {
	baseRoute := BaseRoute[TRequest, TResponse]{
		validator: validator,
		handler:   handlerFunc,
		wrapper:   wrapper,
		path:      path,
		method:    method,
	}

	return &baseRoute
}
