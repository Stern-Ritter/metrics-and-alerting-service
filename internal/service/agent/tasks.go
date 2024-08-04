package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"

	"github.com/cenkalti/backoff/v4"
	"github.com/shirou/gopsutil/v3/mem"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	pb "github.com/Stern-Ritter/metrics-and-alerting-service/proto/gen/metrics"
)

const (
	ipKey = "X-Real-IP"
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

	a.metricsCh <- metricsBatch
}

// StartSendMetricsWorkerPool starts a pool of workers to send metrics statistics.
func (a *Agent) StartSendMetricsWorkerPool(wg *sync.WaitGroup,
	worker func(id int, metricsCh <-chan []metrics.Metrics, wg *sync.WaitGroup)) {
	sendRateLimit := a.Config.RateLimit
	if sendRateLimit <= 0 {
		a.Logger.Error("Rate limit can't be less than or equal to zero",
			zap.String("event", "start send metrics worker pool"))
		sendRateLimit = 1
	}

	for w := 1; w <= sendRateLimit; w++ {
		go worker(w, a.metricsCh, wg)
		wg.Add(1)
	}

	a.Logger.Debug("Worker pool started",
		zap.String("event", "starting send metrics worker pool"))
}

// SendMetricsWithHTTPWorker is a worker function that sends metrics using HTTP.
func (a *Agent) SendMetricsWithHTTPWorker(id int, metricsCh <-chan []metrics.Metrics, wg *sync.WaitGroup) {
	a.Logger.Debug("Worker started", zap.Int("worker id", id),
		zap.String("event", "starting send metrics worker"))

	for metricsBatch := range metricsCh {
		body, err := json.Marshal(metricsBatch)
		if err != nil {
			a.Logger.Error(err.Error(), zap.Int("worker id", id),
				zap.String("event", "JSON encoding metrics"))
			continue
		}

		ipAddr, err := getIPAddr()
		if err != nil {
			a.Logger.Error(err.Error(), zap.Int("worker id", id),
				zap.String("event", "get ip address"))
		}
		headers := map[string]string{
			"Content-Type": "application/json",
			ipKey:          ipAddr,
		}
		sendMetricsBatch := func() error {
			resp, err := sendPostRequest(a.HTTPClient, a.Config.SendMetricsEndPoint, headers, body)
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
	wg.Done()
}

// SendMetricsWithGrpcWorker is a worker function that sends metrics using gRPC.
func (a *Agent) SendMetricsWithGrpcWorker(id int, metricsCh <-chan []metrics.Metrics, wg *sync.WaitGroup) {
	a.Logger.Debug("Worker started", zap.Int("worker id", id),
		zap.String("event", "starting send metrics worker"))

	for metricsBatch := range metricsCh {
		metricsData := metrics.MetricsToRepeatedMetricData(metricsBatch)

		ipAddr, err := getIPAddr()
		if err != nil {
			a.Logger.Error(err.Error(), zap.Int("worker id", id),
				zap.String("event", "get ip address"))
		}
		md := metadata.Pairs(ipKey, ipAddr)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		sendMetricData := func() error {
			updateMetricsBatchRequest := &pb.UpdateMetricsBatchRequest{
				Metrics: metricsData,
			}
			_, err := a.GRPCClient.UpdateMetricsBatch(ctx, updateMetricsBatchRequest)
			if err != nil {
				if e, ok := status.FromError(err); ok {
					if e.Code() == codes.Unavailable || e.Code() == codes.DeadlineExceeded {
						return errors.NewUnsuccessRequestProcessing(
							fmt.Sprintf("unsuccess request sent on url: %s, status code: %d",
								a.Config.SendMetricsEndPoint, e.Code()), nil)
					} else {
						return backoff.Permanent(errors.NewUnsuccessRequestProcessing(
							fmt.Sprintf("unsuccess request sent on url: %s, status code: %d",
								a.Config.SendMetricsEndPoint, e.Code()), nil))
					}
				}
				return errors.NewUnsuccessRequestProcessing(
					fmt.Sprintf("unsuccess request sent on url: %s, error parsing status code: %s",
						a.Config.SendMetricsEndPoint, err), nil)
			}

			return nil
		}

		if sendErr := backoff.Retry(sendMetricData, a.sendMetricsBatchRetryIntervals); sendErr != nil {
			a.Logger.Error(sendErr.Error(), zap.Int("worker id", id),
				zap.String("event", "sending metrics update"))
			continue
		}

		a.Logger.Debug("Success sent metrics update", zap.Int("worker id", id),
			zap.String("event", "sending metrics update"))
	}

	a.Logger.Debug("Worker stopped", zap.Int("worker id", id),
		zap.String("event", "stopping send metrics worker"))
	wg.Done()
}

// StopSendMetricsWorkerPool stops send metrics workers
func (a *Agent) StopSendMetricsWorkerPool() {
	a.Logger.Debug("Worker pool stopped",
		zap.String("event", "starting send metrics worker pool"))
	close(a.metricsCh)
}
