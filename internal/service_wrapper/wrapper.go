package service_wrapper

import (
	"chat_app_backend/application/models/jwt_claims"
	"chat_app_backend/internal/jwt"
	"chat_app_backend/internal/logger"
	"chat_app_backend/internal/sqlc/db"
	"github.com/gin-gonic/gin"
)

type RouteHandler = func(serviceWrapper IServiceWrapper, ctx *gin.Context)

type IServiceWrapper interface {
	WrapRoute(handler RouteHandler) gin.HandlerFunc
	GetJwtHandler() jwt.IHandler[jwt_claims.UserClaims]
	GetDbConnection() db.IDbConnection
	GetLogger() logger.ILogger
	Close() error
}

type ServiceWrapper struct {
	db         db.IDbConnection
	jwtHandler jwt.IHandler[jwt_claims.UserClaims]
	logger     logger.ILogger
}

func (wrapper *ServiceWrapper) GetLogger() logger.ILogger {
	return wrapper.logger
}

func (wrapper *ServiceWrapper) WrapRoute(handler RouteHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handler(wrapper, ctx)
	}
}

func (wrapper *ServiceWrapper) GetDbConnection() db.IDbConnection {
	return wrapper.db
}

func (wrapper *ServiceWrapper) GetJwtHandler() jwt.IHandler[jwt_claims.UserClaims] {
	return wrapper.jwtHandler
}

func (wrapper *ServiceWrapper) Close() error {
	wrapper.db.Close()
	return nil
}

func CreateWrapper(
	db db.IDbConnection,
	jwtHandler jwt.IHandler[jwt_claims.UserClaims],
	logger logger.ILogger,
) IServiceWrapper {
	sw := &ServiceWrapper{}
	sw.db = db
	sw.jwtHandler = jwtHandler
	sw.logger = logger
	return sw
}
