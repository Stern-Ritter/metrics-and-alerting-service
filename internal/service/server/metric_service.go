package server

import (
	"context"
	e "errors"
	"fmt"
	"github.com/Stern-Ritter/metrics-and-alerting-service/migrations"
	"github.com/cenkalti/backoff/v4"
	"github.com/pressly/goose/v3"
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
	storageRetryInterval = backoff.NewExponentialBackOff(
		backoff.WithInitialInterval(1*time.Second),
		backoff.WithRandomizationFactor(0),
		backoff.WithMultiplier(3),
		backoff.WithMaxInterval(5*time.Second),
		backoff.WithMaxElapsedTime(10*time.Second))
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

	update := func() error {
		err = s.Storage.UpdateMetric(ctx, m)
		if isDatabaseConnectionError(err) {
			s.Logger.Error(err.Error(), zap.String("event", "failed try update metric"))
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}

	if err = backoff.Retry(update, storageRetryInterval); err != nil {
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

	update := func() error {
		err := s.Storage.UpdateMetric(ctx, metric)
		if isDatabaseConnectionError(err) {
			s.Logger.Error(err.Error(), zap.String("event", "failed try update metric"))
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}
	if err := backoff.Retry(update, storageRetryInterval); err != nil {
		return metrics.Metrics{}, err
	}

	var m metrics.Metrics
	var err error
	get := func() error {
		m, err = s.Storage.GetMetric(ctx, metric)
		if isDatabaseConnectionError(err) {
			s.Logger.Error(err.Error(), zap.String("event", "failed try get metric"))
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}

	if err := backoff.Retry(get, storageRetryInterval); err != nil {
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

	updateBatch := func() error {
		err := s.Storage.UpdateMetrics(ctx, metrics)
		if isDatabaseConnectionError(err) {
			s.Logger.Error(err.Error(), zap.String("event", "failed try update metric batch"))
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}

	if err := backoff.Retry(updateBatch, storageRetryInterval); err != nil {
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
	get := func() error {
		m, err = s.Storage.GetMetric(ctx, metrics.Metrics{ID: mName, MType: mType})
		if isDatabaseConnectionError(err) {
			s.Logger.Error(err.Error(), zap.String("event", "failed try get metric"))
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}

	if err := backoff.Retry(get, storageRetryInterval); err != nil {
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
	get := func() error {
		m, err = s.Storage.GetMetric(ctx, metric)
		if isDatabaseConnectionError(err) {
			s.Logger.Error(err.Error(), zap.String("event", "failed try get metric"))
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}

	if err := backoff.Retry(get, storageRetryInterval); err != nil {
		return metrics.Metrics{}, err
	}

	return m, err
}

func (s *MetricService) GetMetrics(ctx context.Context) (map[string]metrics.GaugeMetric,
	map[string]metrics.CounterMetric, error) {
	var gauges map[string]metrics.GaugeMetric
	var counters map[string]metrics.CounterMetric
	var err error

	getAll := func() error {
		gauges, counters, err = s.Storage.GetMetrics(ctx)
		if isDatabaseConnectionError(err) {
			s.Logger.Error(err.Error(), zap.String("event", "failed try get metrics"))
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}

	err = backoff.Retry(getAll, storageRetryInterval)

	return gauges, counters, err
}

func (s *MetricService) RestoreStateFromFile(filePath string) error {
	restore := func() error {
		err := s.Storage.Restore(filePath)
		var storageErr errors.FileUnavailable
		if e.As(err, &storageErr) {
			s.Logger.Error(err.Error(), zap.String("event", "failed try restore storage state from file"))
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}

	return backoff.Retry(restore, storageRetryInterval)
}

func (s *MetricService) SaveStateToFile(filePath string) error {
	save := func() error {
		err := s.Storage.Save(filePath)
		var storageErr errors.FileUnavailable
		if e.As(err, &storageErr) {
			s.Logger.Error(err.Error(), zap.String("event", "failed try async save to file storage"))
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}

	return backoff.Retry(save, storageRetryInterval)
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

func (s *MetricService) MigrateDatabase(databaseDsn string) error {
	goose.SetBaseFS(migrations.Migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose failed to set postgres dialect: %w", err)
	}

	db, err := goose.OpenDBWithDriver("pgx", databaseDsn)
	if err != nil {
		return fmt.Errorf("goose failed to open database connection: %w", err)
	}

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("goose failed to migrate database: %w", err)
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("goose failed to close database connection: %w", err)
	}

	return nil
}

func (s *MetricService) PingDatabase(ctx context.Context) error {
	return s.Storage.Ping(ctx)
}

func isDatabaseConnectionError(err error) bool {
	var pgErr *pgconn.PgError
	return err != nil && e.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
}
