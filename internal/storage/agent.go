package storage

import (
	"fmt"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model"
)

type AgentCache interface {
	Storage
	UpdateMonitorMetrics(m *model.Monitor)
}

type AgentMemCache struct {
	MemStorage
}

func NewAgentMemCache(supportedGaugeMetrics map[string]model.GaugeMetric, supportedCounterMetrics map[string]model.CounterMetric) AgentMemCache {
	return AgentMemCache{
		MemStorage: MemStorage{
			gauges:   supportedGaugeMetrics,
			counters: supportedCounterMetrics,
		},
	}
}

func (c *AgentMemCache) UpdateMonitorMetrics(m *model.Monitor) {
	c.updateMonitorMetric(model.NewGauge("Alloc", m.Alloc))
	c.updateMonitorMetric(model.NewGauge("BuckHashSys", m.BuckHashSys))
	c.updateMonitorMetric(model.NewGauge("Frees", m.Frees))
	c.updateMonitorMetric(model.NewGauge("GCCPUFraction", m.GCCPUFraction))
	c.updateMonitorMetric(model.NewGauge("GCSys", m.GCSys))
	c.updateMonitorMetric(model.NewGauge("HeapAlloc", m.HeapAlloc))
	c.updateMonitorMetric(model.NewGauge("HeapIdle", m.HeapIdle))
	c.updateMonitorMetric(model.NewGauge("HeapInuse", m.HeapInuse))
	c.updateMonitorMetric(model.NewGauge("HeapObjects", m.HeapObjects))
	c.updateMonitorMetric(model.NewGauge("HeapReleased", m.HeapReleased))
	c.updateMonitorMetric(model.NewGauge("HeapSys", m.HeapSys))
	c.updateMonitorMetric(model.NewGauge("LastGC", m.LastGC))
	c.updateMonitorMetric(model.NewGauge("Lookups", m.Lookups))
	c.updateMonitorMetric(model.NewGauge("MCacheInuse", m.MCacheInuse))
	c.updateMonitorMetric(model.NewGauge("MCacheSys", m.MCacheSys))
	c.updateMonitorMetric(model.NewGauge("MSpanInuse", m.MSpanInuse))
	c.updateMonitorMetric(model.NewGauge("MSpanSys", m.MSpanSys))
	c.updateMonitorMetric(model.NewGauge("Mallocs", m.Mallocs))
	c.updateMonitorMetric(model.NewGauge("NextGC", m.NextGC))
	c.updateMonitorMetric(model.NewGauge("NumForcedGC", m.NumForcedGC))
	c.updateMonitorMetric(model.NewGauge("NumGC", m.NumGC))
	c.updateMonitorMetric(model.NewGauge("OtherSys", m.OtherSys))
	c.updateMonitorMetric(model.NewGauge("PauseTotalNs", m.PauseTotalNs))
	c.updateMonitorMetric(model.NewGauge("StackInuse", m.StackInuse))
	c.updateMonitorMetric(model.NewGauge("StackSys", m.StackSys))
	c.updateMonitorMetric(model.NewGauge("Sys", m.Sys))
	c.updateMonitorMetric(model.NewGauge("TotalAlloc", m.TotalAlloc))
}

func (c *AgentMemCache) updateMonitorMetric(metric model.GaugeMetric) {
	err := c.UpdateGaugeMetric(metric)
	if err == nil {
		err := c.UpdateCounterMetric(model.NewCounter("PollCount", 1))
		if err != nil {
			fmt.Println(err)
		}
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
