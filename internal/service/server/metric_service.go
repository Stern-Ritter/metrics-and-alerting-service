package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/pressly/goose/v3"

	"github.com/Stern-Ritter/metrics-and-alerting-service/migrations"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"

	er "github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	storage "github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

// MetricService is a service for managing metrics.
type MetricService struct {
	storage              storage.Storage
	logger               *logger.ServerLogger
	storageRetryInterval *backoff.ExponentialBackOff
}

// NewMetricService is constructor for creating a new MetricService.
func NewMetricService(storage storage.Storage, logger *logger.ServerLogger) *MetricService {
	storageRetryInterval := backoff.NewExponentialBackOff(
		backoff.WithInitialInterval(1*time.Second),
		backoff.WithRandomizationFactor(0),
		backoff.WithMultiplier(3),
		backoff.WithMaxInterval(5*time.Second),
		backoff.WithMaxElapsedTime(10*time.Second))

	return &MetricService{storage: storage, logger: logger, storageRetryInterval: storageRetryInterval}
}

// UpdateMetricWithPathVars updates a metric using string params.
func (s *MetricService) UpdateMetricWithPathVars(ctx context.Context, mName string, mType string,
	mValue string, isSyncSaveStorageState bool, filePath string) error {
	m, err := metrics.NewMetricsWithStringValue(mName, mType, mValue)
	if err != nil {
		return err
	}

	update := func() error {
		err = s.storage.UpdateMetric(ctx, m)
		if isDatabaseConnectionError(err) {
			s.logger.Error(err.Error(), zap.String("event", "failed try update metric"))
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}

	if err = backoff.Retry(update, s.storageRetryInterval); err != nil {
		return err
	}

	if isSyncSaveStorageState {
		err := s.SaveStateToFile(filePath)
		if err != nil {
			s.logger.Error(err.Error(), zap.String("event", "sync save to file storage"))
		} else {
			s.logger.Info("Success sync save to file storage", zap.String("event", "sync save to file storage"))
		}
	}
	return nil
}

// UpdateMetricWithBody updates a metric using Metrics object
func (s *MetricService) UpdateMetricWithBody(ctx context.Context, metric metrics.Metrics, isSyncSaveStorageState bool,
	filePath string) (metrics.Metrics, error) {

	update := func() error {
		err := s.storage.UpdateMetric(ctx, metric)
		if isDatabaseConnectionError(err) {
			s.logger.Error(err.Error(), zap.String("event", "failed try update metric"))
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}
	if err := backoff.Retry(update, s.storageRetryInterval); err != nil {
		return metrics.Metrics{}, err
	}

	var m metrics.Metrics
	var err error
	get := func() error {
		m, err = s.storage.GetMetric(ctx, metric)
		if isDatabaseConnectionError(err) {
			s.logger.Error(err.Error(), zap.String("event", "failed try get metric"))
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}

	if err := backoff.Retry(get, s.storageRetryInterval); err != nil {
		return metrics.Metrics{}, err
	}

	if isSyncSaveStorageState {
		err := s.SaveStateToFile(filePath)
		if err != nil {
			s.logger.Error(err.Error(), zap.String("event", "sync save to file storage"))
		} else {
			s.logger.Info("Success sync save to file storage", zap.String("event", "sync save to file storage"))
		}
	}

	return m, nil
}

// UpdateMetricsBatchWithBody updates a slice of metrics.
func (s *MetricService) UpdateMetricsBatchWithBody(ctx context.Context, metrics []metrics.Metrics,
	isSyncSaveStorageState bool, filePath string) error {

	updateBatch := func() error {
		err := s.storage.UpdateMetrics(ctx, metrics)
		if isDatabaseConnectionError(err) {
			s.logger.Error(err.Error(), zap.String("event", "failed try update metric batch"))
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}

	if err := backoff.Retry(updateBatch, s.storageRetryInterval); err != nil {
		return err
	}

	if isSyncSaveStorageState {
		err := s.SaveStateToFile(filePath)
		if err != nil {
			s.logger.Error(err.Error(), zap.String("event", "sync save to file storage"))
		} else {
			s.logger.Info("Success sync save to file storage", zap.String("event", "sync save to file storage"))
		}
	}
	return nil
}

// GetMetricValueByTypeAndName returns a string with value of the metric by metric type and name.
func (s *MetricService) GetMetricValueByTypeAndName(ctx context.Context, mType string, mName string) (string, error) {
	var m metrics.Metrics
	var err error
	get := func() error {
		m, err = s.storage.GetMetric(ctx, metrics.Metrics{ID: mName, MType: mType})
		if isDatabaseConnectionError(err) {
			s.logger.Error(err.Error(), zap.String("event", "failed try get metric"))
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}

	if err := backoff.Retry(get, s.storageRetryInterval); err != nil {
		return "", err
	}

	switch metrics.MetricType(m.MType) {
	case metrics.Gauge:
		return utils.FormatGaugeMetricValue(*m.Value), nil
	case metrics.Counter:
		return utils.FormatCounterMetricValue(*m.Delta), nil
	default:
		return "", er.NewInvalidMetricType(fmt.Sprintf("Invalid metric type: %s", m.MType), nil)
	}
}

// GetMetricHandlerWithBody returns a metric by metric type and name.
func (s *MetricService) GetMetricHandlerWithBody(ctx context.Context, metric metrics.Metrics) (metrics.Metrics, error) {
	var m metrics.Metrics
	var err error
	get := func() error {
		m, err = s.storage.GetMetric(ctx, metric)
		if isDatabaseConnectionError(err) {
			s.logger.Error(err.Error(), zap.String("event", "failed try get metric"))
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}

	if err := backoff.Retry(get, s.storageRetryInterval); err != nil {
		return metrics.Metrics{}, err
	}

	return m, err
}

// GetMetrics returns all metrics.
func (s *MetricService) GetMetrics(ctx context.Context) (map[string]metrics.GaugeMetric,
	map[string]metrics.CounterMetric, error) {
	var gauges map[string]metrics.GaugeMetric
	var counters map[string]metrics.CounterMetric
	var err error

	getAll := func() error {
		gauges, counters, err = s.storage.GetMetrics(ctx)
		if isDatabaseConnectionError(err) {
			s.logger.Error(err.Error(), zap.String("event", "failed try get metrics"))
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}

	err = backoff.Retry(getAll, s.storageRetryInterval)

	return gauges, counters, err
}

// RestoreStateFromFile restores the storage state from a file.
func (s *MetricService) RestoreStateFromFile(filePath string) error {
	restore := func() error {
		err := s.storage.Restore(filePath)
		var storageErr er.FileUnavailable
		if errors.As(err, &storageErr) {
			s.logger.Error(err.Error(), zap.String("event", "failed try restore storage state from file"))
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}

	return backoff.Retry(restore, s.storageRetryInterval)
}

// SaveStateToFile saves the storage state to a file.
func (s *MetricService) SaveStateToFile(filePath string) error {
	save := func() error {
		err := s.storage.Save(filePath)
		var storageErr er.FileUnavailable
		if errors.As(err, &storageErr) {
			s.logger.Error(err.Error(), zap.String("event", "failed try async save to file storage"))
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}

	return backoff.Retry(save, s.storageRetryInterval)
}

// SetSaveStateToFileInterval sets an interval to save the storage state to a file.
func (s *MetricService) SetSaveStateToFileInterval(filePath string, storeInterval int) {
	if storeInterval <= 0 {
		return
	}

	s.logger.Info("Start async save to file storage", zap.String("event", "start async save to file storage"),
		zap.String("file path", filePath), zap.Int("interval", storeInterval))
	go func() {
		ticker := time.NewTicker(time.Duration(storeInterval) * time.Second)
		for range ticker.C {
			err := s.SaveStateToFile(filePath)
			if err != nil {
				s.logger.Error(err.Error(), zap.String("event", "async save to file storage"))
			} else {
				s.logger.Info("Success async save to file storage", zap.String("event", "async save to file storage"))
			}
		}
	}()
}

// MigrateDatabase migrates the database using goose.
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

// PingDatabase checks the connection to the database.
func (s *MetricService) PingDatabase(ctx context.Context) error {
	return s.storage.Ping(ctx)
}

func isDatabaseConnectionError(err error) bool {
	var pgErr *pgconn.PgError
	return err != nil && errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
}
