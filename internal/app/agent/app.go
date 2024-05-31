package agent

import (
	"context"
	"sync"
	"time"

	compress "github.com/Stern-Ritter/metrics-and-alerting-service/internal/compress/agent"
	service "github.com/Stern-Ritter/metrics-and-alerting-service/internal/service/agent"
)

const (
	taskCount = 3
)

func Run(a *service.Agent) error {
	a.HTTPClient.URL(a.Config.SendMetricsURL)
	a.HTTPClient.UseHandler("before dial", compress.GzipMiddleware)
	a.HTTPClient.UseHandler("before dial", a.SignMiddleware)

	wg := sync.WaitGroup{}
	wg.Add(taskCount)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.StartSendMetricsWorkerPool()

	service.SetInterval(ctx, &wg, a.UpdateRuntimeMetrics, time.Duration(a.Config.UpdateMetricsInterval)*time.Second)
	service.SetInterval(ctx, &wg, a.UpdateUtilMetrics, time.Duration(a.Config.UpdateMetricsInterval)*time.Second)
	service.SetInterval(ctx, &wg, a.SendMetrics, time.Duration(a.Config.SendMetricsInterval)*time.Second)

	wg.Wait()
	a.StopTasks()

	return nil
}
