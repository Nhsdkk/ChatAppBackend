package router

import (
	"chat_app_backend/internal/middleware"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/sqlc/db_queries"
	"github.com/gin-gonic/gin"
)

type AuthorizedRoute[TRequest interface{}, TResponse interface{}] struct {
	route BaseRoute[TRequest, TResponse]
}

func (a *AuthorizedRoute[TRequest, TResponse]) getMethod() HttpMethod {
	return a.route.getMethod()
}

func (a *AuthorizedRoute[TRequest, TResponse]) getPath() string {
	return a.route.getPath()
}

func (a *AuthorizedRoute[TRequest, TResponse]) getEndpointHandler(preferredResponseStatus int, env *request_env.RequestEnv) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userAny, exists := ctx.Get(middleware.ClaimsKey)
		if !exists {
			return
		}

		user := userAny.(*db_queries.User)
		env.User = user

		a.route.getEndpointHandler(preferredResponseStatus, env)(ctx)
	}
}
