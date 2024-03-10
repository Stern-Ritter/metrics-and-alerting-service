package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/config"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	tasks "github.com/Stern-Ritter/metrics-and-alerting-service/internal/transport"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/caarlos0/env"
	"github.com/go-resty/resty/v2"
)

const (
	taksCount = 2
)

func main() {
	config, err := getConfig(config.AgentConfig{
		SendMetricsEndPoint: "/update/{type}/{name}/{value}",
	})
	if err != nil {
		log.Fatal(err)
	}

	httpClient := resty.New()
	cache := storage.NewAgentMemCache(model.SupportedGaugeMetrics, model.SupportedCounterMetrics)
	monitor := model.Monitor{}
	random := utils.NewRandom()

	run(httpClient, &cache, &monitor, &random, config)
}

func run(httpClient *resty.Client, cache storage.AgentCache, monitor *model.Monitor, random *utils.Random, config config.AgentConfig) {
	var wg sync.WaitGroup
	wg.Add(taksCount)

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Hour, cancel)

	updateMetricsTask := func() {
		tasks.UpdateMetrics(cache, monitor, random)
	}
	sendMetricsTask := func() {
		tasks.SendMetrics(httpClient, config.SendMetricsURL, config.SendMetricsEndPoint, cache)
	}

	tasks.SetInterval(ctx, &wg, updateMetricsTask, time.Duration(config.UpdateMetricsInterval)*time.Second)
	tasks.SetInterval(ctx, &wg, sendMetricsTask, time.Duration(config.SendMetricsInterval)*time.Second)

	wg.Wait()
}

func getConfig(c config.AgentConfig) (config.AgentConfig, error) {
	err := c.ParseFlags()
	if err != nil {
		return c, err
	}

	err = env.Parse(&c)
	if err != nil {
		return c, err
	}

	return c, nil
}
