package main

import (
	"context"
	"time"

	pkgredis "github.com/HOangAG2207/GoBeK03/pkg/redis"
)

func main() {
	ctx := context.Background()

	redisClient, err := pkgredis.NewRedisClient("")
	if err != nil {
		panic(err)
	}
	redisClientCache, err := pkgredis.NewRedisClient("CACHE")
	if err != nil {
		panic(err)
	}

	redisClient.Set(ctx, "key_normal", "1234", time.Hour)
	redisClientCache.Set(ctx, "key_cache", "12345", time.Hour)
}
