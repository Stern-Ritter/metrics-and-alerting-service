package storage

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	er "github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

type AgentCache interface {
	UpdateGaugeMetric(metric metrics.GaugeMetric) (metrics.GaugeMetric, error)
	UpdateCounterMetric(metric metrics.CounterMetric) (metrics.CounterMetric, error)
	UpdateMonitorMetrics(m *monitors.Monitor)
	ResetMetricValue(metricType, metricName string) error
	GetMetrics() (map[string]metrics.GaugeMetric, map[string]metrics.CounterMetric)
}

type AgentMemCache struct {
	gaugesMu   sync.Mutex
	countersMu sync.Mutex

	gauges   map[string]metrics.GaugeMetric
	counters map[string]metrics.CounterMetric

	Logger *zap.Logger
}

func NewAgentMemCache(supportedGaugeMetrics map[string]metrics.GaugeMetric, supportedCounterMetrics map[string]metrics.CounterMetric,
	logger *zap.Logger) AgentMemCache {
	return AgentMemCache{
		gauges:   supportedGaugeMetrics,
		counters: supportedCounterMetrics,
		Logger:   logger,
	}
}

func (c *AgentMemCache) UpdateGaugeMetric(metric metrics.GaugeMetric) (metrics.GaugeMetric, error) {
	c.gaugesMu.Lock()
	defer c.gaugesMu.Unlock()

	err := c.checkGaugeMetricNameWhenUpdate(metric.Name)
	if err != nil {
		return metrics.GaugeMetric{}, err
	}

	if savedMetric, exists := c.gauges[metric.Name]; exists {
		savedMetric.SetValue(metric.GetValue())
		c.gauges[metric.Name] = savedMetric
	} else {
		c.gauges[metric.Name] = metric
	}

	return c.gauges[metric.Name], nil
}

func (c *AgentMemCache) UpdateCounterMetric(metric metrics.CounterMetric) (metrics.CounterMetric, error) {
	c.countersMu.Lock()
	defer c.countersMu.Unlock()

	err := c.checkCounterMetricNameWhenUpdate(metric.Name)
	if err != nil {
		return metrics.CounterMetric{}, err
	}

	if savedMetric, exists := c.counters[metric.Name]; exists {
		savedMetric.SetValue(metric.GetValue())
		c.counters[metric.Name] = savedMetric
	} else {
		c.counters[metric.Name] = metric
	}

	return c.counters[metric.Name], nil
}

func (c *AgentMemCache) ResetMetricValue(metricType, metricName string) error {
	c.countersMu.Lock()
	c.gaugesMu.Lock()
	defer c.countersMu.Unlock()
	defer c.gaugesMu.Unlock()

	switch metrics.MetricType(metricType) {
	case metrics.Gauge:
		err := c.checkGaugeMetricNameWhenReset(metricName)
		if err != nil {
			return err
		}

		savedMetric := c.gauges[metricName]
		savedMetric.SetValue(0)
		c.gauges[savedMetric.Name] = savedMetric
	case metrics.Counter:
		err := c.checkCounterMetricNameWhenReset(metricName)
		if err != nil {
			return err
		}

		savedMetric := c.counters[metricName]
		savedMetric.ClearValue()
		c.counters[savedMetric.Name] = savedMetric

	default:
		return er.NewInvalidMetricType(fmt.Sprintf("Invalid metric type: %s", metricType), nil)
	}
	return nil
}

func (c *AgentMemCache) GetMetrics() (map[string]metrics.GaugeMetric, map[string]metrics.CounterMetric) {
	c.gaugesMu.Lock()
	gauges := utils.CopyMap(c.gauges)
	c.gaugesMu.Unlock()

	c.countersMu.Lock()
	counters := utils.CopyMap(c.counters)
	c.countersMu.Unlock()

	return gauges, counters
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
		c.Logger.Error(err.Error(), zap.String("event", "update monitor metric"))
		return
	}

	_, err = c.UpdateCounterMetric(metrics.NewCounter("PollCount", 1))
	if err != nil {
		c.Logger.Error(err.Error(), zap.String("event", "update PollCount counter metric"))
	}
}

func (c *AgentMemCache) checkGaugeMetricNameWhenUpdate(name string) error {
	_, exists := c.gauges[name]
	if !exists {
		return er.NewInvalidMetricName(fmt.Sprintf("Invalid metric name: %s", name), nil)
	}
	return nil
}

func (c *AgentMemCache) checkCounterMetricNameWhenUpdate(name string) error {
	_, exists := c.counters[name]
	if !exists {
		return er.NewInvalidMetricName(fmt.Sprintf("Invalid metric name: %s", name), nil)
	}
	return nil
}

func (c *AgentMemCache) checkGaugeMetricNameWhenReset(name string) error {
	_, exists := c.gauges[name]
	if !exists {
		return er.NewInvalidMetricName(fmt.Sprintf("Gauge metric with name: %s not exists", name), nil)
	}

	return nil
}

func (c *AgentMemCache) checkCounterMetricNameWhenReset(name string) error {
	_, exists := c.counters[name]
	if !exists {
		return er.NewInvalidMetricName(fmt.Sprintf("Counter metric with name: %s not exists", name), nil)
	}

	return nil
}
