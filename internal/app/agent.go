package app

import (
	"context"
	"sync"
	"time"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/config"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	tasks "github.com/Stern-Ritter/metrics-and-alerting-service/internal/transport"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/go-resty/resty/v2"
)

const (
	taksCount = 2
)

type MonitoringAgent struct {
	httpClient *resty.Client
	cache      storage.AgentCache
	monitor    *model.Monitor
	random     *utils.Random
	config     config.AgentConfig
}

func NewMonitoringAgent(httpClient *resty.Client, cache storage.AgentCache, monitor *model.Monitor, random *utils.Random,
	config config.AgentConfig) MonitoringAgent {
	return MonitoringAgent{httpClient, cache, monitor, random, config}
}

func (m *MonitoringAgent) Run() {
	var wg sync.WaitGroup
	wg.Add(taksCount)

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Hour, cancel)

	updateMetricsTask := func() { tasks.UpdateMetrics(m.cache, m.monitor, m.random) }
	sendMetricsTask := func() {
		tasks.SendMetrics(m.httpClient, m.config.SendMetricsURL.String(), m.config.SendMetricsEndPoint, m.cache)
	}

	tasks.SetInterval(ctx, &wg, updateMetricsTask, m.config.UpdateMetricsInterval)
	tasks.SetInterval(ctx, &wg, sendMetricsTask, m.config.SendMetricsInterval)

	wg.Wait()
}
