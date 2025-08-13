package router

import (
	"chat_app_backend/internal/middleware"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/sqlc/db_queries"

	"github.com/gin-gonic/gin"
)

type AuthorizedRoute[TRequest interface{}, TResponse interface{}] struct {
	Route IRoute
}

func (a *AuthorizedRoute[TRequest, TResponse]) getMethod() HttpMethod {
	return a.Route.getMethod()
}

func (a *AuthorizedRoute[TRequest, TResponse]) getPath() string {
	return a.Route.getPath()
}

func (a *AuthorizedRoute[TRequest, TResponse]) getEndpointHandler(preferredResponseStatus int, env *request_env.RequestEnv) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userAny, exists := ctx.Get(middleware.ClaimsKey)
		if !exists {
			return
		}

		user := userAny.(*db_queries.User)
		env.User = user

		a.Route.getEndpointHandler(preferredResponseStatus, env)(ctx)
	}
}
