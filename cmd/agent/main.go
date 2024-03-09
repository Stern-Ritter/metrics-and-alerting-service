package main

import (
	"flag"
	"time"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/app"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/config"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/go-resty/resty/v2"
)

var monitoringAgentConfig = config.AgentConfig{
	SendMetricsURL: config.URL{
		Host: "https://localhost",
		Port: 8080,
	},
	SendMetricsEndPoint: "/update/{type}/{name}/{value}",
}

func parseFlags() {
	flag.Var(&monitoringAgentConfig.SendMetricsURL, "a", "address and port to run server in format <host>:<port>")
	flag.DurationVar(&monitoringAgentConfig.UpdateMetricsInterval, "p", time.Duration(2)*time.Second, "interval to update metrics in seconds")
	flag.DurationVar(&monitoringAgentConfig.SendMetricsInterval, "r", time.Duration(10)*time.Second, "interval for sending metrics to the server in seconds")
	flag.Parse()
}

func main() {
	parseFlags()

	httpClient := resty.New()
	cache := storage.NewAgentMemCache(model.SupportedGaugeMetrics, model.SupportedCounterMetrics)
	monitor := model.Monitor{}
	random := utils.NewRandom()

	agent := app.NewMonitoringAgent(httpClient, &cache, &monitor, &random, monitoringAgentConfig)

	agent.Run()
}
