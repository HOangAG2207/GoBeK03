package pkgredis

import "github.com/redis/go-redis/v9"

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
	return redisClient, nil
}
