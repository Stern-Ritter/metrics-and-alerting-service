package server

import (
	"context"
	"fmt"
	"time"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	storage "github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"go.uber.org/zap"
)

type MetricService struct {
	Storage storage.Storage
	Logger  *logger.ServerLogger
}

func NewMetricService(storage storage.Storage, logger *logger.ServerLogger) *MetricService {
	return &MetricService{Storage: storage, Logger: logger}
}

func (s *MetricService) UpdateMetricWithPathVars(ctx context.Context, mName string, mType string,
	mValue string, isSyncSaveStorageState bool, storageFilePath string) error {
	m, err := metrics.NewMetricsWithStringValue(mName, mType, mValue)
	if err != nil {
		return err
	}

	err = s.Storage.UpdateMetric(ctx, m)
	if err != nil {
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

func (s *MetricService) UpdateMetricWithBody(ctx context.Context, metric metrics.Metrics, isSyncSaveStorageState bool,
	storageFilePath string) (metrics.Metrics, error) {
	err := s.Storage.UpdateMetric(ctx, metric)
	if err != nil {
		return metrics.Metrics{}, err
	}

	m, err := s.Storage.GetMetric(ctx, metric)
	if err != nil {
		return metrics.Metrics{}, err
	}

	if isSyncSaveStorageState {
		err := s.Storage.Save(storageFilePath)
		if err != nil {
			s.Logger.Error(err.Error(), zap.String("event", "sync save to file storage"))
		} else {
			s.Logger.Info("Success sync save to file storage", zap.String("event", "sync save to file storage"))
		}
	}

	return m, nil
}

func (s *MetricService) UpdateMetricsBatchWithBody(ctx context.Context, metrics []metrics.Metrics,
	isSyncSaveStorageState bool, storageFilePath string) error {
	err := s.Storage.UpdateMetrics(ctx, metrics)
	if err != nil {
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

func (s *MetricService) GetMetricValueByTypeAndName(ctx context.Context, mType string, mName string) (string, error) {
	m, err := s.Storage.GetMetric(ctx, metrics.Metrics{ID: mName, MType: mType})
	if err != nil {
		return "", err
	}

	switch metrics.MetricType(m.MType) {
	case metrics.Gauge:
		return utils.FormatGaugeMetricValue(*m.Value), nil
	case metrics.Counter:
		return utils.FormatCounterMetricValue(*m.Delta), nil
	default:
		return "", errors.NewInvalidMetricType(fmt.Sprintf("Invalid metric type: %s", m.MType), nil)
	}
}

func (s *MetricService) GetMetricHandlerWithBody(ctx context.Context, metric metrics.Metrics) (metrics.Metrics, error) {
	return s.Storage.GetMetric(ctx, metric)
}

func (s *MetricService) GetMetrics(ctx context.Context) (map[string]metrics.GaugeMetric, map[string]metrics.CounterMetric, error) {
	return s.Storage.GetMetrics(ctx)
}

func (s *MetricService) RestoreMetricsFromStorage(storageFilePath string) error {
	return s.Storage.Restore(storageFilePath)
}

func (s *MetricService) SetMetricsSaveInterval(storageFilePath string, storeInterval int) {
	if storeInterval <= 0 {
		return
	}

	s.Logger.Info("Start async save to file storage", zap.String("event", "start async save to file storage"),
		zap.String("file name", storageFilePath), zap.Int("interval", storeInterval))
	go func() {
		ticker := time.NewTicker(time.Duration(storeInterval) * time.Second)
		for range ticker.C {
			if err := s.Storage.Save(storageFilePath); err != nil {
				s.Logger.Error(err.Error(), zap.String("event", "async save to file storage"))
			} else {
				s.Logger.Info("Success async save to file storage", zap.String("event", "async save to file storage"))
			}
		}
	}()
}

func (s *MetricService) PingDatabase(ctx context.Context) error {
	return s.Storage.Ping(ctx)
}
