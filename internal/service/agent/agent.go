package agent

import (
	"crypto/rsa"
	"time"

	"github.com/cenkalti/backoff/v4"
	"gopkg.in/h2non/gentleman.v2"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/agent"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	cache "github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	pb "github.com/Stern-Ritter/metrics-and-alerting-service/proto/gen/metrics"
)

// Agent is monitoring agent that collects and sends metrics statistics to the server.
type Agent struct {
	HTTPClient                     *gentleman.Client
	GRPCClient                     pb.MetricsClient
	Cache                          cache.AgentCache
	RuntimeMonitor                 *monitors.RuntimeMonitor
	UtilMonitor                    *monitors.UtilMonitor
	Random                         *utils.Random
	Config                         *config.AgentConfig
	metricsCh                      chan []metrics.Metrics
	doneCh                         chan struct{}
	sendMetricsBatchRetryIntervals *backoff.ExponentialBackOff
	rsaPublicKey                   *rsa.PublicKey
	Logger                         *logger.AgentLogger
}

// NewAgent is constructor for creating a new Agent.
func NewAgent(cache cache.AgentCache, runtimeMonitor *monitors.RuntimeMonitor,
	utilMonitor *monitors.UtilMonitor, random *utils.Random, config *config.AgentConfig, rsaPublicKey *rsa.PublicKey,
	logger *logger.AgentLogger) *Agent {

	sendMetricsBatchRetryIntervals := backoff.NewExponentialBackOff(
		backoff.WithInitialInterval(1*time.Second),
		backoff.WithRandomizationFactor(0),
		backoff.WithMultiplier(3),
		backoff.WithMaxInterval(5*time.Second),
		backoff.WithMaxElapsedTime(10*time.Second))

	metricsCh := make(chan []metrics.Metrics, config.MetricsBufferSize)
	doneCh := make(chan struct{})

	return &Agent{
		Cache:                          cache,
		RuntimeMonitor:                 runtimeMonitor,
		UtilMonitor:                    utilMonitor,
		Random:                         random,
		Config:                         config,
		metricsCh:                      metricsCh,
		doneCh:                         doneCh,
		sendMetricsBatchRetryIntervals: sendMetricsBatchRetryIntervals,
		rsaPublicKey:                   rsaPublicKey,
		Logger:                         logger}
}

// SetHTTPClient sets the HTTP client for the Agent.
func (a *Agent) SetHTTPClient(client *gentleman.Client) {
	a.HTTPClient = client
}

// SetGRPCClient sets the gRPC client for the Agent.
func (a *Agent) SetGRPCClient(grpcClient pb.MetricsClient) {
	a.GRPCClient = grpcClient
}
