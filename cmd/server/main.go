package main

import (
	"log"

	app "github.com/Stern-Ritter/metrics-and-alerting-service/internal/app/server"
	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	service "github.com/Stern-Ritter/metrics-and-alerting-service/internal/service/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"go.uber.org/zap"
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

	storage := storage.NewServerMemStorage(logger)
	server := service.NewServer(&storage, &config, logger)

	err = app.Run(server)
	if err != nil {
		logger.Fatal(err.Error(), zap.String("event", "start server"))
	}
}
