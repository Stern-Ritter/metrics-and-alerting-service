package main

import (
	"log"

	"go.uber.org/zap"

	app "github.com/Stern-Ritter/metrics-and-alerting-service/internal/app/server"
	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
)

func main() {
	config, err := app.GetConfig(config.ServerConfig{
		LoggerLvl: "info",
	})
	if err != nil {
		log.Fatalf("%+v", err)
	}

	logger, err := logger.Initialize(config.LoggerLvl)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	err = app.Run(&config, logger)
	if err != nil {
		logger.Fatal(err.Error(), zap.String("event", "start server"))
	}
}
