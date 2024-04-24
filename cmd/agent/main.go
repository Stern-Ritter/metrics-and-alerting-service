package main

import (
	"log"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	app "github.com/Stern-Ritter/metrics-and-alerting-service/internal/app/agent"
	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/agent"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	service "github.com/Stern-Ritter/metrics-and-alerting-service/internal/service/agent"
	storage "github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

func main() {
	config, err := app.GetConfig(config.AgentConfig{
		SendMetricsEndPoint: "/updates",
		LoggerLvl:           "info",
	})
	if err != nil {
		log.Fatalf("%+v", err)
	}

	logger, err := logger.Initialize(config.LoggerLvl)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	httpClient := resty.New()
	cache := storage.NewAgentMemCache(metrics.SupportedGaugeMetrics, metrics.SupportedCounterMetrics, logger)
	monitor := monitors.Monitor{}
	random := utils.NewRandom()
	agent := service.NewAgent(httpClient, &cache, &monitor, &random, &config, logger)

	err = app.Run(agent)
	if err != nil {
		logger.Fatal(err.Error(), zap.String("event", "start agent"))
	}
}
