package agent

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/cenkalti/backoff/v4"
	"github.com/shirou/gopsutil/v3/mem"
	"go.uber.org/zap"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
)

// UpdateRuntimeMetrics task that collects runtime metrics statistics and updates the cache.
func (a *Agent) UpdateRuntimeMetrics() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	a.RuntimeMonitor.Update(&ms)
	a.Cache.UpdateRuntimeMonitorMetrics(a.RuntimeMonitor)

	randomValue, _ := a.Random.Float(0.1, 99.99)
	_, err := a.Cache.UpdateGaugeMetric(metrics.NewGauge("RandomValue", randomValue))
	if err != nil {
		a.Logger.Error(err.Error(), zap.String("event", "update RandomValue gauge metric"))
	}
}

// UpdateUtilMetrics task that collects utilization metrics and updates the cache
func (a *Agent) UpdateUtilMetrics() {
	ms, err := mem.VirtualMemory()
	if err != nil {
		a.Logger.Error(err.Error(), zap.String("event", "update util gauge metrics"))
		return
	}
	err = a.UtilMonitor.Update(ms)
	if err != nil {
		a.Logger.Error(err.Error(), zap.String("event", "update util gauge metrics"))
		return
	}
	a.Cache.UpdateUtilMonitorMetrics(a.UtilMonitor)
}

// SendMetrics task that gets all metrics from the cache, resets the PollCount counter metric
// and sends the metrics statistics to the server.
func (a *Agent) SendMetrics() {
	gauges, counters := a.Cache.GetMetrics()

	err := a.Cache.ResetMetricValue(string(metrics.Counter), "PollCount")
	if err != nil {
		a.Logger.Error(err.Error(), zap.String("event", "reset PollCount counter metric"))
	}

	metricsBatch := make([]metrics.Metrics, 0, len(gauges)+len(counters))
	for _, gaugeMetric := range gauges {
		metricsBatch = append(metricsBatch, metrics.GaugeMetricToMetrics(gaugeMetric))

	}
	for _, counterMetric := range counters {
		metricsBatch = append(metricsBatch, metrics.CounterMetricToMetrics(counterMetric))
	}

	select {
	case <-a.doneCh:
		close(a.metricsCh)
		a.Logger.Info("send metrics task stopped")
		return
	case a.metricsCh <- metricsBatch:
	}
}

// StartSendMetricsWorkerPool starts a pool of workers to send metrics statistics.
func (a *Agent) StartSendMetricsWorkerPool() {
	sendRateLimit := a.Config.RateLimit
	if sendRateLimit <= 0 {
		a.Logger.Error("Rate limit can't be less than or equal to zero",
			zap.String("event", "start send metrics worker pool"))
		sendRateLimit = 1
	}

	for w := 1; w <= sendRateLimit; w++ {
		go a.sendMetricsWorker(w, a.metricsCh)
	}

	a.Logger.Debug("Worker pool started",
		zap.String("event", "starting send metrics worker pool"))
}

func (a *Agent) sendMetricsWorker(id int, metricsCh <-chan []metrics.Metrics) {
	a.Logger.Debug("Worker started", zap.Int("worker id", id),
		zap.String("event", "starting send metrics worker"))

	for metricsBatch := range metricsCh {
		body, err := json.Marshal(metricsBatch)
		if err != nil {
			a.Logger.Error(err.Error(), zap.Int("worker id", id),
				zap.String("event", "JSON encoding metrics"))
			continue
		}

		sendMetricsBatch := func() error {
			resp, err := sendPostRequest(a.HTTPClient, a.Config.SendMetricsEndPoint, "application/json", body)
			if err == nil && !resp.Ok {
				return errors.NewUnsuccessRequestProcessing(fmt.Sprintf("unsuccess request sent on url: %s, status code: %d",
					a.Config.SendMetricsEndPoint, resp.StatusCode), nil)
			} else if err != nil {
				return backoff.Permanent(err)
			}

			return nil
		}

		if sendErr := backoff.Retry(sendMetricsBatch, a.sendMetricsBatchRetryIntervals); sendErr != nil {
			a.Logger.Error(sendErr.Error(), zap.Int("worker id", id),
				zap.String("event", "sending metrics update"))
			continue
		}

		a.Logger.Debug("Success sent metrics update", zap.Int("worker id", id),
			zap.String("event", "sending metrics update"))
	}
	a.Logger.Debug("Worker stopped", zap.Int("worker id", id),
		zap.String("event", "stopping send metrics worker"))
}

// StopTasks stops all Agent tasks
func (a *Agent) StopTasks() {
	close(a.doneCh)
}
