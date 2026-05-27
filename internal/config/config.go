package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppPort     string `envconfig:"APP_PORT" default:"8081"`
	ServiceName string `envconfig:"APP_SERVICE_NAME" default:"GoBe-K03-echo"`
	InstanceID  string `envconfig:"INSTANCE_ID" default:""`
}

func NewConfig() (*Config, error) {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}
	var cfg Config
	if err := envconfig.Process("api", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
