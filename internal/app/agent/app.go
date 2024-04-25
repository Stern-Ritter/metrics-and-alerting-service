package agent

import (
	"context"
	"sync"
	"time"

	compress "github.com/Stern-Ritter/metrics-and-alerting-service/internal/compress/agent"
	service "github.com/Stern-Ritter/metrics-and-alerting-service/internal/service/agent"
)

const (
	taskCount = 2
)

func Run(a *service.Agent) error {
	a.HTTPClient.URL(a.Config.SendMetricsURL)
	a.HTTPClient.UseHandler("before dial", compress.GzipMiddleware)
	a.HTTPClient.UseHandler("before dial", a.SignMiddleware)

	wg := sync.WaitGroup{}
	wg.Add(taskCount)

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Hour, cancel)

	updateMetricsTask := func() {
		a.UpdateMetrics()
	}
	sendMetricsTask := func() {
		a.SendMetrics()
	}

	service.SetInterval(ctx, &wg, updateMetricsTask, time.Duration(a.Config.UpdateMetricsInterval)*time.Second)
	service.SetInterval(ctx, &wg, sendMetricsTask, time.Duration(a.Config.SendMetricsInterval)*time.Second)

	wg.Wait()

	return nil
}
