package main

import (
	"log"

	app "github.com/Stern-Ritter/metrics-and-alerting-service/internal/app/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model"
	service "github.com/Stern-Ritter/metrics-and-alerting-service/internal/service/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/go-resty/resty/v2"
)

func main() {
	httpClient := resty.New()
	cache := storage.NewAgentMemCache(model.SupportedGaugeMetrics, model.SupportedCounterMetrics)
	monitor := model.Monitor{}
	random := utils.NewRandom()
	agent := service.NewAgent(httpClient, &cache, &monitor, &random)

	err := app.Run(agent)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
