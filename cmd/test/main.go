package main

import (
	pkglogger "github.com/HOangAG2207/GoBeK03/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	// ctx := context.Background()

	// redisClient, err := pkgredis.NewRedisClient("")
	// if err != nil {
	// 	panic(err)
	// }
	// redisClientCache, err := pkgredis.NewRedisClient("CACHE")
	// if err != nil {
	// 	panic(err)
	// }

	// redisClient.Set(ctx, "key_normal", "1234", time.Hour)
	// redisClientCache.Set(ctx, "key_cache", "12345", time.Hour)
	_ = godotenv.Load()
	pkglogger.SetLogLevel()

	log.Debug().Int("run-time", 100).Msg("This is a debug message")
}
