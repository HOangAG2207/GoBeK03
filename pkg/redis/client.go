package pkgredis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(envprefix string) (*redis.Client, error) {
	cfg, err := newConfig(envprefix)

	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// ✅ verify connection
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return redisClient, nil
}
