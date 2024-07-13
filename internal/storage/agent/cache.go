package storage

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	er "github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

// AgentCache defines an interface for getting and updating metrics statistics in the agent's cache.
type AgentCache interface {
	UpdateGaugeMetric(metric metrics.GaugeMetric) (metrics.GaugeMetric, error)
	UpdateCounterMetric(metric metrics.CounterMetric) (metrics.CounterMetric, error)
	UpdateRuntimeMonitorMetrics(m *monitors.RuntimeMonitor)
	UpdateUtilMonitorMetrics(m *monitors.UtilMonitor)
	ResetMetricValue(metricType, metricName string) error
	GetMetrics() (map[string]metrics.GaugeMetric, map[string]metrics.CounterMetric)
}

// AgentMemCache is an in-memory implementation of the AgentCache interface.
type AgentMemCache struct {
	gaugesMu   sync.RWMutex
	countersMu sync.RWMutex

	gauges   map[string]metrics.GaugeMetric
	counters map[string]metrics.CounterMetric

	Logger *logger.AgentLogger
}

// NewAgentMemCache is constructor for creating a new AgentMemCache with the provided supported gauge and counter metrics.
func NewAgentMemCache(supportedGaugeMetrics map[string]metrics.GaugeMetric, supportedCounterMetrics map[string]metrics.CounterMetric,
	logger *logger.AgentLogger) AgentMemCache {
	return AgentMemCache{
		gauges:   supportedGaugeMetrics,
		counters: supportedCounterMetrics,
		Logger:   logger,
	}
}

// UpdateGaugeMetric updates or adds a gauge metric in the cache and returns the updated metric.
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

// UpdateCounterMetric updates or adds a counter metric in the cache and returns the updated metric.
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

// ResetMetricValue resets the value of the specified metric to 0.
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

// GetMetrics returns all the gauge and counter metrics from the cache.
func (c *AgentMemCache) GetMetrics() (map[string]metrics.GaugeMetric, map[string]metrics.CounterMetric) {
	c.gaugesMu.RLock()
	gauges := utils.CopyMap(c.gauges)
	c.gaugesMu.RUnlock()

	c.countersMu.RLock()
	counters := utils.CopyMap(c.counters)
	c.countersMu.RUnlock()

	return gauges, counters
}

// UpdateRuntimeMonitorMetrics updates the runtime monitor metrics in the cache.
func (c *AgentMemCache) UpdateRuntimeMonitorMetrics(m *monitors.RuntimeMonitor) {
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

// UpdateUtilMonitorMetrics updates the utilization monitor metrics in the cache.
func (c *AgentMemCache) UpdateUtilMonitorMetrics(m *monitors.UtilMonitor) {
	c.updateMonitorMetric(metrics.NewGauge("TotalMemory", m.TotalMemory))
	c.updateMonitorMetric(metrics.NewGauge("FreeMemory", m.FreeMemory))
	c.updateMonitorMetric(metrics.NewGauge("CPUutilization1", m.CPUutilization1))
}

func (c *AgentMemCache) updateMonitorMetric(metric metrics.GaugeMetric) {
	_, err := c.UpdateGaugeMetric(metric)
	if err != nil {
		c.Logger.Error(err.Error(), zap.String("event", "update monitor metric"),
			zap.String("metric type", string(metric.Type)), zap.String("metric name", metric.Name))
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
