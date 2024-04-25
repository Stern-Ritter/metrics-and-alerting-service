package agent

import (
	"time"

	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"

	"gopkg.in/h2non/gentleman.v2"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	cache "github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

type Agent struct {
	HTTPClient                     *gentleman.Client
	Cache                          cache.AgentCache
	Monitor                        *monitors.Monitor
	Random                         *utils.Random
	Config                         *config.AgentConfig
	Logger                         *zap.Logger
	sendMetricsBatchRetryIntervals *backoff.ExponentialBackOff
}

func NewAgent(httpClient *gentleman.Client, cache cache.AgentCache, monitor *monitors.Monitor,
	random *utils.Random, config *config.AgentConfig, logger *zap.Logger) *Agent {
	sendMetricsBatchRetryIntervals := backoff.NewExponentialBackOff(
		backoff.WithInitialInterval(1*time.Second),
		backoff.WithRandomizationFactor(0),
		backoff.WithMultiplier(3),
		backoff.WithMaxInterval(5*time.Second),
		backoff.WithMaxElapsedTime(10*time.Second))

	return &Agent{httpClient, cache, monitor, random, config, logger,
		sendMetricsBatchRetryIntervals}
}
