package storage

import (
	"fmt"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	"go.uber.org/zap"
)

type AgentCache interface {
	Storage
	UpdateMonitorMetrics(m *monitors.Monitor)
}

type AgentMemCache struct {
	MemStorage
}

func NewAgentMemCache(supportedGaugeMetrics map[string]metrics.GaugeMetric,
	supportedCounterMetrics map[string]metrics.CounterMetric) AgentMemCache {
	return AgentMemCache{
		MemStorage: MemStorage{
			gauges:   supportedGaugeMetrics,
			counters: supportedCounterMetrics,
		},
	}
}

func (c *AgentMemCache) UpdateMonitorMetrics(m *monitors.Monitor) {
	c.updateMonitorMetric(metrics.NewGauge("Alloc", m.Alloc))
	c.updateMonitorMetric(metrics.NewGauge("BuckHashSys", m.BuckHashSys))
	c.updateMonitorMetric(metrics.NewGauge("Frees", m.Frees))
	c.updateMonitorMetric(metrics.NewGauge("GCCPUFraction", m.GCCPUFraction))
	c.updateMonitorMetric(metrics.NewGauge("GCSys", m.GCSys))
	c.updateMonitorMetric(metrics.NewGauge("HeapAlloc", m.HeapAlloc))
	c.updateMonitorMetric(metrics.NewGauge("HeapIdle", m.HeapIdle))
	c.updateMonitorMetric(metrics.NewGauge("HeapInuse", m.HeapInuse))
	c.updateMonitorMetric(metrics.NewGauge("HeapObjects", m.HeapObjects))
	c.updateMonitorMetric(metrics.NewGauge("HeapReleased", m.HeapReleased))
	c.updateMonitorMetric(metrics.NewGauge("HeapSys", m.HeapSys))
	c.updateMonitorMetric(metrics.NewGauge("LastGC", m.LastGC))
	c.updateMonitorMetric(metrics.NewGauge("Lookups", m.Lookups))
	c.updateMonitorMetric(metrics.NewGauge("MCacheInuse", m.MCacheInuse))
	c.updateMonitorMetric(metrics.NewGauge("MCacheSys", m.MCacheSys))
	c.updateMonitorMetric(metrics.NewGauge("MSpanInuse", m.MSpanInuse))
	c.updateMonitorMetric(metrics.NewGauge("MSpanSys", m.MSpanSys))
	c.updateMonitorMetric(metrics.NewGauge("Mallocs", m.Mallocs))
	c.updateMonitorMetric(metrics.NewGauge("NextGC", m.NextGC))
	c.updateMonitorMetric(metrics.NewGauge("NumForcedGC", m.NumForcedGC))
	c.updateMonitorMetric(metrics.NewGauge("NumGC", m.NumGC))
	c.updateMonitorMetric(metrics.NewGauge("OtherSys", m.OtherSys))
	c.updateMonitorMetric(metrics.NewGauge("PauseTotalNs", m.PauseTotalNs))
	c.updateMonitorMetric(metrics.NewGauge("StackInuse", m.StackInuse))
	c.updateMonitorMetric(metrics.NewGauge("StackSys", m.StackSys))
	c.updateMonitorMetric(metrics.NewGauge("Sys", m.Sys))
	c.updateMonitorMetric(metrics.NewGauge("TotalAlloc", m.TotalAlloc))
}

func (c *AgentMemCache) updateMonitorMetric(metric metrics.GaugeMetric) {
	_, err := c.UpdateGaugeMetric(metric)
	if err != nil {
		logger.Log.Error(err.Error(), zap.String("event", "update monitor metric"))
		return
	}

	_, err = c.UpdateCounterMetric(metrics.NewCounter("PollCount", 1))
	if err != nil {
		logger.Log.Error(err.Error(), zap.String("event", "update PollCount counter metric"))
	}
}

func (c *AgentMemCache) CheckGaugeMetricNameWhenUpdate(name string) error {
	_, exists := c.gauges[name]
	if !exists {
		return errors.NewInvalidMetricName(fmt.Sprintf("Invalid metric name: %s", name), nil)
	}
	return nil
}

func (c *AgentMemCache) CheckCounterMetricNameWhenUpdate(name string) error {
	_, exists := c.counters[name]
	if !exists {
		return errors.NewInvalidMetricName(fmt.Sprintf("Invalid metric name: %s", name), nil)
	}
	return nil
}
