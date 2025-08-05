package middleware

import (
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/exceptions/common_exceptions"
	"chat_app_backend/internal/middleware/configs/rate_limiter"
	"chat_app_backend/internal/service_wrapper"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

func RateLimiterMiddleware(config *rate_limiter.RateLimiterConfig, services service_wrapper.IServiceWrapper) gin.HandlerFunc {
	if config.RefillPerMinute == 0 || config.MaxRequests == 0 {
		services.GetLogger().
			CreateErrorMessage(exceptions.CreateTrackableExceptionFromStringF("refill per minute and max requests can't be zero")).
			WithFatal().
			Log()

		return nil
	}

	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		newEntry := false
		lastRefill, err := services.GetRedisClient().Get(ctx, fmt.Sprintf("%s:last_refill", ip)).Time()
		switch {
		case errors.Is(err, redis.Nil):
			lastRefill = time.Now()
			newEntry = true
		case err != nil:
			_ = ctx.Error(exceptions.WrapErrorWithTrackableException(err))
			ctx.Abort()
			return
		}

		storedCount, err := services.GetRedisClient().Get(ctx, fmt.Sprintf("%s:current_rate", ip)).Int64()
		switch {
		case errors.Is(err, redis.Nil):
			storedCount = 0
		case err != nil:
			_ = ctx.Error(exceptions.WrapErrorWithTrackableException(err))
			ctx.Abort()
			return
		}

		elapsedMinutes := int64(time.Since(lastRefill).Minutes())

		if newEntry || elapsedMinutes > (config.MaxRequests-storedCount)/config.RefillPerMinute {
			storedCount = config.MaxRequests
		} else {
			storedCount += elapsedMinutes * config.RefillPerMinute
		}

		storedCount = max(storedCount-1, -1)

		pipe := services.GetRedisClient().TxPipeline()
		pipe.Set(ctx, fmt.Sprintf("%s:current_rate", ip), storedCount, 0)

		if newEntry || elapsedMinutes != 0 {
			var timeToSet time.Time
			if newEntry {
				timeToSet = time.Now()
			} else {
				timeToSet = lastRefill.Add(time.Minute * time.Duration(elapsedMinutes))
			}
			pipe.Set(ctx, fmt.Sprintf("%s:last_refill", ip), timeToSet, 0)
		}

		if _, err := pipe.Exec(ctx); err != nil {
			_ = ctx.Error(exceptions.WrapErrorWithTrackableException(err))
			ctx.Abort()
			return
		}

		if storedCount < 0 {
			_ = ctx.Error(
				common_exceptions.TooManyRequestsException{
					BaseRestException: exceptions.BaseRestException{
						ITrackableException: exceptions.CreateTrackableExceptionFromStringF("too many requests from %s ip", ip),
						Message:             "",
					},
				},
			)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
