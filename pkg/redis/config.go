package pkgredis

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Address  string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
}

func newConfig(envprefix string) (*config, error) {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}
	var cfg config
	if err := envconfig.Process(envprefix, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
