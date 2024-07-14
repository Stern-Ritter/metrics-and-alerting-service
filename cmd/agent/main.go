package main

import (
	"fmt"
	"log"

	"go.uber.org/zap"

	app "github.com/Stern-Ritter/metrics-and-alerting-service/internal/app/agent"
	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/agent"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/agent"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printBuildInfo()

	cfg, err := app.GetConfig(config.AgentConfig{
		SendMetricsURL:        "localhost:8080",
		SendMetricsEndPoint:   "/updates",
		UpdateMetricsInterval: 2,
		SendMetricsInterval:   5,
		MetricsBufferSize:     12,
		RateLimit:             1,
		LoggerLvl:             "info",
	})
	if err != nil {
		log.Fatalf("%+v", err)
	}

	logger, err := logger.Initialize(cfg.LoggerLvl)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	err = app.Run(&cfg, logger)
	if err != nil {
		logger.Fatal(err.Error(), zap.String("event", "start agent"))
	}
}

func printBuildInfo() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
}
