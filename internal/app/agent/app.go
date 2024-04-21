package agent

import (
	"context"
	"sync"
	"time"

	compress "github.com/Stern-Ritter/metrics-and-alerting-service/internal/compress/agent"
	service "github.com/Stern-Ritter/metrics-and-alerting-service/internal/service/agent"
)

const (
	taksCount = 2
)

func Run(a *service.Agent) error {
	a.HTTPClient.OnAfterResponse(compress.GzipMiddleware)

	wg := sync.WaitGroup{}
	wg.Add(taksCount)

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
