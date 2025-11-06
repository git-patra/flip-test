package main

import (
	"boilerplate-go/config"
	"boilerplate-go/internal/app/server"
	"boilerplate-go/internal/pkg/statements/infrastructure/bus"

	"log"
)

func MainHttpHandler(cfg *config.AppConfig, x bus.Exchange) {
	serverHttp := server.NewServer(cfg, x)
	err := serverHttp.Start()

	if err != nil {
		log.Fatalf("error starting server: %s", err)
		return
	}

}
