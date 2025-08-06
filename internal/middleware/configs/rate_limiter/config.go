package rate_limiter

type RateLimiterConfig struct {
	RefillPerMinute int64 `env:"REFILL_PER_MINUTE"`
	MaxRequests     int64 `env:"MAX_REQUESTS"`
}
