package storage

import (
	"fmt"
	"strconv"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/metrics"
)

type MemStorage struct {
	gauges   map[string]metrics.Gauge
	counters map[string]metrics.Counter
}

func NewMemStore() MemStorage {
	return MemStorage{
		gauges:   make(map[string]metrics.Gauge),
		counters: make(map[string]metrics.Counter),
	}
}

func (s *MemStorage) UpdateMetric(metricType, metricName, metricValue string) error {
	switch metricType {
	case metrics.TypeGauge:
		value, err := parseGaugeTypeMetricValue(metricValue)
		if err != nil {
			return err
		}
		updateGaugeTypeMetric(s.gauges, metricName, value)
	case metrics.TypeCounter:
		value, err := parseCounterTypeMetricValue(metricValue)
		if err != nil {
			return err
		}
		updateCounterTypeMetric(s.counters, metricName, value)
	default:
		return errors.NewInvalidMetricValue(fmt.Sprintf("Invalid metric type: %s\n", metricType), nil)
	}
	return nil
}

func parseGaugeTypeMetricValue(metricValue string) (float64, error) {
	value, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		return 0, errors.NewInvalidMetricValue(
			fmt.Sprintf("The value for the %s metric should be of float64 type", metrics.TypeGauge),
			err)
	}

	return value, nil
}

func updateGaugeTypeMetric(gauges map[string]metrics.Gauge, metricName string, value float64) {
	if metric, exists := gauges[metricName]; exists {
		metric.SetValue(value)
		gauges[metricName] = metric
	} else {
		newMetric := metrics.NewGauge(value)
		gauges[metricName] = newMetric
	}
}

func parseCounterTypeMetricValue(metricValue string) (int64, error) {
	value, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		return 0, errors.NewInvalidMetricValue(
			fmt.Sprintf("The value for the %s metric should be of int64 type", metrics.TypeCounter),
			err)
	}

	return value, nil
}

func updateCounterTypeMetric(counters map[string]metrics.Counter, metricName string, value int64) {
	if metric, exists := counters[metricName]; exists {
		metric.SetValue(value)
		counters[metricName] = metric
	} else {
		newMetric := metrics.NewCounter(value)
		counters[metricName] = newMetric
	}
}
