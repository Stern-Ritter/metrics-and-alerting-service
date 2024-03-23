package agent

import (
	"context"
	"sync"
	"time"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/agent"
	service "github.com/Stern-Ritter/metrics-and-alerting-service/internal/service/agent"
)

const (
	taksCount = 2
)

func Run(a *service.Agent) error {
	config, err := getConfig(config.AgentConfig{
		SendMetricsEndPoint: "/update",
	})
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(taksCount)

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Hour, cancel)

	updateMetricsTask := func() {
		service.UpdateMetrics(a.Cache, a.Monitor, a.Random)
	}
	sendMetricsTask := func() {
		service.SendMetrics(a.HTTPClient, config.SendMetricsURL, config.SendMetricsEndPoint, a.Cache)
	}

	service.SetInterval(ctx, &wg, updateMetricsTask, time.Duration(config.UpdateMetricsInterval)*time.Second)
	service.SetInterval(ctx, &wg, sendMetricsTask, time.Duration(config.SendMetricsInterval)*time.Second)

	wg.Wait()

	return nil
}
