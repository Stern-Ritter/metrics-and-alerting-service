package server

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	validGaugeMetricType   = string(metrics.Gauge)
	validCounterMetricType = string(metrics.Counter)
	defaultMetricName      = "metricName"

	initGaugeMetricValue               = float64(12.2)
	validGaugeMetricValue              = float64(11.1)
	updatedNonExistingGaugeMetricValue = float64(11.1)
	updatedExistingGaugeMetricValue    = float64(11.1)

	initCounterMetricValue               = int64(12)
	validCounterMetricValue              = int64(11)
	updatedNonExistingCounterMetricValue = int64(11)
	updatedExistingCounterMetricValue    = int64(23)
)

func TestUpdateMetric(t *testing.T) {
	testCases := []struct {
		name string

		gaugesInitState   map[string]metrics.GaugeMetric
		countersInitState map[string]metrics.CounterMetric

		metric metrics.Metrics

		gaugesUpdatedState   map[string]metrics.GaugeMetric
		countersUpdatedState map[string]metrics.CounterMetric
		storageError         error
	}{
		{
			name:              "should return error when update metric with invalid type",
			gaugesInitState:   make(map[string]metrics.GaugeMetric),
			countersInitState: make(map[string]metrics.CounterMetric),

			metric: metrics.Metrics{ID: defaultMetricName, MType: "unknown", Value: &validGaugeMetricValue},

			gaugesUpdatedState:   make(map[string]metrics.GaugeMetric),
			countersUpdatedState: make(map[string]metrics.CounterMetric),
			storageError:         errors.InvalidMetricType{},
		},
		{
			name:              "should correct create metric when update non existing gauge metric with valid value",
			gaugesInitState:   make(map[string]metrics.GaugeMetric),
			countersInitState: make(map[string]metrics.CounterMetric),

			metric: metrics.Metrics{ID: defaultMetricName, MType: validGaugeMetricType, Value: &validGaugeMetricValue},

			gaugesUpdatedState: map[string]metrics.GaugeMetric{defaultMetricName: metrics.NewGauge(defaultMetricName,
				updatedNonExistingGaugeMetricValue)},
			countersUpdatedState: make(map[string]metrics.CounterMetric),
			storageError:         nil,
		},
		{
			name:              "should correct update metric when update existing gauge metric with valid value",
			gaugesInitState:   map[string]metrics.GaugeMetric{defaultMetricName: metrics.NewGauge(defaultMetricName, initGaugeMetricValue)},
			countersInitState: make(map[string]metrics.CounterMetric),

			metric: metrics.Metrics{ID: defaultMetricName, MType: validGaugeMetricType, Value: &validGaugeMetricValue},

			gaugesUpdatedState: map[string]metrics.GaugeMetric{defaultMetricName: metrics.NewGauge(defaultMetricName,
				updatedExistingGaugeMetricValue)},
			countersUpdatedState: make(map[string]metrics.CounterMetric),
			storageError:         nil,
		},

		{
			name:              "should correct create metric when update non existing counter metric with valid value",
			gaugesInitState:   make(map[string]metrics.GaugeMetric),
			countersInitState: make(map[string]metrics.CounterMetric),

			metric: metrics.Metrics{ID: defaultMetricName, MType: validCounterMetricType, Delta: &validCounterMetricValue},

			gaugesUpdatedState: make(map[string]metrics.GaugeMetric),
			countersUpdatedState: map[string]metrics.CounterMetric{
				defaultMetricName: metrics.NewCounter(defaultMetricName, updatedNonExistingCounterMetricValue),
			},
			storageError: nil,
		},

		{
			name:              "should correct update metric when update existing counter metric with valid value",
			gaugesInitState:   make(map[string]metrics.GaugeMetric),
			countersInitState: map[string]metrics.CounterMetric{defaultMetricName: metrics.NewCounter(defaultMetricName, initCounterMetricValue)},

			metric: metrics.Metrics{ID: defaultMetricName, MType: validCounterMetricType, Delta: &validCounterMetricValue},

			gaugesUpdatedState: make(map[string]metrics.GaugeMetric),
			countersUpdatedState: map[string]metrics.CounterMetric{
				defaultMetricName: metrics.NewCounter(defaultMetricName, updatedExistingCounterMetricValue),
			},
			storageError: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := logger.Initialize("info")
			require.NoError(t, err, "Error init logger")
			metricsStorage := NewMemoryStorage(logger)
			metricsStorage.SetGaugeMetircs(tt.gaugesInitState)
			metricsStorage.SetCounterMetrics(tt.countersInitState)

			err = metricsStorage.UpdateMetric(context.TODO(), tt.metric)
			if tt.storageError != nil {
				require.Error(t, err)
				assert.IsType(t, tt.storageError, err)
			} else {
				require.NoError(t, err)
			}

			gauges, counters, err := metricsStorage.GetMetrics(context.TODO())
			require.NoError(t, err)
			assert.True(t, reflect.DeepEqual(tt.gaugesUpdatedState, gauges))
			assert.True(t, reflect.DeepEqual(tt.countersUpdatedState, counters))
		})
	}
}

func TestGetMetric(t *testing.T) {
	testCases := []struct {
		name              string
		gaugesInitState   map[string]metrics.GaugeMetric
		countersInitState map[string]metrics.CounterMetric
		requestedMetric   metrics.Metrics
		returnedMetric    metrics.Metrics
		err               error
	}{
		{
			name:              "should correct return gauge metric when get existing gauge metric",
			gaugesInitState:   map[string]metrics.GaugeMetric{"first": metrics.NewGauge("first", initGaugeMetricValue)},
			countersInitState: make(map[string]metrics.CounterMetric),
			requestedMetric:   metrics.Metrics{ID: "first", MType: string(metrics.Gauge)},
			returnedMetric:    metrics.Metrics{ID: "first", MType: string(metrics.Gauge), Value: &initGaugeMetricValue},
			err:               nil,
		},
		{
			name:              "should return error when get non existing gauge metric",
			gaugesInitState:   map[string]metrics.GaugeMetric{"first": metrics.NewGauge("first", 64)},
			countersInitState: make(map[string]metrics.CounterMetric),
			requestedMetric:   metrics.Metrics{ID: "second", MType: string(metrics.Gauge)},
			err:               errors.InvalidMetricName{},
		},
		{
			name:              "should correct return counter metric when get existing counter metric",
			gaugesInitState:   make(map[string]metrics.GaugeMetric),
			countersInitState: map[string]metrics.CounterMetric{"first": metrics.NewCounter("first", initCounterMetricValue)},
			requestedMetric:   metrics.Metrics{ID: "first", MType: string(metrics.Counter)},
			returnedMetric:    metrics.Metrics{ID: "first", MType: string(metrics.Counter), Delta: &initCounterMetricValue},
			err:               nil,
		},

		{
			name:              "should return error when get non existing gauge metric",
			gaugesInitState:   make(map[string]metrics.GaugeMetric),
			countersInitState: map[string]metrics.CounterMetric{"first": metrics.NewCounter("first", 64)},
			requestedMetric:   metrics.Metrics{ID: "second", MType: string(metrics.Counter)},
			err:               errors.InvalidMetricName{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := logger.Initialize("info")
			require.NoError(t, err, "Error init logger")
			metricsStorage := NewMemoryStorage(logger)
			metricsStorage.SetGaugeMetircs(tt.gaugesInitState)
			metricsStorage.SetCounterMetrics(tt.countersInitState)

			metric, err := metricsStorage.GetMetric(context.TODO(), tt.requestedMetric)
			if tt.err != nil {
				require.Error(t, err)
				assert.IsType(t, tt.err, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.returnedMetric, metric)
			}
		})
	}
}

func TestSaveAndLoadMetrics(t *testing.T) {
	testCases := []struct {
		name                     string
		gaugesInitState          map[string]metrics.GaugeMetric
		countersInitState        map[string]metrics.CounterMetric
		requestedMetric          metrics.Metrics
		udpatedMetric            metrics.Metrics
		initStateMetricValue     float64
		updatedStateMetricValue  float64
		restoredStateMetricValue float64
	}{{
		name:                     "should success save and load gauge metric using file storage",
		gaugesInitState:          map[string]metrics.GaugeMetric{"first": metrics.NewGauge("first", 0)},
		countersInitState:        make(map[string]metrics.CounterMetric),
		requestedMetric:          metrics.Metrics{ID: "first", MType: string(metrics.Gauge)},
		udpatedMetric:            metrics.Metrics{ID: "first", MType: string(metrics.Gauge), Value: &initGaugeMetricValue},
		initStateMetricValue:     0,
		updatedStateMetricValue:  initGaugeMetricValue,
		restoredStateMetricValue: initGaugeMetricValue,
	}, {
		name:                     "should success save and load counter metric using file storage",
		gaugesInitState:          make(map[string]metrics.GaugeMetric),
		countersInitState:        map[string]metrics.CounterMetric{"first": metrics.NewCounter("first", initCounterMetricValue)},
		requestedMetric:          metrics.Metrics{ID: "first", MType: string(metrics.Counter)},
		udpatedMetric:            metrics.Metrics{ID: "first", MType: string(metrics.Counter), Delta: &initCounterMetricValue},
		initStateMetricValue:     float64(initCounterMetricValue),
		updatedStateMetricValue:  float64(initCounterMetricValue + initCounterMetricValue),
		restoredStateMetricValue: float64(initCounterMetricValue + initCounterMetricValue),
	}}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.CreateTemp(t.TempDir(), "file-storage-*.json")
			require.NoError(t, err)

			logger, err := logger.Initialize("info")
			require.NoError(t, err, "Error init logger")
			metricsStorage := NewMemoryStorage(logger)
			metricsStorage.SetGaugeMetircs(tt.gaugesInitState)
			metricsStorage.SetCounterMetrics(tt.countersInitState)

			initMetric, err := metricsStorage.GetMetric(context.TODO(), tt.requestedMetric)
			require.NoError(t, err)
			value, err := initMetric.GetValue()
			require.NoError(t, err)
			assert.Equal(t, tt.initStateMetricValue, value)

			err = metricsStorage.UpdateMetric(context.TODO(), tt.udpatedMetric)
			require.NoError(t, err)
			updatedMetric, err := metricsStorage.GetMetric(context.TODO(), tt.requestedMetric)
			require.NoError(t, err)
			value, err = updatedMetric.GetValue()
			require.NoError(t, err)
			assert.Equal(t, tt.updatedStateMetricValue, value)

			err = metricsStorage.Save(file.Name())
			require.NoError(t, err)

			metricsStorage.SetGaugeMetircs(map[string]metrics.GaugeMetric{})
			metricsStorage.SetCounterMetrics(map[string]metrics.CounterMetric{})
			_, err = metricsStorage.GetMetric(context.TODO(), tt.requestedMetric)
			require.Error(t, err)

			err = metricsStorage.Restore(file.Name())
			require.NoError(t, err)

			restoredFromFileStorageMetric, err := metricsStorage.GetMetric(context.TODO(), tt.requestedMetric)
			require.NoError(t, err)
			value, err = restoredFromFileStorageMetric.GetValue()
			require.NoError(t, err)
			assert.Equal(t, tt.restoredStateMetricValue, value)
		})
	}
}
