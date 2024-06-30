package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	er "github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

// StorageState is the state of the in-memory implementation of the Storage interface.
type StorageState struct {
	Gauges   map[string]metrics.GaugeMetric   `json:"gauges"`
	Counters map[string]metrics.CounterMetric `json:"counters"`
}

// MemoryStorage is an in-memory implementation of the Storage interface.
type MemoryStorage struct {
	gaugesMu   sync.Mutex
	countersMu sync.Mutex

	gauges   map[string]metrics.GaugeMetric
	counters map[string]metrics.CounterMetric

	Logger *logger.ServerLogger
}

// NewMemoryStorage is constructor for creating a new MemoryStorage.
func NewMemoryStorage(logger *logger.ServerLogger) *MemoryStorage {
	return &MemoryStorage{
		gauges:   make(map[string]metrics.GaugeMetric),
		counters: make(map[string]metrics.CounterMetric),
		Logger:   logger,
	}
}

// UpdateMetric updates a single metric in the memory storage.
func (s *MemoryStorage) UpdateMetric(ctx context.Context, metric metrics.Metrics) error {
	switch metrics.MetricType(metric.MType) {
	case metrics.Gauge:
		m := metrics.MetricsToGaugeMetric(metric)
		_, err := s.updateGaugeMetric(m)
		if err != nil {
			return err
		}
	case metrics.Counter:
		m := metrics.MetricsToCounterMetric(metric)
		_, err := s.updateCounterMetric(m)
		if err != nil {
			return err
		}
	default:
		return er.NewInvalidMetricType(fmt.Sprintf("Invalid metric type: %s", metric.MType), nil)
	}

	return nil
}

// UpdateMetrics updates multiple metrics in the memory storage.
func (s *MemoryStorage) UpdateMetrics(ctx context.Context, metricsBatch []metrics.Metrics) error {
	for _, metric := range metricsBatch {
		switch metrics.MetricType(metric.MType) {
		case metrics.Gauge:
			m := metrics.MetricsToGaugeMetric(metric)
			_, err := s.updateGaugeMetric(m)
			if err != nil {
				return err
			}
		case metrics.Counter:
			m := metrics.MetricsToCounterMetric(metric)
			_, err := s.updateCounterMetric(m)
			if err != nil {
				return err
			}
		default:
			return er.NewInvalidMetricType(fmt.Sprintf("Invalid metric type: %s", metric.MType), nil)
		}
	}

	return nil
}

func (s *MemoryStorage) updateGaugeMetric(metric metrics.GaugeMetric) (metrics.GaugeMetric, error) {
	s.gaugesMu.Lock()
	defer s.gaugesMu.Unlock()

	if savedMetric, exists := s.gauges[metric.Name]; exists {
		savedMetric.SetValue(metric.GetValue())
		s.gauges[metric.Name] = savedMetric
	} else {
		s.gauges[metric.Name] = metric
	}

	return s.gauges[metric.Name], nil
}

func (s *MemoryStorage) updateCounterMetric(metric metrics.CounterMetric) (metrics.CounterMetric, error) {
	s.countersMu.Lock()
	defer s.countersMu.Unlock()

	if savedMetric, exists := s.counters[metric.Name]; exists {
		savedMetric.SetValue(metric.GetValue())
		s.counters[metric.Name] = savedMetric
	} else {
		s.counters[metric.Name] = metric
	}

	return s.counters[metric.Name], nil
}

// GetMetric gets a single metric from the memory storage.
func (s *MemoryStorage) GetMetric(ctx context.Context, metric metrics.Metrics) (metrics.Metrics, error) {
	switch metrics.MetricType(metric.MType) {
	case metrics.Gauge:
		s.gaugesMu.Lock()
		m, exists := s.gauges[metric.ID]
		s.gaugesMu.Unlock()
		if !exists {
			return metrics.Metrics{}, er.NewInvalidMetricName(fmt.Sprintf("Gauge metric with name: %s not exists", metric.ID), nil)
		}
		return metrics.GaugeMetricToMetrics(m), nil

	case metrics.Counter:
		s.countersMu.Lock()
		m, exists := s.counters[metric.ID]
		s.countersMu.Unlock()
		if !exists {
			return metrics.Metrics{}, er.NewInvalidMetricName(fmt.Sprintf("Counter metric with name: %s not exists", metric.ID), nil)
		}
		return metrics.CounterMetricToMetrics(m), nil

	default:
		return metrics.Metrics{}, er.NewInvalidMetricType(fmt.Sprintf("Invalid metric type: %s", metric.MType), nil)
	}
}

// GetMetrics gets all metrics from the memory storage.
func (s *MemoryStorage) GetMetrics(ctx context.Context) (map[string]metrics.GaugeMetric, map[string]metrics.CounterMetric, error) {
	s.gaugesMu.Lock()
	gauges := utils.CopyMap(s.gauges)
	s.gaugesMu.Unlock()

	s.countersMu.Lock()
	counters := utils.CopyMap(s.counters)
	s.countersMu.Unlock()

	return gauges, counters, nil
}

// Ping checks the connection to the memory storage.
// This operation is not supported for the in-memory implementation of the Storage
func (s *MemoryStorage) Ping(ctx context.Context) error {
	return fmt.Errorf("the database is disabled")
}

// Restore restores the memory storage state from the file.
func (s *MemoryStorage) Restore(fname string) error {
	file, err := os.OpenFile(fname, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return er.NewFileUnavailable(fmt.Sprintf("can not open file %s to restore state: %v", fname, err), err)
	}
	defer file.Close()

	state := StorageState{
		Gauges:   make(map[string]metrics.GaugeMetric),
		Counters: make(map[string]metrics.CounterMetric),
	}

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil
	}
	data := scanner.Bytes()

	err = json.Unmarshal(data, &state)
	if err != nil {
		return err
	}

	s.gaugesMu.Lock()
	s.gauges = state.Gauges
	s.gaugesMu.Unlock()

	s.countersMu.Lock()
	s.counters = state.Counters
	s.countersMu.Unlock()

	return nil
}

// Save saves the memory storage state to the file.
func (s *MemoryStorage) Save(fname string) error {
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return er.NewFileUnavailable(fmt.Sprintf("can not open file %s to save state: %v", fname, err), err)
	}
	defer file.Close()

	s.gaugesMu.Lock()
	s.countersMu.Lock()
	state := StorageState{
		Gauges:   s.gauges,
		Counters: s.counters,
	}

	data, err := json.Marshal(&state)
	if err != nil {
		return err
	}
	s.gaugesMu.Unlock()
	s.countersMu.Unlock()

	_, err = file.Write(data)
	return err
}

// SetGaugeMetircs sets the gauge metrics in the memory storage.
func (s *MemoryStorage) SetGaugeMetircs(gauges map[string]metrics.GaugeMetric) {
	s.gaugesMu.Lock()
	defer s.gaugesMu.Unlock()
	s.gauges = gauges
}

// SetCounterMetrics sets the counter metrics in the memory storage.
func (s *MemoryStorage) SetCounterMetrics(counters map[string]metrics.CounterMetric) {
	s.countersMu.Lock()
	defer s.countersMu.Unlock()
	s.counters = counters
}
