package agent

import (
	"fmt"
	"runtime"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/go-resty/resty/v2"
)

func UpdateMetrics(cache storage.AgentCache, monitor *monitors.Monitor, rand *utils.Random) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	monitor.Update(&ms)
	randomValue, _ := rand.Float(0.1, 99.99)

	cache.UpdateMonitorMetrics(monitor)
	err := cache.UpdateGaugeMetric(metrics.NewGauge("RandomValue", randomValue))
	if err != nil {
		fmt.Println(err)
	}
}

func SendMetrics(client *resty.Client, url string, endpoint string, cache storage.AgentCache) {
	gauges, counters := cache.GetMetrics()

	err := cache.ResetMetricValue(string(metrics.Gauge), "PollCount")
	if err != nil {
		fmt.Println(err)
	}

	for _, metric := range gauges {
		_, err := sendPostRequest(client, url, endpoint, "text/plain",
			map[string]string{"type": string(metric.Type), "name": metric.Name, "value": utils.FormatGaugeMetricValue(metric.GetValue())})
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, metric := range counters {
		_, err := sendPostRequest(client, url, endpoint, "text/plain",
			map[string]string{"type": string(metric.Type), "name": metric.Name, "value": utils.FormatCounterMetricValue(metric.GetValue())})
		if err != nil {
			fmt.Println(err)
		}
	}
}
