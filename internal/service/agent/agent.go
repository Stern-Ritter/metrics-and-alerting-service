package agent

import (
	"time"

	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"

	"gopkg.in/h2non/gentleman.v2"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	cache "github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

// Agent is monitoring agent that collects and sends metrics statistics to the server.
type Agent struct {
	HTTPClient                     *gentleman.Client
	Cache                          cache.AgentCache
	RuntimeMonitor                 *monitors.RuntimeMonitor
	UtilMonitor                    *monitors.UtilMonitor
	Random                         *utils.Random
	Config                         *config.AgentConfig
	metricsCh                      chan []metrics.Metrics
	doneCh                         chan struct{}
	sendMetricsBatchRetryIntervals *backoff.ExponentialBackOff
	Logger                         *zap.Logger
}

// NewAgent is constructor for creating a new Agent.
func NewAgent(httpClient *gentleman.Client, cache cache.AgentCache, runtimeMonitor *monitors.RuntimeMonitor,
	utilMonitor *monitors.UtilMonitor, random *utils.Random, config *config.AgentConfig, logger *zap.Logger) *Agent {

	sendMetricsBatchRetryIntervals := backoff.NewExponentialBackOff(
		backoff.WithInitialInterval(1*time.Second),
		backoff.WithRandomizationFactor(0),
		backoff.WithMultiplier(3),
		backoff.WithMaxInterval(5*time.Second),
		backoff.WithMaxElapsedTime(10*time.Second))

	metricsCh := make(chan []metrics.Metrics, config.MetricsBufferSize)
	doneCh := make(chan struct{})

	return &Agent{
		HTTPClient:                     httpClient,
		Cache:                          cache,
		RuntimeMonitor:                 runtimeMonitor,
		UtilMonitor:                    utilMonitor,
		Random:                         random,
		Config:                         config,
		metricsCh:                      metricsCh,
		doneCh:                         doneCh,
		sendMetricsBatchRetryIntervals: sendMetricsBatchRetryIntervals,
		Logger:                         logger}
}
