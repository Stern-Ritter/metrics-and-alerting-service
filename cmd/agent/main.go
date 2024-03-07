package main

import (
	"net/http"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/app"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/config"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

func main() {
	httpClient := &http.Client{}
	cache := storage.NewAgentMemCache(model.SupportedGaugeMetrics, model.SupportedCounterMetrics)
	monitor := model.Monitor{}
	random := utils.NewRandom()

	agent := app.NewMonitoringAgent(httpClient, &cache, &monitor, &random, config.MonitoringAgentConfig)

	agent.Run()
}
