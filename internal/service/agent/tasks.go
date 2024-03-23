package agent

import (
	"encoding/json"
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
	_, err := cache.UpdateGaugeMetric(metrics.NewGauge("RandomValue", randomValue))
	if err != nil {
		fmt.Println(err)
	}
}

func SendMetrics(client *resty.Client, url string, endpoint string, cache storage.AgentCache) {
	gauges, counters := cache.GetMetrics()

	err := cache.ResetMetricValue(string(metrics.Counter), "PollCount")
	if err != nil {
		fmt.Println(err)
	}

	for _, gaugeMetric := range gauges {
		metric := metrics.GaugeMetricToMetrics(gaugeMetric)
		body, err := json.Marshal(metric)
		if err != nil {
			fmt.Println(err)
			continue
		}

		_, err = sendPostRequest(client, url, endpoint, "application/json", body)
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, counterMetric := range counters {
		metric := metrics.CounterMetricToMetrics(counterMetric)
		body, err := json.Marshal(metric)
		if err != nil {
			fmt.Println(err)
			continue
		}

		_, err = sendPostRequest(client, url, endpoint, "application/json", body)
		if err != nil {
			fmt.Println(err)
		}
	}
}
