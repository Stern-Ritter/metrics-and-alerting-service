package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	er "github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
)

type DBStorage struct {
	db     *sql.DB
	Logger *logger.ServerLogger
}

func NewDBStorage(db *sql.DB, logger *logger.ServerLogger) *DBStorage {
	return &DBStorage{db: db, Logger: logger}
}

func (s *DBStorage) UpdateMetric(ctx context.Context, metric metrics.Metrics) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	err = s.updateMetricInTx(ctx, tx, metric)
	if err != nil {
		//nolint:errcheck
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *DBStorage) UpdateMetrics(ctx context.Context, metricsBatch []metrics.Metrics) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	for _, metric := range metricsBatch {
		err = s.updateMetricInTx(ctx, tx, metric)
		if err != nil {
			//nolint:errcheck
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *DBStorage) updateMetricInTx(ctx context.Context, tx *sql.Tx, metric metrics.Metrics) error {
	mValue, err := metric.GetValue()
	if err != nil {
		return err
	}

	row := tx.QueryRowContext(ctx, `
		SELECT
			id,
			value
		FROM metrics
		WHERE
			name = $1 AND
			type = $2
	`, metric.ID, metric.MType)

	var mID int64
	var mSavedValue float64
	err = row.Scan(&mID, &mSavedValue)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		err = saveMetric(ctx, tx, metric.ID, metric.MType, mValue)
	} else {
		switch metrics.MetricType(metric.MType) {
		case metrics.Gauge:
			err = updateMetric(ctx, tx, mID, mValue)
		case metrics.Counter:
			err = updateMetric(ctx, tx, mID, mValue+mSavedValue)
		}
	}

	return err
}

func saveMetric(ctx context.Context, tx *sql.Tx, mName string, mType string, mValue float64) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO metrics
		(name, type, value)
		VALUES
		($1, $2, $3)
	`, mName, mType, mValue)

	return err
}

func updateMetric(ctx context.Context, tx *sql.Tx, mID int64, mValue float64) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE metrics
		SET value = $1 WHERE id = $2
	`, mValue, mID)

	return err
}

func (s *DBStorage) GetMetric(ctx context.Context, metric metrics.Metrics) (metrics.Metrics, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT
			name,
			type,
			value
		FROM metrics
		WHERE
			name = $1 AND
			type = $2
	`, metric.ID, metric.MType)

	var mName string
	var mType string
	var mValue float64
	err := row.Scan(&mName, &mType, &mValue)

	if err != nil {
		return metrics.Metrics{}, er.NewInvalidMetricName(fmt.Sprintf("Metric with name: %s not exists", metric.ID), nil)
	}

	m, err := metrics.NewMetricsWithNumberValue(mName, mType, mValue)
	if err != nil {
		return metrics.Metrics{}, err
	}

	return m, err
}

func (s *DBStorage) GetMetrics(ctx context.Context) (map[string]metrics.GaugeMetric, map[string]metrics.CounterMetric, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			name,
			type,
			value
		FROM metrics
		WHERE
			type IN($1, $2)
	`, metrics.Gauge, metrics.Counter)

	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	gauges := make(map[string]metrics.GaugeMetric)
	counters := make(map[string]metrics.CounterMetric)

	for rows.Next() {
		var mName string
		var mType string
		var mValue float64

		if err := rows.Scan(&mName, &mType, &mValue); err != nil {
			return nil, nil, err
		}

		switch metrics.MetricType(mType) {
		case metrics.Gauge:
			gauge := metrics.NewGauge(mName, mValue)
			gauges[mName] = gauge
		case metrics.Counter:
			counter := metrics.NewCounter(mName, int64(mValue))
			counters[mName] = counter
		}
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return gauges, counters, nil
}

func (s *DBStorage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *DBStorage) Restore(fName string) error {
	return fmt.Errorf("can not restore database storage state from file")
}

func (s *DBStorage) Save(fName string) error {
	return fmt.Errorf("can not save database storage state to file: %s", fName)
}
