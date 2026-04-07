package main

import (
	"go-clean-arch/internal/di"
	"go-clean-arch/pkg/config"
	"go-clean-arch/pkg/log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("failed to load config", "error", err)
		panic(err)
	}

	server, err := di.InitializeAPI(cfg)
	if err != nil {
		log.Error("failed to initialize server", "error", err)
		panic(err)
	}

	if err := server.Start(); err != nil {
		log.Error("failed to start server", "error", err)
		panic(err)
	}
}
