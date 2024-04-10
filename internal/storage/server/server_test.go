package server

// import (
// 	"os"
// 	"testing"

// 	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
// 	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func TestSaveAndLoadGaugeMetrics(t *testing.T) {
// 	testCases := []struct {
// 		name                 string
// 		metric               metrics.GaugeMetric
// 		initStateMetricValue float64
// 		gaugesInitState      map[string]metrics.GaugeMetric
// 	}{{
// 		name:                 "should success save and load gauge metric using file storage #1",
// 		metric:               metrics.NewGauge("first", 11.12),
// 		initStateMetricValue: 0,
// 		gaugesInitState: map[string]metrics.GaugeMetric{
// 			"first": metrics.NewGauge("first", 0),
// 		},
// 	}, {
// 		name:                 "should succes save and load gauge metric useing file storage #2",
// 		metric:               metrics.NewGauge("second", 22.11),
// 		initStateMetricValue: 0,
// 		gaugesInitState: map[string]metrics.GaugeMetric{
// 			"second": metrics.NewGauge("second", 0),
// 		},
// 	},
// 	}

// 	for _, tt := range testCases {
// 		t.Run(tt.name, func(t *testing.T) {
// 			file, err := os.CreateTemp(t.TempDir(), "file-storage-*.json")
// 			require.NoError(t, err)

// 			logger, err := logger.Initialize("info")
// 			require.NoError(t, err, "Error init logger")

// 			storage := NewServerMemStorage(logger)
// 			storage.SetGaugeMetircs(tt.gaugesInitState)

// 			initMetric, err := storage.GetGaugeMetric(tt.metric.Name)
// 			require.NoError(t, err)
// 			assert.Equal(t, tt.initStateMetricValue, initMetric.GetValue())

// 			updatedMetric, err := storage.UpdateGaugeMetric(tt.metric)
// 			require.NoError(t, err)
// 			assert.Equal(t, tt.metric.Value, updatedMetric.GetValue())

// 			err = storage.Save(file.Name())
// 			require.NoError(t, err)

// 			storage.SetGaugeMetircs(map[string]metrics.GaugeMetric{})
// 			_, err = storage.GetGaugeMetric(tt.metric.Name)
// 			require.Error(t, err)

// 			err = storage.Restore(file.Name())
// 			require.NoError(t, err)

// 			restoredFromFileStorageMetric, err := storage.GetGaugeMetric(tt.metric.Name)
// 			require.NoError(t, err)
// 			assert.Equal(t, tt.metric.Value, restoredFromFileStorageMetric.Value)
// 		})
// 	}
// }

// func TestSaveAndLoadCounterMetrics(t *testing.T) {
// 	testCases := []struct {
// 		name                 string
// 		metric               metrics.CounterMetric
// 		initStateMetricValue int64
// 		countersInitState    map[string]metrics.CounterMetric
// 	}{{
// 		name:                 "should success save and load counter metric using file storage #1",
// 		metric:               metrics.NewCounter("first", 11),
// 		initStateMetricValue: 0,
// 		countersInitState: map[string]metrics.CounterMetric{
// 			"first": metrics.NewCounter("first", 0),
// 		},
// 	}, {
// 		name:                 "should succes save and load counter metric using file storage #2",
// 		metric:               metrics.NewCounter("second", 22),
// 		initStateMetricValue: 0,
// 		countersInitState: map[string]metrics.CounterMetric{
// 			"second": metrics.NewCounter("second", 0),
// 		},
// 	},
// 	}

// 	for _, tt := range testCases {
// 		t.Run(tt.name, func(t *testing.T) {
// 			file, err := os.CreateTemp(t.TempDir(), "file-storage-*.json")
// 			require.NoError(t, err)

// 			logger, err := logger.Initialize("info")
// 			require.NoError(t, err, "Error init logger")

// 			storage := NewServerMemStorage(logger)
// 			storage.SetCounterMetrics(tt.countersInitState)

// 			initMetric, err := storage.GetCounterMetric(tt.metric.Name)
// 			require.NoError(t, err)
// 			assert.Equal(t, tt.initStateMetricValue, initMetric.GetValue())

// 			updatedMetric, err := storage.UpdateCounterMetric(tt.metric)
// 			require.NoError(t, err)
// 			assert.Equal(t, tt.metric.Value, updatedMetric.GetValue())

// 			err = storage.Save(file.Name())
// 			require.NoError(t, err)

// 			storage.SetCounterMetrics(map[string]metrics.CounterMetric{})
// 			_, err = storage.GetCounterMetric(tt.metric.Name)
// 			require.Error(t, err)

// 			err = storage.Restore(file.Name())
// 			require.NoError(t, err)

// 			restoredFromFileStorageMetric, err := storage.GetCounterMetric(tt.metric.Name)
// 			require.NoError(t, err)
// 			assert.Equal(t, tt.metric.Value, restoredFromFileStorageMetric.Value)
// 		})
// 	}
// }
