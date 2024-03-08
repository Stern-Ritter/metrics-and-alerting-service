package main

import (
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/app"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/config"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/go-resty/resty/v2"
)

func main() {
	httpClient := resty.New()
	cache := storage.NewAgentMemCache(model.SupportedGaugeMetrics, model.SupportedCounterMetrics)
	monitor := model.Monitor{}
	random := utils.NewRandom()

	agent := app.NewMonitoringAgent(httpClient, &cache, &monitor, &random, config.MonitoringAgentConfig)

	agent.Run()
}
