package main

import (
	"log"

	"github.com/HOangAG2207/GoBeK03/internal/api"
	"github.com/HOangAG2207/GoBeK03/internal/config"
)

func main() {
	cfg, err := config.NewConfig() // load config first
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	app := api.NewEngine(&api.EngineOpts{
		Cfg: cfg,
	})

	if err := app.Start(); err != nil {
		panic(err)
	}
}
