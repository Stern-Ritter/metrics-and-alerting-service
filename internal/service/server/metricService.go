package server

import (
	"context"
	e "errors"
	"fmt"
	"time"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	storage "github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

var (
	storageRetryInterval = []int{1, 3, 5}
)

type MetricService struct {
	Storage storage.Storage
	Logger  *logger.ServerLogger
}

func NewMetricService(storage storage.Storage, logger *logger.ServerLogger) *MetricService {
	return &MetricService{Storage: storage, Logger: logger}
}

func (s *MetricService) UpdateMetricWithPathVars(ctx context.Context, mName string, mType string,
	mValue string, isSyncSaveStorageState bool, filePath string) error {
	m, err := metrics.NewMetricsWithStringValue(mName, mType, mValue)
	if err != nil {
		return err
	}

	for _, interval := range storageRetryInterval {
		err = s.Storage.UpdateMetric(ctx, m)
		if isDatabaseConnectionError(err) {
			s.Logger.Error(err.Error(), zap.String("event", "failed try update metric"))
			time.Sleep(time.Duration(interval) * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		return err
	}

	if isSyncSaveStorageState {
		err := s.SaveStateToFile(filePath)
		if err != nil {
			s.Logger.Error(err.Error(), zap.String("event", "sync save to file storage"))
		} else {
			s.Logger.Info("Success sync save to file storage", zap.String("event", "sync save to file storage"))
		}
	}
	return nil
}

func (s *MetricService) UpdateMetricWithBody(ctx context.Context, metric metrics.Metrics, isSyncSaveStorageState bool,
	filePath string) (metrics.Metrics, error) {

	var err error
	for _, interval := range storageRetryInterval {
		err = s.Storage.UpdateMetric(ctx, metric)
		if isDatabaseConnectionError(err) {
			s.Logger.Error(err.Error(), zap.String("event", "failed try update metric"))
			time.Sleep(time.Duration(interval) * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		return metrics.Metrics{}, err
	}

	var m metrics.Metrics
	for _, interval := range storageRetryInterval {
		m, err = s.Storage.GetMetric(ctx, metric)
		if isDatabaseConnectionError(err) {
			s.Logger.Error(err.Error(), zap.String("event", "failed try get metric"))
			time.Sleep(time.Duration(interval) * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		return metrics.Metrics{}, err
	}

	if isSyncSaveStorageState {
		err := s.SaveStateToFile(filePath)
		if err != nil {
			s.Logger.Error(err.Error(), zap.String("event", "sync save to file storage"))
		} else {
			s.Logger.Info("Success sync save to file storage", zap.String("event", "sync save to file storage"))
		}
	}

	return m, nil
}

func (s *MetricService) UpdateMetricsBatchWithBody(ctx context.Context, metrics []metrics.Metrics,
	isSyncSaveStorageState bool, filePath string) error {
	var err error
	for _, interval := range storageRetryInterval {
		err = s.Storage.UpdateMetrics(ctx, metrics)
		if isDatabaseConnectionError(err) {
			s.Logger.Error(err.Error(), zap.String("event", "failed try update metric batch"))
			time.Sleep(time.Duration(interval) * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		return err
	}

	if isSyncSaveStorageState {
		err := s.SaveStateToFile(filePath)
		if err != nil {
			s.Logger.Error(err.Error(), zap.String("event", "sync save to file storage"))
		} else {
			s.Logger.Info("Success sync save to file storage", zap.String("event", "sync save to file storage"))
		}
	}
	return nil
}

func (s *MetricService) GetMetricValueByTypeAndName(ctx context.Context, mType string, mName string) (string, error) {
	var m metrics.Metrics
	var err error
	for _, interval := range storageRetryInterval {
		m, err = s.Storage.GetMetric(ctx, metrics.Metrics{ID: mName, MType: mType})
		if isDatabaseConnectionError(err) {
			s.Logger.Error(err.Error(), zap.String("event", "failed try get metric"))
			time.Sleep(time.Duration(interval) * time.Second)
		} else {
			break
		}
	}
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
	var m metrics.Metrics
	var err error
	for _, interval := range storageRetryInterval {
		m, err = s.Storage.GetMetric(ctx, metric)
		if isDatabaseConnectionError(err) {
			s.Logger.Error(err.Error(), zap.String("event", "failed try get metric"))
			time.Sleep(time.Duration(interval) * time.Second)
		} else {
			break
		}
	}
	return m, err
}

func (s *MetricService) GetMetrics(ctx context.Context) (map[string]metrics.GaugeMetric,
	map[string]metrics.CounterMetric, error) {
	var gauges map[string]metrics.GaugeMetric
	var counters map[string]metrics.CounterMetric
	var err error
	for _, interval := range storageRetryInterval {
		gauges, counters, err = s.Storage.GetMetrics(ctx)
		if isDatabaseConnectionError(err) {
			s.Logger.Error(err.Error(), zap.String("event", "failed try get metrics"))
			time.Sleep(time.Duration(interval) * time.Second)
		} else {
			break
		}
	}
	return gauges, counters, err
}

func (s *MetricService) RestoreStateFromFile(filePath string) error {
	var err error
	for _, interval := range storageRetryInterval {
		err = s.Storage.Restore(filePath)
		var storageErr errors.FileUnavailable
		if !e.As(err, &storageErr) {
			return err
		}
		s.Logger.Error(err.Error(), zap.String("event", "failed try restore storage state from file"))
		time.Sleep(time.Duration(interval) * time.Second)
	}
	return err
}

func (s *MetricService) SaveStateToFile(filePath string) error {
	var err error
	for _, interval := range storageRetryInterval {
		err = s.Storage.Save(filePath)
		var storageErr errors.FileUnavailable
		if !e.As(err, &storageErr) {
			err = storageErr
			break
		}
		s.Logger.Error(err.Error(), zap.String("event", "failed try async save to file storage"))
		time.Sleep(time.Duration(interval) * time.Second)
	}
	return err
}

func (s *MetricService) SetSaveStateToFileInterval(filePath string, storeInterval int) {
	if storeInterval <= 0 {
		return
	}

	s.Logger.Info("Start async save to file storage", zap.String("event", "start async save to file storage"),
		zap.String("file path", filePath), zap.Int("interval", storeInterval))
	go func() {
		ticker := time.NewTicker(time.Duration(storeInterval) * time.Second)
		for range ticker.C {
			err := s.SaveStateToFile(filePath)
			if err != nil {
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

func isDatabaseConnectionError(err error) bool {
	var pgErr *pgconn.PgError
	return err != nil && e.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
}
