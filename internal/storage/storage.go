package storage

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

type Storage interface {
	UpdateGaugeMetric(metric model.GaugeMetric) error
	UpdateCounterMetric(metric model.CounterMetric) error
	UpdateMetric(metricType, metricName, metricValue string) error
	ResetMetricValue(metricType, metricName string) error
	GetMetrics() (map[string]model.GaugeMetric, map[string]model.CounterMetric)
}

type MemStorage struct {
	gaugesMu   sync.Mutex
	gauges     map[string]model.GaugeMetric
	countersMu sync.Mutex
	counters   map[string]model.CounterMetric
}

func (s *MemStorage) UpdateGaugeMetric(metric model.GaugeMetric) error {
	s.gaugesMu.Lock()
	defer s.gaugesMu.Unlock()

	err := s.checkGaugeMetricNameWhenUpdate(metric.Name)
	if err != nil {
		return err
	}

	savedMetric, exists := s.gauges[metric.Name]
	if exists {
		savedMetric.SetValue(metric.GetValue())
		s.gauges[savedMetric.Name] = savedMetric
	} else {
		s.gauges[metric.Name] = metric
	}

	return nil
}

func (s *MemStorage) checkGaugeMetricNameWhenUpdate(name string) error {
	return nil
}

func (s *MemStorage) UpdateCounterMetric(metric model.CounterMetric) error {
	s.countersMu.Lock()
	defer s.countersMu.Unlock()

	err := s.checkCounterMetricNameWhenUpdate(metric.Name)
	if err != nil {
		return err
	}

	savedMetric, exists := s.counters[metric.Name]
	if exists {
		savedMetric.SetValue(metric.GetValue())
		s.counters[savedMetric.Name] = savedMetric
	} else {
		s.counters[metric.Name] = metric
	}

	return nil
}

func (s *MemStorage) checkCounterMetricNameWhenUpdate(name string) error {
	return nil
}

func (s *MemStorage) UpdateMetric(metricType, metricName, metricValue string) error {
	switch model.MetricType(metricType) {
	case model.Gauge:
		value, err := parseGaugeMetricValue(metricValue)
		if err != nil {
			return err
		}
		metric := model.NewGauge(metricName, value)
		s.UpdateGaugeMetric(metric)
	case model.Counter:
		value, err := parseCounterMetricValue(metricValue)
		if err != nil {
			return err
		}
		metric := model.NewCounter(metricName, value)
		s.UpdateCounterMetric(metric)

	default:
		return errors.NewInvalidMetricType(fmt.Sprintf("Invalid metric type: %s", metricType), nil)
	}
	return nil
}

func parseGaugeMetricValue(v string) (float64, error) {
	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, errors.NewInvalidMetricValue(
			fmt.Sprintf("The value for the %s metric should be of float64 type", model.Gauge), err)
	}

	return value, nil
}

func parseCounterMetricValue(v string) (int64, error) {
	value, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, errors.NewInvalidMetricValue(
			fmt.Sprintf("The value for the %s metric should be of int64 type", model.Counter), err)
	}

	return value, nil
}

func (s *MemStorage) ResetMetricValue(metricType, metricName string) error {
	switch model.MetricType(metricType) {
	case model.Gauge:
		err := s.checkGaugeMetricNameWhenReset(metricName)
		if err != nil {
			return err
		}

		savedMetric := s.gauges[metricName]
		savedMetric.SetValue(0)
		s.gauges[savedMetric.Name] = savedMetric
	case model.Counter:
		err := s.checkCounterMetricNameWhenReset(metricName)
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

func (s *MemStorage) checkGaugeMetricNameWhenReset(name string) error {
	_, exists := s.gauges[name]
	if !exists {
		return errors.NewInvalidMetricName(fmt.Sprintf("Gauge metric with name: %s not exists", name), nil)
	}

	return nil
}

func (s *MemStorage) checkCounterMetricNameWhenReset(name string) error {
	_, exists := s.counters[name]
	if !exists {
		return errors.NewInvalidMetricName(fmt.Sprintf("Counter metric with name: %s not exists", name), nil)
	}

	return nil
}

func (s *MemStorage) GetMetrics() (map[string]model.GaugeMetric, map[string]model.CounterMetric) {
	s.gaugesMu.Lock()
	gauges := utils.CopyMap(s.gauges)
	s.gaugesMu.Unlock()

	s.countersMu.Lock()
	counters := utils.CopyMap(s.counters)
	s.countersMu.Unlock()

	return gauges, counters
}
