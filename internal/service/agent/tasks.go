package agent

import (
	"encoding/json"
	"fmt"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	cache "github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/cenkalti/backoff/v4"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var (
	sendMetricsBatchRetryIntervals = backoff.NewExponentialBackOff(
		backoff.WithInitialInterval(1*time.Second),
		backoff.WithRandomizationFactor(0),
		backoff.WithMultiplier(3),
		backoff.WithMaxInterval(5*time.Second),
		backoff.WithMaxElapsedTime(10*time.Second))
)

func UpdateMetrics(cache cache.AgentCache, monitor *monitors.Monitor, rand *utils.Random, logger *zap.Logger) {
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

func SendMetrics(client *resty.Client, url string, endpoint string, cache cache.AgentCache, logger *zap.Logger) {
	gauges, counters := cache.GetMetrics()

	err := cache.ResetMetricValue(string(metrics.Counter), "PollCount")
	if err != nil {
		logger.Error(err.Error(), zap.String("event", "reset PollCount counter metric"))
	}

	metricsBatch := make([]metrics.Metrics, 0)

	for _, gaugeMetric := range gauges {
		metricsBatch = append(metricsBatch, metrics.GaugeMetricToMetrics(gaugeMetric))

	}

	for _, counterMetric := range counters {
		metricsBatch = append(metricsBatch, metrics.CounterMetricToMetrics(counterMetric))

	}

	body, err := json.Marshal(metricsBatch)
	if err != nil {
		logger.Error(err.Error(), zap.String("event", "JSON encoding metrics batch"))
	}

	sendMetricsBatch := func() error {
		resp, err := sendPostRequest(client, url, endpoint, "application/json", body)
		if err == nil && resp.StatusCode() != http.StatusOK {
			return errors.NewUnsuccessRequestProccessing(fmt.Sprintf("Url: %s, Status code: %d",
				strings.Join([]string{url, endpoint}, ""), resp.StatusCode()), nil)
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}

	if sendErr := backoff.Retry(sendMetricsBatch, sendMetricsBatchRetryIntervals); sendErr != nil {
		logger.Error(sendErr.Error(), zap.String("event", "send update metrics batch"))
	}
}
