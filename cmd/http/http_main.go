package main

import (
	"boilerplate-go/config"
	"boilerplate-go/internal/app/server"

	"log"
)

func MainHttpHandler(cfg *config.AppConfig) {
	serverHttp := server.NewServer(cfg)
	err := serverHttp.Start()

	if err != nil {
		log.Fatalf("error starting server: %s", err)
		return
	}

}
