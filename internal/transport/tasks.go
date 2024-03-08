package transport

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/go-resty/resty/v2"
)

func UpdateMetrics(cache storage.AgentCache, monitor *model.Monitor, rand *utils.Random) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	monitor.Update(&ms)
	randomValue, _ := rand.Float(0.1, 99.99)

	cache.UpdateMonitorMetrics(monitor)
	cache.UpdateGaugeMetric(model.NewGauge("RandomValue", randomValue))
}

func SendMetrics(client *resty.Client, endpoint string, cache storage.AgentCache) {
	gauges, counters := cache.GetMetrics()
	cache.ResetMetricValue(string(model.Gauge), "PollCount")

	for _, metric := range gauges {
		_, err := sendPostRequest(client, endpoint, "text/plain",
			map[string]string{"type": string(metric.Type), "name": metric.Name, "value": utils.FormatGaugeMetricValue(metric.GetValue())})
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, metric := range counters {
		_, err := sendPostRequest(client, endpoint, "text/plain",
			map[string]string{"type": string(metric.Type), "name": metric.Name, "value": utils.FormatCounterMetricValue(metric.GetValue())})
		if err != nil {
			fmt.Println(err)
		}
	}
}

func sendPostRequest(client *resty.Client, endpoint, contentType string, pathParams map[string]string) (*resty.Response, error) {
	resp, err := client.R().
		SetHeader("Content-Type", contentType).
		SetPathParams(pathParams).
		Get(endpoint)

	return resp, err
}

func SetInterval(ctx context.Context, wg *sync.WaitGroup, task func(), interval time.Duration) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				wg.Done()
				return
			default:
				task()
				time.Sleep(interval)
			}
		}
	}()
}
