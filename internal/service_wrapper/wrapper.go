package service_wrapper

import (
	"chat_app_backend/application/models/jwt_claims"
	"chat_app_backend/internal/jwt"
	"chat_app_backend/internal/logger"
	"chat_app_backend/internal/redis"
	"chat_app_backend/internal/s3"
	"chat_app_backend/internal/sqlc/db"

	"github.com/gin-gonic/gin"
)

type RouteHandler = func(serviceWrapper IServiceWrapper, ctx *gin.Context)

type IServiceWrapper interface {
	WrapRoute(handler RouteHandler) gin.HandlerFunc
	GetJwtHandler() jwt.IHandler[jwt_claims.UserClaims]
	GetDbConnection() db.IDbConnection
	GetLogger() logger.ILogger
	GetRedisClient() *redis.Client
	GetS3Client() s3.IClient
	Close() error
}

type ServiceWrapper struct {
	db          db.IDbConnection
	jwtHandler  jwt.IHandler[jwt_claims.UserClaims]
	logger      logger.ILogger
	redisClient *redis.Client
	s3Client    s3.IClient
}

func (wrapper *ServiceWrapper) GetS3Client() s3.IClient {
	return wrapper.s3Client
}

func (wrapper *ServiceWrapper) GetRedisClient() *redis.Client {
	return wrapper.redisClient
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
	_ = wrapper.redisClient.Close()
	return nil
}

func CreateWrapper(
	db db.IDbConnection,
	jwtHandler jwt.IHandler[jwt_claims.UserClaims],
	logger logger.ILogger,
	redisClient *redis.Client,
) IServiceWrapper {
	sw := &ServiceWrapper{}
	sw.db = db
	sw.jwtHandler = jwtHandler
	sw.logger = logger
	sw.redisClient = redisClient
	return sw
}
