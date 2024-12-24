package main

import (
	"gotrics-server/internal/api"
	"gotrics-server/internal/config"
	"gotrics-server/internal/logger"
	"log"
)

func main() {

	logger.Init()

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	api.StartServer(cfg)
}
