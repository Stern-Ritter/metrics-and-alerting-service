package main

import (
	"log"

	"go.uber.org/zap"
	"gopkg.in/h2non/gentleman.v2"

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
	cfg, err := app.GetConfig(config.AgentConfig{
		SendMetricsEndPoint: "/updates",
		MetricsBufferSize:   12,
		LoggerLvl:           "info",
	})
	if err != nil {
		log.Fatalf("%+v", err)
	}

	logger, err := logger.Initialize(cfg.LoggerLvl)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	httpClient := gentleman.New()
	cache := storage.NewAgentMemCache(metrics.SupportedGaugeMetrics, metrics.SupportedCounterMetrics, logger)
	runtimeMonitor := monitors.RuntimeMonitor{}
	utilMonitor := monitors.UtilMonitor{}
	random := utils.NewRandom()
	agent := service.NewAgent(httpClient, &cache, &runtimeMonitor, &utilMonitor, &random, &cfg, logger)

	err = app.Run(agent)
	if err != nil {
		logger.Fatal(err.Error(), zap.String("event", "start agent"))
	}
}
