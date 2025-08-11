package application

import (
	"chat_app_backend/application/application_config"
	controllers "chat_app_backend/application/controllers/users"
	"chat_app_backend/application/models/jwt_claims"
	"chat_app_backend/internal/configuration"
	"chat_app_backend/internal/env_loader"
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/jwt"
	logger2 "chat_app_backend/internal/logger"
	"chat_app_backend/internal/middleware"
	"chat_app_backend/internal/middleware/configs/rate_limiter"
	"chat_app_backend/internal/redis"
	"chat_app_backend/internal/s3"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/sqlc/db"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

type IApplication interface {
	Configure()
	Serve()
	Close()
}

type Application struct {
	config         *application_config.ApplicationConfig
	engine         *gin.Engine
	server         *http.Server
	serviceWrapper service_wrapper.IServiceWrapper
	configuration  configuration.IConfiguration
}

func (appl *Application) Close() {
	if err := appl.server.Close(); err != nil {
		log.Fatalf("Can't close server: %s", err)
	}

	if err := appl.serviceWrapper.Close(); err != nil {
		log.Fatalf("Can't close services: %s", err)
	}
}

func (appl *Application) Configure() {
	appl.loadConfigurations()
	appl.configureServices()
	appl.createServer()
	appl.configureMiddleware()

	appl.configureRoutes()
}

func (appl *Application) createServer() {
	applConfig, err := appl.configuration.Get(&application_config.ApplicationConfig{})
	if err != nil {
		appl.serviceWrapper.GetLogger().
			CreateErrorMessage(exceptions.WrapErrorWithTrackableException(err)).
			WithFatal().
			Log()
		return
	}

	appl.config = applConfig.(*application_config.ApplicationConfig)
	appl.server = &http.Server{
		Addr:    appl.config.Url,
		Handler: appl.engine,
	}
}

func (appl *Application) configureRoutes() {
	controllers.CreateUserController(appl.engine, appl.serviceWrapper).ConfigureGroup()
}

func (appl *Application) configureMiddleware() {
	rateLimiterConfig, err := appl.configuration.Get(&rate_limiter.RateLimiterConfig{})
	if err != nil {
		appl.serviceWrapper.GetLogger().
			CreateErrorMessage(exceptions.WrapErrorWithTrackableException(err)).
			WithFatal().
			Log()

		return
	}

	appl.engine.Use(
		middleware.RequestLoggingMiddleware(appl.serviceWrapper.GetLogger()),
		middleware.ErrorHandlerMiddleware(appl.serviceWrapper.GetLogger()),
		middleware.RateLimiterMiddleware(rateLimiterConfig.(*rate_limiter.RateLimiterConfig), appl.serviceWrapper, appl.config),
		middleware.AuthorizationMiddleware(
			appl.serviceWrapper.GetJwtHandler(),
			appl.serviceWrapper.GetDbConnection(),
		),
	)
}

func (appl *Application) loadConfigurations() {
	dbConfiguration := &db.PostgresConfig{}
	jwtConfig := &jwt.JwtConfig{}
	redisConfig := &redis.RedisConfig{}
	rateLimiterConfig := &rate_limiter.RateLimiterConfig{}
	s3Config := &s3.S3Config{}
	applicationConfig := &application_config.ApplicationConfig{}
	envLoader := env_loader.CreateLoaderFromEnv()

	dbConfigurationLoadingErr := envLoader.LoadDataIntoStruct(dbConfiguration)
	if dbConfigurationLoadingErr != nil {
		log.Fatal(dbConfigurationLoadingErr)
	}

	jwtConfigurationLoadingError := envLoader.LoadDataIntoStruct(jwtConfig)
	if jwtConfigurationLoadingError != nil {
		log.Fatal(jwtConfigurationLoadingError)
	}

	redisConfigurationLoadingError := envLoader.LoadDataIntoStruct(redisConfig)
	if redisConfigurationLoadingError != nil {
		log.Fatal(redisConfigurationLoadingError)
	}

	rateLimiterConfigurationLoadingError := envLoader.LoadDataIntoStruct(rateLimiterConfig)
	if rateLimiterConfigurationLoadingError != nil {
		log.Fatal(rateLimiterConfigurationLoadingError)
	}

	applicationConfigLoadingError := envLoader.LoadDataIntoStruct(applicationConfig)
	if applicationConfigLoadingError != nil {
		log.Fatal(applicationConfigLoadingError)
	}

	s3ConfigLoadingError := envLoader.LoadDataIntoStruct(s3Config)
	if s3ConfigLoadingError != nil {
		log.Fatal(s3ConfigLoadingError)
	}

	appl.configuration = configuration.CreateConfiguration().
		AddConfiguration(jwtConfig).
		AddConfiguration(dbConfiguration).
		AddConfiguration(redisConfig).
		AddConfiguration(rateLimiterConfig).
		AddConfiguration(applicationConfig).
		AddConfiguration(s3Config)
}

func (appl *Application) configureServices() {
	ctx := context.Background()

	logger := logger2.CreateLogger(os.Stdout)

	dbConnection, dbBuildError := configuration.BuildFromConfiguration[db.Connection](
		appl.configuration,
		db.CreateConnection,
		&ctx,
	)

	if dbBuildError != nil {
		logger.
			CreateErrorMessage(exceptions.WrapErrorWithTrackableException(dbBuildError)).
			WithFatal().
			Log()

		return
	}

	jwtHandler, jwtHandlerBuildError := configuration.BuildFromConfiguration[jwt.Handler[jwt_claims.UserClaims]](
		appl.configuration,
		jwt.CreateJwtHandler[jwt_claims.UserClaims],
	)

	if jwtHandlerBuildError != nil {
		logger.
			CreateErrorMessage(exceptions.WrapErrorWithTrackableException(jwtHandlerBuildError)).
			WithFatal().
			Log()

		return
	}

	redisClient, redisClientBuildError := configuration.BuildFromConfiguration[redis.Client](
		appl.configuration,
		redis.CreateRedisClient,
		&ctx,
	)

	if redisClientBuildError != nil {
		logger.
			CreateErrorMessage(exceptions.WrapErrorWithTrackableException(redisClientBuildError)).
			WithFatal().
			Log()

		return
	}

	s3Client, s3ClientCreationError := configuration.BuildFromConfiguration[s3.Client](
		appl.configuration,
		s3.CreateClient,
	)

	if s3ClientCreationError != nil {
		logger.
			CreateErrorMessage(exceptions.WrapErrorWithTrackableException(s3ClientCreationError)).
			WithFatal().
			Log()

		return
	}

	appl.serviceWrapper = service_wrapper.CreateWrapper(
		dbConnection,
		jwtHandler,
		logger,
		redisClient,
		s3Client,
	)
}

func (appl *Application) Serve() {
	go func() {
		if err := appl.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("run error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appl.Close()
}

func Create() *Application {
	engine := gin.New()
	return &Application{
		engine:         engine,
		serviceWrapper: nil,
		configuration:  nil,
	}
}
