package agent

import (
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/go-resty/resty/v2"
)

type Agent struct {
	HTTPClient *resty.Client
	Cache      storage.AgentCache
	Monitor    *monitors.Monitor
	Random     *utils.Random
}

func NewAgent(httpClient *resty.Client, cache storage.AgentCache, monitor *monitors.Monitor, random *utils.Random) *Agent {
	return &Agent{httpClient, cache, monitor, random}
}
