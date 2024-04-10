package server

import (
	"database/sql"
	e "errors"
	"fmt"
	"os"

	"context"

	errors "github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/jackc/pgx/v5"
)

type DBStorage struct {
	conn   *sql.DB
	Logger *logger.ServerLogger
}

func NewDBStorage(conn *sql.DB, logger *logger.ServerLogger) *DBStorage {
	return &DBStorage{conn: conn, Logger: logger}
}

func (s *DBStorage) UpdateMetric(ctx context.Context, metric metrics.Metrics) error {
	mValue, err := metric.GetValue()
	if err != nil {
		return err
	}

	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	row := s.conn.QueryRowContext(ctx, `
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

	var mId int64
	var mSavedValue float64
	err = row.Scan(&mId, &mValue)

	if err != nil {
		if !e.Is(err, sql.ErrNoRows) {
			return err
		}
		err = saveMetric(ctx, tx, metric.ID, metric.MType, mValue)
	} else {
		switch metrics.MetricType(metric.MType) {
		case metrics.Gauge:
			err = updateMetric(ctx, tx, mId, mValue)
		case metrics.Counter:
			err = updateMetric(ctx, tx, mId, mSavedValue+mSavedValue)
		}
	}

	if err != nil {
		return err
	}

	return tx.Commit()
}

func saveMetric(ctx context.Context, tx *sql.Tx, mName string, mType string, mValue float64) error {
	mTypeId, err := getMetricTypeId(ctx, tx, mType)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO metrics
		(name, type_id, value)
		VALUES
		(@name, @type_id, @value)
	`, pgx.NamedArgs{
		"name":    mName,
		"type_id": mTypeId,
		"value":   mValue,
	})

	return err
}

func updateMetric(ctx context.Context, tx *sql.Tx, mId int64, mValue float64) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE metrics
		SET value = @value WHERE id = @id
	`, pgx.NamedArgs{
		"id":    mId,
		"value": mValue,
	})

	return err
}

func getMetricTypeId(ctx context.Context, tx *sql.Tx, mType string) (int64, error) {
	row := tx.QueryRowContext(ctx, `
		SELECT
			id
		FROM metric_types
		WHERE name = @type
	`, pgx.NamedArgs{"type": mType})

	var mTypeId int64
	err := row.Scan(&mTypeId)
	if err != nil {
		return 0, errors.NewInvalidMetricType(fmt.Sprintf("Invalid metric type: %s", mType), err)
	}

	return mTypeId, nil
}

func (s *DBStorage) GetMetric(ctx context.Context, metric metrics.Metrics) (metrics.Metrics, error) {
	row := s.conn.QueryRowContext(ctx, `
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
		return metrics.Metrics{}, err
	}

	m, err := metrics.NewMetrics(mName, mType, fmt.Sprint(mValue))
	if err != nil {
		return metrics.Metrics{}, err
	}

	return m, err
}

func (s *DBStorage) GetMetrics(ctx context.Context) (map[string]metrics.GaugeMetric, map[string]metrics.CounterMetric, error) {
	rows, err := s.conn.QueryContext(ctx, `
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
	return s.conn.PingContext(ctx)
}

func (s *DBStorage) Restore(fname string) error {
	return fmt.Errorf("can not restore database storage state from file")
}

func (s *DBStorage) Save(fname string) error {
	return fmt.Errorf("can not save database storage state to file")
}

func (s *DBStorage) InitDatabase(initScriptPath string) error {
	data, err := os.ReadFile(initScriptPath)
	if err != nil {
		return err
	}

	initScript := string(data)
	_, err = s.conn.Exec(initScript)
	return err
}
