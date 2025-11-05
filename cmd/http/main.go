package main

import (
	"boilerplate-go/config"
	"boilerplate-go/internal/pkg/statements"
	"boilerplate-go/internal/pkg/statements/infrastructure/bus"
	"context"
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

	exchange := bus.NewExchange()
	// start consumers per queue
	_ = statements.InitEventConsumers(context.Background(), exchange)
	logrus.Info("[bus] consumers started")

	wg.Add(1)
	go func() {
		defer wg.Done()
		logrus.Info("Starting HTTP handler")
		MainHttpHandler(cfg, exchange)
	}()

	wg.Wait()
}
