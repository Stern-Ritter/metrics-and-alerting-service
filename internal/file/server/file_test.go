package server

import (
	"os"
	"testing"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveAndLoadGaugeMetrics(t *testing.T) {
	testCases := []struct {
		name                 string
		metric               metrics.GaugeMetric
		initStateMetricValue float64
		gaugesInitState      map[string]metrics.GaugeMetric
	}{{
		name:                 "should success save and load gauge metric using file storage #1",
		metric:               metrics.NewGauge("first", 11.12),
		initStateMetricValue: 0,
		gaugesInitState: map[string]metrics.GaugeMetric{
			"first": metrics.NewGauge("first", 0),
		},
	}, {
		name:                 "should succes save and load gauge metric useing file storage #2",
		metric:               metrics.NewGauge("second", 22.11),
		initStateMetricValue: 0,
		gaugesInitState: map[string]metrics.GaugeMetric{
			"second": metrics.NewGauge("second", 0),
		},
	},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.CreateTemp(t.TempDir(), "file-storage-*.json")
			require.NoError(t, err)

			storage := storage.NewServerMemStorage()
			storage.SetGaugeMetircs(tt.gaugesInitState)

			fileStorage, err := NewServerFileStorage(file.Name(), &storage)
			require.NoError(t, err)

			initMetric, err := storage.GetGaugeMetric(tt.metric.Name)
			require.NoError(t, err)
			assert.Equal(t, tt.initStateMetricValue, initMetric.GetValue())

			updatedMetric, err := storage.UpdateGaugeMetric(tt.metric)
			require.NoError(t, err)
			assert.Equal(t, tt.metric.Value, updatedMetric.GetValue())

			err = fileStorage.Save()
			require.NoError(t, err)

			storage.SetGaugeMetircs(map[string]metrics.GaugeMetric{})
			_, err = storage.GetGaugeMetric(tt.metric.Name)
			require.Error(t, err)

			err = fileStorage.Load()
			require.NoError(t, err)

			restoredFromFileStorageMetric, err := storage.UpdateGaugeMetric(tt.metric)
			require.NoError(t, err)
			assert.Equal(t, tt.metric.Value, restoredFromFileStorageMetric.Value)
		})
	}
}

func TestSaveAndLoadCounterMetrics(t *testing.T) {
	testCases := []struct {
		name                 string
		metric               metrics.CounterMetric
		initStateMetricValue int64
		countersInitState    map[string]metrics.CounterMetric
	}{{
		name:                 "should success save and load counter metric using file storage #1",
		metric:               metrics.NewCounter("first", 11),
		initStateMetricValue: 0,
		countersInitState: map[string]metrics.CounterMetric{
			"first": metrics.NewCounter("first", 0),
		},
	}, {
		name:                 "should succes save and load gauge metric useing file storage #2",
		metric:               metrics.NewCounter("second", 22),
		initStateMetricValue: 0,
		countersInitState: map[string]metrics.CounterMetric{
			"second": metrics.NewCounter("second", 0),
		},
	},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.CreateTemp(t.TempDir(), "file-storage-*.json")
			require.NoError(t, err)

			storage := storage.ServerMemStorage{}
			storage.SetCounterMetrics(tt.countersInitState)

			fileStorage, err := NewServerFileStorage(file.Name(), &storage)
			require.NoError(t, err)

			initMetric, err := storage.GetCounterMetric(tt.metric.Name)
			require.NoError(t, err)
			assert.Equal(t, tt.initStateMetricValue, initMetric.GetValue())

			updatedMetric, err := storage.UpdateCounterMetric(tt.metric)
			require.NoError(t, err)
			assert.Equal(t, tt.metric.Value, updatedMetric.GetValue())

			err = fileStorage.Save()
			require.NoError(t, err)

			storage.SetCounterMetrics(map[string]metrics.CounterMetric{})
			_, err = storage.GetCounterMetric(tt.metric.Name)
			require.Error(t, err)

			err = fileStorage.Load()
			require.NoError(t, err)

			restoredFromFileStorageMetric, err := storage.UpdateCounterMetric(tt.metric)
			require.NoError(t, err)
			assert.Equal(t, tt.metric.Value, restoredFromFileStorageMetric.Value)
		})
	}
}

func TestClose(t *testing.T) {
	type want struct {
		gaugeMetric       metrics.GaugeMetric
		gaugeMetricError  bool
		counterMetric     metrics.CounterMetric
		couterMetricError bool
	}

	testCases := []struct {
		name              string
		gaugesInitState   map[string]metrics.GaugeMetric
		countersInitState map[string]metrics.CounterMetric
		want              want
	}{{
		name:              "should success save metrics on file storage close #1",
		gaugesInitState:   map[string]metrics.GaugeMetric{"first": metrics.NewGauge("first", 11.11)},
		countersInitState: map[string]metrics.CounterMetric{"first": metrics.NewCounter("first", 12)},
		want: want{
			gaugeMetric:       metrics.NewGauge("first", 11.11),
			gaugeMetricError:  false,
			counterMetric:     metrics.NewCounter("first", 12),
			couterMetricError: false,
		},
	}, {
		name:              "should success save metrics on file storage close #2",
		gaugesInitState:   make(map[string]metrics.GaugeMetric),
		countersInitState: make(map[string]metrics.CounterMetric),
		want: want{
			gaugeMetric:       metrics.GaugeMetric{},
			gaugeMetricError:  true,
			counterMetric:     metrics.CounterMetric{},
			couterMetricError: true,
		},
	},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			f1, err := os.CreateTemp(t.TempDir(), "file-storage-*.json")
			require.NoError(t, err)

			s1 := storage.NewServerMemStorage()
			s1.SetGaugeMetircs(tt.gaugesInitState)
			s1.SetCounterMetrics(tt.countersInitState)

			fs, err := NewServerFileStorage(f1.Name(), &s1)
			require.NoError(t, err)

			initGaugeMetric1, err := s1.GetGaugeMetric(tt.want.gaugeMetric.Name)
			if tt.want.gaugeMetricError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want.gaugeMetric.Value, initGaugeMetric1.GetValue())
			}

			initCounterMetric1, err := s1.GetCounterMetric(tt.want.counterMetric.Name)
			if tt.want.couterMetricError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want.counterMetric.Value, initCounterMetric1.GetValue())
			}

			fs.Close()

			f2, err := os.CreateTemp(t.TempDir(), "file-storage-*.json")
			require.NoError(t, err)

			s2 := storage.NewServerMemStorage()
			s2.SetGaugeMetircs(tt.gaugesInitState)
			s2.SetCounterMetrics(tt.countersInitState)

			_, err = NewServerFileStorage(f2.Name(), &s2)
			require.NoError(t, err)

			initGaugeMetric2, err := s2.GetGaugeMetric(tt.want.gaugeMetric.Name)
			if tt.want.gaugeMetricError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want.gaugeMetric.Value, initGaugeMetric2.GetValue())
			}

			initCounterMetric2, err := s2.GetCounterMetric(tt.want.counterMetric.Name)
			if tt.want.couterMetricError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want.counterMetric.Value, initCounterMetric2.GetValue())
			}
		})
	}

}
