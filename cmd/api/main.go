package main

import (
	"log"

	"go-clean-arch/internal/di"
	"go-clean-arch/pkg/config"
)

func main() {
	config, error := config.LoadConfig()

	if error != nil {
		log.Fatal("error loading config: ", error)
	}

	server, err := di.InitailizeApi(config)
	if err != nil {
		log.Fatal("error initializing server: ", err)
	}

	if err := server.Start(); err != nil {
		log.Fatal("error starting server: ", err)
	}
}
