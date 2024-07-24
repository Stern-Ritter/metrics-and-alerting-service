package server

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
)

func TestDBStorage_UpdateMetric(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "unexpected error when mock db connection")
	defer db.Close()

	lg := &logger.ServerLogger{}
	storage := NewDBStorage(db, lg)

	metric := metrics.Metrics{
		ID:    "testMetric",
		MType: "gauge",
		Value: floatPtr(100),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT id, value FROM metrics WHERE name = \$1 AND type = \$2`).
		WithArgs(metric.ID, metric.MType).
		WillReturnRows(sqlmock.NewRows(nil))
	mock.ExpectExec(`INSERT INTO metrics \(name, type, value\) VALUES \(\$1, \$2, \$3\)`).
		WithArgs(metric.ID, metric.MType, *metric.Value).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = storage.UpdateMetric(context.Background(), metric)
	assert.NoError(t, err, "unexpected error when update metric")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err, "not all queued expectations were met in order")
}

func TestDBStorage_UpdateMetrics(t *testing.T) {
	t.Run("Update gauge metrics", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err, "unexpected error when mock db connection")
		defer db.Close()

		lg := &logger.ServerLogger{}
		storage := NewDBStorage(db, lg)

		metricsBatch := []metrics.Metrics{
			{
				ID:    "first",
				MType: "gauge",
				Value: floatPtr(10.11),
			},
			{
				ID:    "second",
				MType: "gauge",
				Value: floatPtr(20.22),
			},
		}

		mock.ExpectBegin()
		for _, metric := range metricsBatch {
			mock.ExpectQuery(`SELECT id, value FROM metrics WHERE name = \$1 AND type = \$2`).
				WithArgs(metric.ID, metric.MType).
				WillReturnRows(sqlmock.NewRows(nil))
			mock.ExpectExec(`INSERT INTO metrics \(name, type, value\) VALUES \(\$1, \$2, \$3\)`).
				WithArgs(metric.ID, metric.MType, *metric.Value).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}
		mock.ExpectCommit()

		err = storage.UpdateMetrics(context.Background(), metricsBatch)
		assert.NoError(t, err, "unexpected error when update metric")

		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, err, "not all queued expectations were met in order")
	})

	t.Run("Update counter metrics", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err, "unexpected error when mock db connection")
		defer db.Close()

		lg := &logger.ServerLogger{}
		storage := NewDBStorage(db, lg)

		metricsBatch := []metrics.Metrics{
			{
				ID:    "first",
				MType: "counter",
				Delta: int64Ptr(10),
			},
			{
				ID:    "first",
				MType: "counter",
				Delta: int64Ptr(20),
			},
		}

		mock.ExpectBegin()
		for _, metric := range metricsBatch {
			mock.ExpectQuery(`SELECT id, value FROM metrics WHERE name = \$1 AND type = \$2`).
				WithArgs(metric.ID, metric.MType).
				WillReturnRows(sqlmock.NewRows(nil))
			mock.ExpectExec(`INSERT INTO metrics \(name, type, value\) VALUES \(\$1, \$2, \$3\)`).
				WithArgs(metric.ID, metric.MType, floatPtr(float64(*metric.Delta))).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}
		mock.ExpectCommit()

		err = storage.UpdateMetrics(context.Background(), metricsBatch)
		assert.NoError(t, err, "unexpected error when update metric")

		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, err, "not all queued expectations were met in order")
	})
}

func TestDBStorage_GetMetric(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "unexpected error when mock db connection")
	defer db.Close()

	lg := &logger.ServerLogger{}
	storage := NewDBStorage(db, lg)

	metric := metrics.Metrics{
		ID:    "first",
		MType: "gauge",
	}
	metricValue := 10.11

	mock.ExpectQuery(`SELECT name, type, value FROM metrics WHERE name = \$1 AND type = \$2`).
		WithArgs(metric.ID, metric.MType).
		WillReturnRows(sqlmock.NewRows([]string{"name", "type", "value"}).AddRow(metric.ID, metric.MType, metricValue))

	res, err := storage.GetMetric(context.Background(), metric)
	assert.NoError(t, err, "unexpected error when get metric")
	assert.Equal(t, metric.ID, res.ID, "metric id should be %d, got %d", metric.ID, res.ID)
	assert.Equal(t, metric.MType, res.MType, "metric type should be %s, got %s", metric.MType, res.MType)
	assert.Equal(t, metricValue, *res.Value, "metric value should be %f, got %f", metricValue, res.Value)
	assert.NoError(t, mock.ExpectationsWereMet(), "not all queued expectations were met in order")
}

func TestDBStorage_GetMetrics(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "unexpected error when mock db connection")
	defer db.Close()

	lg := &logger.ServerLogger{}
	storage := NewDBStorage(db, lg)

	gaugeMetric := metrics.GaugeMetric{
		Metric: metrics.Metric{
			Name: "first",
			Type: "gauge",
		},
		Value: 10.11,
	}
	counterMetric := metrics.CounterMetric{
		Metric: metrics.Metric{
			Name: "second",
			Type: "counter",
		},
		Value: 10,
	}
	mock.ExpectQuery(`SELECT name, type, value FROM metrics WHERE type IN\(\$1, \$2\)`).
		WithArgs("gauge", "counter").
		WillReturnRows(sqlmock.NewRows([]string{"name", "type", "value"}).
			AddRow(gaugeMetric.Metric.Name, gaugeMetric.Metric.Type, gaugeMetric.Value).
			AddRow(counterMetric.Metric.Name, counterMetric.Metric.Type, counterMetric.Value))

	gauges, counters, err := storage.GetMetrics(context.Background())
	assert.NoError(t, err, "unexpected error when get metrics")
	assert.Equal(t, 1, len(gauges), "should return %d gauge metrics, got %d", 1, len(gauges))
	assert.Equal(t, 1, len(counters), "should return %d counter metrics, got %d", 1, len(counters))
	assert.Equal(t, gaugeMetric, gauges[gaugeMetric.Metric.Name], "should return gauge metric: %v, got: %v",
		gaugeMetric, gauges[gaugeMetric.Metric.Name])
	assert.Equal(t, counterMetric, counters[counterMetric.Metric.Name], "should return counter metric: %v, got: %v",
		counterMetric, counters[counterMetric.Metric.Name])
	assert.NoError(t, mock.ExpectationsWereMet(), "not all queued expectations were met in order")
}

func floatPtr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}
