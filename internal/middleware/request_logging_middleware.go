package middleware

import (
	"bytes"
	"chat_app_backend/internal/logger"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

func RequestLoggingMiddleware(logger logger.ILogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		incomingRequest := ctx.Request

		bodyRaw := ctx.Request.Body
		body := "EMPTY"

		if byteData, err := io.ReadAll(bodyRaw); err == nil {
			body = string(byteData)
			_ = incomingRequest.Body.Close()
			incomingRequest.Body = io.NopCloser(bytes.NewBuffer(byteData))
		}

		logger.
			CreateInfoMessageF(
				`REQUEST [%s] %s
Headers: %v
RequestBody: %v`,
				incomingRequest.Method,
				incomingRequest.URL,
				incomingRequest.Header,
				body,
			).Log()

		start := time.Now()
		ctx.Next()
		duration := time.Since(start)

		logger.
			CreateInfoMessageF(
				`RESPONSE [%s] %s 
Status: %d
Time taken: %d ms`,
				incomingRequest.Method,
				incomingRequest.URL,
				ctx.Writer.Status(),
				duration.Milliseconds(),
			).Log()
	}
}
