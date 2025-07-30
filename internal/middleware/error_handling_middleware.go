package middleware

import (
	"chat_app_backend/application/models/exception"
	"chat_app_backend/internal/logger"
	"errors"
	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware(logger logger.ILogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if ctx.Writer.Written() {
			return
		}

		if len(ctx.Errors) == 0 {
			return
		}

		err := ctx.Errors.Last().Err

		var restException exception.IRestException
		switch {
		case errors.As(err, &restException):
		default:
			restException = exception.ServerException{
				Err: err,
			}
		}

		logger.
			CreateErrorMessage(restException).
			Log()

		ctx.JSON(restException.GetHttpStatusCode(), restException.GetResponse())
	}
}
