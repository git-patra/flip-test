package main

import (
	"boilerplate-go/config"
	"log"
	"sync"

	"github.com/sirupsen/logrus"
)

func main() {
	configFilePath := ".env.yaml"
	cfg, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("error loading config: %s", err)
		return
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		logrus.Info("Starting HTTP handler")
		MainHttpHandler(cfg)
	}()

	wg.Wait()
}
