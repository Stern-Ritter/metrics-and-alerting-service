package server

import (
	"fmt"

	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"go.uber.org/zap"
)

type MetricService struct {
	Storage storage.ServerStorage
	Logger  *logger.ServerLogger
}

func NewMetricService(storage storage.ServerStorage, logger *logger.ServerLogger) *MetricService {
	return &MetricService{Storage: storage, Logger: logger}
}

func (s *MetricService) UpdateMetricWithPathVars(metricType string, metricName string,
	metricValue string, isSyncSaveStorageState bool, storageFilePath string) error {
	if err := s.Storage.UpdateMetric(metricType, metricName, metricValue); err != nil {
		return err
	}

	if isSyncSaveStorageState {
		err := s.Storage.Save(storageFilePath)
		if err != nil {
			s.Logger.Error(err.Error(), zap.String("event", "sync save to file storage"))
		} else {
			s.Logger.Info("Success sync save to file storage", zap.String("event", "sync save to file storage"))
		}
	}
	return nil
}

func (s *MetricService) UpdateMetricWithBody(metric metrics.Metrics, isSyncSaveStorageState bool,
	storageFilePath string) (metrics.Metrics, error) {
	switch metrics.MetricType(metric.MType) {
	case metrics.Gauge:
		updatedMetric, err := s.Storage.UpdateGaugeMetric(metrics.MetricsToGaugeMetric(metric))
		if err != nil {
			return metric, err
		}

		value := updatedMetric.GetValue()
		metric.Value = &value

	case metrics.Counter:
		updatedMetric, err := s.Storage.UpdateCounterMetric(metrics.MetricsToCounterMetric(metric))
		if err != nil {
			return metric, err
		}

		delta := updatedMetric.GetValue()
		metric.Delta = &delta

	default:
		return metric, fmt.Errorf("invalid metric type: %s", metric.MType)
	}

	if isSyncSaveStorageState {
		err := s.Storage.Save(storageFilePath)
		if err != nil {
			s.Logger.Error(err.Error(), zap.String("event", "sync save to file storage"))
		} else {
			s.Logger.Info("Success sync save to file storage", zap.String("event", "sync save to file storage"))
		}
	}

	return metric, nil
}

func (s *MetricService) GetMetricValueByTypeAndName(metricType string, metricName string) (string, error) {
	return s.Storage.GetMetricValueByTypeAndName(metricType, metricName)
}

func (s *MetricService) GetMetricHandlerWithBody(metric metrics.Metrics) (metrics.Metrics, error) {
	switch metrics.MetricType(metric.MType) {
	case metrics.Gauge:
		savedMetric, err := s.Storage.GetGaugeMetric(metric.ID)
		if err != nil {
			metric.Value = &metrics.ZeroGaugeMetricValue
			break
		}

		value := savedMetric.GetValue()
		metric.Value = &value

	case metrics.Counter:
		savedMetric, err := s.Storage.GetCounterMetric(metric.ID)
		if err != nil {
			metric.Delta = &metrics.ZeroCounterMetricValue
			break
		}

		delta := savedMetric.GetValue()
		metric.Delta = &delta

	default:
		return metric, fmt.Errorf("invalid metric type: %s", metric.MType)
	}

	return metric, nil
}

func (s *MetricService) GetMetrics() (map[string]metrics.GaugeMetric, map[string]metrics.CounterMetric) {
	return s.Storage.GetMetrics()
}

func (s *MetricService) RestoreMetricsFromStorage(storageFilePath string) error {
	return s.Storage.Restore(storageFilePath)
}

func (s *MetricService) SetMetricsSaveInterval(storageFilePath string, storeInterval int) {
	s.Storage.SetSaveInterval(storageFilePath, storeInterval)
}
