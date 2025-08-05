package middleware

import (
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/exceptions/common_exceptions"
	"chat_app_backend/internal/middleware/configs/rate_limiter"
	redisinternal "chat_app_backend/internal/redis"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

func RateLimiterMiddleware(config *rate_limiter.RateLimiterConfig, redisClient *redisinternal.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		newEntry := false
		lastRefill, err := redisClient.Get(ctx, fmt.Sprintf("%s:last_refill", ip)).Time()
		switch {
		case errors.Is(err, redis.Nil):
			lastRefill = time.Now()
			newEntry = true
		case err != nil:
			_ = ctx.Error(exceptions.WrapErrorWithTrackableException(err))
			ctx.Abort()
			return
		}

		storedCount, err := redisClient.Get(ctx, fmt.Sprintf("%s:current_rate", ip)).Int64()
		switch {
		case errors.Is(err, redis.Nil):
			storedCount = 0
		case err != nil:
			_ = ctx.Error(exceptions.WrapErrorWithTrackableException(err))
			ctx.Abort()
			return
		}

		diff := int64(time.Since(lastRefill).Minutes())

		if newEntry || diff > (config.MaxRequests-storedCount)/config.RefillPerMinute {
			storedCount = config.MaxRequests
		} else {
			storedCount += diff * config.RefillPerMinute
		}

		storedCount = max(storedCount-1, -1)

		pipe := redisClient.TxPipeline()
		pipe.Set(ctx, fmt.Sprintf("%s:current_rate", ip), storedCount, 0)

		if newEntry || diff != 0 {
			var timeToSet time.Time
			if newEntry {
				timeToSet = time.Now()
			} else {
				timeToSet = lastRefill.Add(time.Minute * time.Duration(diff))
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
