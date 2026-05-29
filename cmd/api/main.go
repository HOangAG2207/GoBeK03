package main

import (
	"log"

	"github.com/HOangAG2207/GoBeK03/internal/api"
	pkgredis "github.com/HOangAG2207/GoBeK03/pkg/redis"
)

// @title           GoBe K03 project API
// @version         1.0
// @description     API for GoBe K03
// @host            localhost:8080
// @BasePath 		/api
func main() {
	cfg, err := api.NewConfig() // load config first
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	redisClient, err := pkgredis.NewRedisClient("")
	if err != nil {
		panic(err)
	}
	app := api.NewEngine(&api.EngineOpts{
		Cfg:         cfg,
		RedisClient: redisClient,
	})
	if err := app.Start(); err != nil {
		panic(err)
	}
}
