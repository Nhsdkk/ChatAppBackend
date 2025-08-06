package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func CreateRedisClient(config *RedisConfig, ctx context.Context) (*Client, error) {
	client := redis.NewClient(
		&redis.Options{
			Addr:       fmt.Sprintf("%s:%d", config.Host, config.Port),
			ClientName: config.ClientName,
			Username:   config.Username,
			Password:   config.Password,
			DB:         config.DB,
		},
	)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{
		Client: client,
	}, nil
}
