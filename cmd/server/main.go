package main

import (
	"log"

	app "github.com/Stern-Ritter/metrics-and-alerting-service/internal/app/server"
	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
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
