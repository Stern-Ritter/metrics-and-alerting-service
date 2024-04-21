package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"

	er "github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
)

func (a *Agent) UpdateMetrics() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	a.Monitor.Update(&ms)
	randomValue, _ := a.Random.Float(0.1, 99.99)

	a.Cache.UpdateMonitorMetrics(a.Monitor)
	_, err := a.Cache.UpdateGaugeMetric(metrics.NewGauge("RandomValue", randomValue))
	if err != nil {
		a.Logger.Error(err.Error(), zap.String("event", "update RandomValue gauge metric"))
	}
}

func (a *Agent) SendMetrics() {
	gauges, counters := a.Cache.GetMetrics()

	err := a.Cache.ResetMetricValue(string(metrics.Counter), "PollCount")
	if err != nil {
		a.Logger.Error(err.Error(), zap.String("event", "reset PollCount counter metric"))
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
		a.Logger.Error(err.Error(), zap.String("event", "JSON encoding metrics batch"))
	}

	sendMetricsBatch := func() error {
		resp, err := sendPostRequest(a.HTTPClient, a.Config.SendMetricsURL, a.Config.SendMetricsEndPoint,
			"application/json", body)

		if err == nil && resp.StatusCode() != http.StatusOK {
			return er.NewUnsuccessRequestProccessing(fmt.Sprintf("Url: %s, Status code: %d",
				strings.Join([]string{a.Config.SendMetricsURL, a.Config.SendMetricsEndPoint}, ""),
				resp.StatusCode()), nil)
		} else if err != nil {
			return backoff.Permanent(err)
		}

		return nil
	}

	if sendErr := backoff.Retry(sendMetricsBatch, a.sendMetricsBatchRetryIntervals); sendErr != nil {
		a.Logger.Error(sendErr.Error(), zap.String("event", "send update metrics batch"))
	}
}
