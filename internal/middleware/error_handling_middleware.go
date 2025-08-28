package middleware

import (
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/exceptions/common_exceptions"
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

		var err exceptions.ITrackableException
		errors.As(ctx.Errors.Last().Err, &err)

		var restException exceptions.IRestException
		switch {
		case errors.As(err, &restException):
		default:
			restException = common_exceptions.ServerException{
				BaseRestException: exceptions.BaseRestException{
					ITrackableException: err,
					Message:             "",
				},
			}
		}

		logger.
			CreateErrorMessage(restException).
			Log()

		ctx.JSON(restException.GetHttpStatusCode(), restException.GetResponse())
	}
}
