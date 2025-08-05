package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func CreateRedisClient(config *RedisConfig) (*Client, error) {
	client := redis.NewClient(
		&redis.Options{
			Addr:       fmt.Sprintf("%s:%d", config.Host, config.Port),
			ClientName: config.ClientName,
			Username:   config.Username,
			Password:   config.Password,
			DB:         config.DB,
		},
	)

	return &Client{
		Client: client,
	}, nil
}
