package agent

import (
	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type Agent struct {
	HTTPClient *resty.Client
	Cache      storage.AgentCache
	Monitor    *monitors.Monitor
	Random     *utils.Random
	Config     *config.AgentConfig
	Logger     *zap.Logger
}

func NewAgent(httpClient *resty.Client, cache storage.AgentCache, monitor *monitors.Monitor,
	random *utils.Random, config *config.AgentConfig, logger *zap.Logger) *Agent {
	return &Agent{httpClient, cache, monitor, random, config, logger}
}
