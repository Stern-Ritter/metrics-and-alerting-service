package storage

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

type MemStorage struct {
	gaugesMu   sync.Mutex
	gauges     map[string]metrics.GaugeMetric
	countersMu sync.Mutex
	counters   map[string]metrics.CounterMetric
}

func (s *MemStorage) UpdateGaugeMetric(metric metrics.GaugeMetric) (metrics.GaugeMetric, error) {
	s.gaugesMu.Lock()
	defer s.gaugesMu.Unlock()

	err := s.checkGaugeMetricNameWhenUpdate(metric.Name)
	if err != nil {
		return metrics.GaugeMetric{}, err
	}

	if savedMetric, exists := s.gauges[metric.Name]; exists {
		savedMetric.SetValue(metric.GetValue())
		s.gauges[metric.Name] = savedMetric
	} else {
		s.gauges[metric.Name] = metric
	}

	return s.gauges[metric.Name], nil
}

func (s *MemStorage) checkGaugeMetricNameWhenUpdate(name string) error {
	return nil
}

func (s *MemStorage) UpdateCounterMetric(metric metrics.CounterMetric) (metrics.CounterMetric, error) {
	s.countersMu.Lock()
	defer s.countersMu.Unlock()

	err := s.checkCounterMetricNameWhenUpdate(metric.Name)
	if err != nil {
		return metrics.CounterMetric{}, nil
	}

	if savedMetric, exists := s.counters[metric.Name]; exists {
		savedMetric.SetValue(metric.GetValue())
		s.counters[metric.Name] = savedMetric
	} else {
		s.counters[metric.Name] = metric
	}

	return s.counters[metric.Name], nil
}

func (s *MemStorage) checkCounterMetricNameWhenUpdate(name string) error {
	return nil
}

func (s *MemStorage) UpdateMetric(metricType, metricName, metricValue string) error {
	switch metrics.MetricType(metricType) {
	case metrics.Gauge:
		value, err := parseGaugeMetricValue(metricValue)
		if err != nil {
			return err
		}
		metric := metrics.NewGauge(metricName, value)
		_, err = s.UpdateGaugeMetric(metric)
		if err != nil {
			return err
		}
	case metrics.Counter:
		value, err := parseCounterMetricValue(metricValue)
		if err != nil {
			return err
		}
		metric := metrics.NewCounter(metricName, value)
		_, err = s.UpdateCounterMetric(metric)
		if err != nil {
			return err
		}

	default:
		return errors.NewInvalidMetricType(fmt.Sprintf("Invalid metric type: %s", metricType), nil)
	}
	return nil
}

func parseGaugeMetricValue(v string) (float64, error) {
	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, errors.NewInvalidMetricValue(
			fmt.Sprintf("The value for the %s metric should be of float64 type", metrics.Gauge), err)
	}

	return value, nil
}

func parseCounterMetricValue(v string) (int64, error) {
	value, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, errors.NewInvalidMetricValue(
			fmt.Sprintf("The value for the %s metric should be of int64 type", metrics.Counter), err)
	}

	return value, nil
}

func (s *MemStorage) ResetMetricValue(metricType, metricName string) error {
	switch metrics.MetricType(metricType) {
	case metrics.Gauge:
		err := s.CheckGaugeMetricNameWhenReset(metricName)
		if err != nil {
			return err
		}

		savedMetric := s.gauges[metricName]
		savedMetric.SetValue(0)
		s.gauges[savedMetric.Name] = savedMetric
	case metrics.Counter:
		err := s.CheckCounterMetricNameWhenReset(metricName)
		if err != nil {
			return err
		}

		savedMetric := s.counters[metricName]
		savedMetric.ClearValue()
		s.counters[savedMetric.Name] = savedMetric

	default:
		return errors.NewInvalidMetricType(fmt.Sprintf("Invalid metric type: %s", metricType), nil)
	}
	return nil
}

func (s *MemStorage) CheckGaugeMetricNameWhenReset(name string) error {
	_, exists := s.gauges[name]
	if !exists {
		return errors.NewInvalidMetricName(fmt.Sprintf("Gauge metric with name: %s not exists", name), nil)
	}

	return nil
}

func (s *MemStorage) CheckCounterMetricNameWhenReset(name string) error {
	_, exists := s.counters[name]
	if !exists {
		return errors.NewInvalidMetricName(fmt.Sprintf("Counter metric with name: %s not exists", name), nil)
	}

	return nil
}

func (s *MemStorage) GetMetrics() (map[string]metrics.GaugeMetric, map[string]metrics.CounterMetric) {
	s.gaugesMu.Lock()
	gauges := utils.CopyMap(s.gauges)
	s.gaugesMu.Unlock()

	s.countersMu.Lock()
	counters := utils.CopyMap(s.counters)
	s.countersMu.Unlock()

	return gauges, counters
}

func (s *MemStorage) GetMetricValueByTypeAndName(metricType, metricName string) (string, error) {
	var value string
	var err error

	switch metrics.MetricType(metricType) {
	case metrics.Gauge:
		s.gaugesMu.Lock()
		metric, exists := s.gauges[metricName]
		s.gaugesMu.Unlock()
		if !exists {
			err = errors.NewInvalidMetricName(fmt.Sprintf("Gauge metric with name: %s not exists", metricName), nil)
			break
		}
		value = utils.FormatGaugeMetricValue(metric.GetValue())
	case metrics.Counter:
		s.countersMu.Lock()
		metric, exists := s.counters[metricName]
		s.countersMu.Unlock()
		if !exists {
			err = errors.NewInvalidMetricName(fmt.Sprintf("Counter metric with name: %s not exists", metricName), nil)
			break
		}
		value = utils.FormatCounterMetricValue(metric.GetValue())

	default:
		err = errors.NewInvalidMetricType(fmt.Sprintf("Invalid metric type: %s", metricType), nil)
	}

	return value, err
}

func (s *MemStorage) GetGaugeMetric(metricName string) (metrics.GaugeMetric, error) {
	s.gaugesMu.Lock()
	metric, exists := s.gauges[metricName]
	s.gaugesMu.Unlock()

	if !exists {
		return metrics.GaugeMetric{}, errors.NewInvalidMetricName(fmt.Sprintf("Gauge metric with name: %s not exists", metricName), nil)
	}

	return metric, nil
}

func (s *MemStorage) GetCounterMetric(metricName string) (metrics.CounterMetric, error) {
	s.countersMu.Lock()
	metric, exists := s.counters[metricName]
	s.countersMu.Unlock()

	if !exists {
		return metrics.CounterMetric{}, errors.NewInvalidMetricName(fmt.Sprintf("Counter metric with name: %s not exists", metricName), nil)
	}

	return metric, nil
}

func (s *MemStorage) SetGaugeMetircs(gauges map[string]metrics.GaugeMetric) {
	s.gaugesMu.Lock()
	defer s.gaugesMu.Unlock()
	s.gauges = gauges
}

func (s *MemStorage) SetCounterMetrics(counters map[string]metrics.CounterMetric) {
	s.countersMu.Lock()
	defer s.countersMu.Unlock()
	s.counters = counters
}
