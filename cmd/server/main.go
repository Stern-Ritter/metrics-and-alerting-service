package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"

	"go.uber.org/zap"

	app "github.com/Stern-Ritter/metrics-and-alerting-service/internal/app/server"
	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printBuildInfo()

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

func printBuildInfo() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
}
