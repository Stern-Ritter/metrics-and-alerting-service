package server

import (
	e "errors"
	"fmt"

	"context"

	errors "github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/jackc/pgx/v5"
)

type DBStorage struct {
	conn   *pgx.Conn
	Logger *logger.ServerLogger
}

func NewDBStorage(conn *pgx.Conn, logger *logger.ServerLogger) *DBStorage {
	return &DBStorage{conn: conn, Logger: logger}
}

func (s *DBStorage) Bootstrap(ctx context.Context) error {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		DROP TABLE IF EXISTS metrics;
		DROP TABLE IF EXISTS metric_types;
    `)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE metric_types (
			id SERIAL PRIMARY KEY,
			name VARCHAR(256) NOT NULL
		);
	`)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, `
		CREATE UNIQUE INDEX metric_type_name_idx ON metric_types (name);
	`)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE metrics (
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(256) NOT NULL,
			type_id INTEGER NOT NULL,
			value DOUBLE PRECISION NOT NULL,
			FOREIGN KEY (type_id) REFERENCES metric_types (id)
		);
	`)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, `
		CREATE INDEX metric_name_idx ON metrics (name);
	`)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO metric_types (name) VALUES ('gauge'), ('counter');
	`)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (s *DBStorage) UpdateMetric(ctx context.Context, metric metrics.Metrics) error {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}

	err = s.updateMetricInTx(ctx, tx, metric)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (s *DBStorage) UpdateMetrics(ctx context.Context, metricsBatch []metrics.Metrics) error {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}

	for _, metric := range metricsBatch {
		err = s.updateMetricInTx(ctx, tx, metric)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *DBStorage) updateMetricInTx(ctx context.Context, tx pgx.Tx, metric metrics.Metrics) error {
	mValue, err := metric.GetValue()
	if err != nil {
		return err
	}

	row := tx.QueryRow(ctx, `
		SELECT
			m.id,
			m.value
		FROM metrics as m
		LEFT JOIN metric_types mt ON m.type_id = mt.id
		WHERE
			m.name = @name AND
			mt.name = @type
	`, pgx.NamedArgs{
		"name": metric.ID,
		"type": metric.MType,
	})

	var mID int64
	var mSavedValue float64
	err = row.Scan(&mID, &mSavedValue)

	if err != nil {
		if !e.Is(err, pgx.ErrNoRows) {
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

func saveMetric(ctx context.Context, tx pgx.Tx, mName string, mType string, mValue float64) error {
	mTypeID, err := getMetricTypeID(ctx, tx, mType)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO metrics
		(name, type_id, value)
		VALUES
		(@name, @type_id, @value)
	`, pgx.NamedArgs{
		"name":    mName,
		"type_id": mTypeID,
		"value":   mValue,
	})

	return err
}

func updateMetric(ctx context.Context, tx pgx.Tx, mID int64, mValue float64) error {
	_, err := tx.Exec(ctx, `
		UPDATE metrics
		SET value = @value WHERE id = @id
	`, pgx.NamedArgs{
		"id":    mID,
		"value": mValue,
	})

	return err
}

func getMetricTypeID(ctx context.Context, tx pgx.Tx, mType string) (int64, error) {
	row := tx.QueryRow(ctx, `
		SELECT
			id
		FROM metric_types
		WHERE name = @type
	`, pgx.NamedArgs{"type": mType})

	var mTypeID int64
	err := row.Scan(&mTypeID)
	if err != nil {
		return 0, errors.NewInvalidMetricType(fmt.Sprintf("Invalid metric type: %s", mType), err)
	}

	return mTypeID, nil
}

func (s *DBStorage) GetMetric(ctx context.Context, metric metrics.Metrics) (metrics.Metrics, error) {
	row := s.conn.QueryRow(ctx, `
		SELECT
			m.name,
			mt.name as type,
			m.value
		FROM metrics as m
		LEFT JOIN metric_types mt ON m.type_id = mt.id
		WHERE
			m.name = @name AND
			mt.name = @type
	`, pgx.NamedArgs{
		"name": metric.ID,
		"type": metric.MType,
	})

	var mName string
	var mType string
	var mValue float64
	err := row.Scan(&mName, &mType, &mValue)

	if err != nil {
		return metrics.Metrics{}, errors.NewInvalidMetricName(fmt.Sprintf("Metric with name: %s not exists", metric.ID), nil)
	}

	m, err := metrics.NewMetricsWithNumberValue(mName, mType, mValue)
	if err != nil {
		return metrics.Metrics{}, err
	}

	return m, err
}

func (s *DBStorage) GetMetrics(ctx context.Context) (map[string]metrics.GaugeMetric, map[string]metrics.CounterMetric, error) {
	rows, err := s.conn.Query(ctx, `
		SELECT
			m.name,
			mt.name as type,
			m.value
		FROM metrics as m
		LEFT JOIN metric_types mt ON m.type_id = mt.id
		WHERE
			mt.name IN(@gaugeType, @counterType)
	`, pgx.NamedArgs{
		"gaugeType":   metrics.Gauge,
		"counterType": metrics.Counter,
	})

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
	return s.conn.Ping(ctx)
}

func (s *DBStorage) Restore(fname string) error {
	return fmt.Errorf("can not restore database storage state from file")
}

func (s *DBStorage) Save(fname string) error {
	return fmt.Errorf("can not save database storage state to file")
}
