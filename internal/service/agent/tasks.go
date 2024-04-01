package agent

import (
	"encoding/json"
	"runtime"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func UpdateMetrics(cache storage.AgentCache, monitor *monitors.Monitor, rand *utils.Random, logger *zap.Logger) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	monitor.Update(&ms)
	randomValue, _ := rand.Float(0.1, 99.99)

	cache.UpdateMonitorMetrics(monitor)
	_, err := cache.UpdateGaugeMetric(metrics.NewGauge("RandomValue", randomValue))
	if err != nil {
		logger.Error(err.Error(), zap.String("event", "update RandomValue gauge metric"))
	}
}

func SendMetrics(client *resty.Client, url string, endpoint string, cache storage.AgentCache, logger *zap.Logger) {
	gauges, counters := cache.GetMetrics()

	err := cache.ResetMetricValue(string(metrics.Counter), "PollCount")
	if err != nil {
		logger.Error(err.Error(), zap.String("event", "reset PollCount counter metric"))
	}

	for _, gaugeMetric := range gauges {
		metric := metrics.GaugeMetricToMetrics(gaugeMetric)
		body, err := json.Marshal(metric)
		if err != nil {
			logger.Error(err.Error(), zap.String("event", "JSON encoding gauge metric"))
			continue
		}

		_, err = sendPostRequest(client, url, endpoint, "application/json", body)
		if err != nil {
			logger.Error(err.Error(), zap.String("event", "send update gauge metric"))
		}
	}

	for _, counterMetric := range counters {
		metric := metrics.CounterMetricToMetrics(counterMetric)
		body, err := json.Marshal(metric)
		if err != nil {
			logger.Error(err.Error(), zap.String("event", "JSON encoding counter metric"))
			continue
		}

		_, err = sendPostRequest(client, url, endpoint, "application/json", body)
		if err != nil {
			logger.Error(err.Error(), zap.String("event", "send update counter metric"))
		}
	}
}
