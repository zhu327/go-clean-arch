package main

import (
	"go-wire/pkg/config"
	"go-wire/pkg/di"
	"log"
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

	if server.Start() != nil {
		log.Fatal("error starting server: ", err)
	}

}
