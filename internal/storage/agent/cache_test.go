package storage

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	er "github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
)

const (
	validGaugeMetricType   = string(metrics.Gauge)
	validCounterMetricType = string(metrics.Counter)
	invalidMetricType      = "Invalid"

	defaultMetricName = "metricName"

	initGaugeMetricValue               = 12.2
	parsedInitGaugeMetricValue         = "12.2"
	invalidGaugeMetricValue            = "eleven point one"
	validGaugeMetricValue              = "11.1"
	parsedValidGaugeMetricValue        = 11.1
	updatedNonExistingGaugeMetricValue = 11.1
	updatedExistingGaugeMetricValue    = 11.1

	initCounterMetricValue               = int64(12)
	parsedInitCounterMetricValue         = "12"
	invalidCounterMetricValue            = "11.1"
	validCounterMetricValue              = "11"
	parsedValidCounterMetricValue        = 11
	updatedNonExistingCounterMetricValue = int64(11)
	updatedExistingCounterMetricValue    = int64(23)

	resetedGaugeMetricValue   = 0.0
	resetedCounterMetricValue = 0
)

func TestUpdateGaugeMetric(t *testing.T) {
	type args struct {
		metricName  string
		metricValue float64
	}

	testCases := []struct {
		name               string
		args               args
		gaugesInitState    map[string]metrics.GaugeMetric
		gaugesUpdatedState map[string]metrics.GaugeMetric
		err                error
	}{
		{
			name: "should return error when update non existing gauge metric",
			args: args{
				metricName:  defaultMetricName,
				metricValue: parsedValidGaugeMetricValue,
			},
			gaugesInitState: make(map[string]metrics.GaugeMetric),
			err:             er.InvalidMetricName{},
		},

		{
			name: "should correct update metric when update existing gauge metric",
			args: args{
				metricName:  defaultMetricName,
				metricValue: parsedValidGaugeMetricValue,
			},
			gaugesInitState: map[string]metrics.GaugeMetric{
				defaultMetricName: metrics.NewGauge(defaultMetricName, initGaugeMetricValue),
			},
			gaugesUpdatedState: map[string]metrics.GaugeMetric{
				defaultMetricName: metrics.NewGauge(defaultMetricName, updatedExistingGaugeMetricValue),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := logger.Initialize("info")
			require.NoError(t, err, "Error init logger")
			metricsStorage := NewAgentMemCache(tt.gaugesInitState, make(map[string]metrics.CounterMetric), logger)

			_, err = metricsStorage.UpdateGaugeMetric(metrics.NewGauge(tt.args.metricName, tt.args.metricValue))
			if tt.err == nil {
				require.NoError(t, err)
				gauges, _ := metricsStorage.GetMetrics()
				assert.True(t, reflect.DeepEqual(tt.gaugesUpdatedState, gauges))
			} else {
				require.Error(t, err)
				assert.IsType(t, tt.err, err)
			}
		})
	}
}

func TestUpdateCounterMetric(t *testing.T) {
	type args struct {
		metricName  string
		metricValue int64
	}

	testCases := []struct {
		name                 string
		args                 args
		countersInitState    map[string]metrics.CounterMetric
		countersUpdatedState map[string]metrics.CounterMetric
		err                  error
	}{
		{
			name: "should return error when update non existing counter metric",
			args: args{
				metricName:  defaultMetricName,
				metricValue: parsedValidCounterMetricValue,
			},
			countersInitState: make(map[string]metrics.CounterMetric),
			err:               er.InvalidMetricName{},
		},

		{
			name: "should correct update metric when update existing counter metric",
			args: args{
				metricName:  defaultMetricName,
				metricValue: parsedValidCounterMetricValue,
			},
			countersInitState: map[string]metrics.CounterMetric{
				defaultMetricName: metrics.NewCounter(defaultMetricName, initCounterMetricValue),
			},
			countersUpdatedState: map[string]metrics.CounterMetric{
				defaultMetricName: metrics.NewCounter(defaultMetricName, updatedExistingCounterMetricValue),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := logger.Initialize("info")
			require.NoError(t, err, "Error init logger")
			metricsStorage := NewAgentMemCache(make(map[string]metrics.GaugeMetric), tt.countersInitState, logger)

			_, err = metricsStorage.UpdateCounterMetric(metrics.NewCounter(tt.args.metricName, tt.args.metricValue))

			if tt.err == nil {
				require.NoError(t, err)
				_, counters := metricsStorage.GetMetrics()
				assert.True(t, reflect.DeepEqual(tt.countersUpdatedState, counters))
			} else {
				require.Error(t, err)
				assert.IsType(t, tt.err, err)
			}
		})
	}
}

func TestUpdateMonitorMetrics(t *testing.T) {
	type state struct {
		gauges   map[string]metrics.GaugeMetric
		counters map[string]metrics.CounterMetric
	}

	testCases := []struct {
		name          string
		monitor       *monitors.Monitor
		expectedState state
	}{{
		name: "should correct update monitor metrics and related counters",
		monitor: &monitors.Monitor{
			Alloc:         1.0,
			BuckHashSys:   2.0,
			Frees:         3.0,
			GCCPUFraction: 4.0,
			GCSys:         5.0,
			HeapAlloc:     6.0,
			HeapIdle:      7.0,
			HeapInuse:     8.0,
			HeapObjects:   9.0,
			HeapReleased:  10.0,
			HeapSys:       11.0,
			LastGC:        12.0,
			Lookups:       13.0,
			MCacheInuse:   14.0,
			MCacheSys:     15.0,
			MSpanInuse:    16.0,
			MSpanSys:      17.0,
			Mallocs:       18.0,
			NextGC:        19.0,
			NumForcedGC:   20.0,
			NumGC:         21.0,
			OtherSys:      22.0,
			PauseTotalNs:  23.0,
			StackInuse:    24.0,
			StackSys:      25.0,
			Sys:           26.0,
			TotalAlloc:    27.0,
		},
		expectedState: state{
			gauges: map[string]metrics.GaugeMetric{
				"Alloc":         metrics.NewGauge("Alloc", 1.0),
				"BuckHashSys":   metrics.NewGauge("BuckHashSys", 2.0),
				"Frees":         metrics.NewGauge("Frees", 3.0),
				"GCCPUFraction": metrics.NewGauge("GCCPUFraction", 4.0),
				"GCSys":         metrics.NewGauge("GCSys", 5.0),
				"HeapAlloc":     metrics.NewGauge("HeapAlloc", 6.0),
				"HeapIdle":      metrics.NewGauge("HeapIdle", 7.0),
				"HeapInuse":     metrics.NewGauge("HeapInuse", 8.0),
				"HeapObjects":   metrics.NewGauge("HeapObjects", 9.0),
				"HeapReleased":  metrics.NewGauge("HeapReleased", 10.0),
				"HeapSys":       metrics.NewGauge("HeapSys", 11.0),
				"LastGC":        metrics.NewGauge("LastGC", 12.0),
				"Lookups":       metrics.NewGauge("Lookups", 13.0),
				"MCacheInuse":   metrics.NewGauge("MCacheInuse", 14.0),
				"MCacheSys":     metrics.NewGauge("MCacheSys", 15.0),
				"MSpanInuse":    metrics.NewGauge("MSpanInuse", 16.0),
				"MSpanSys":      metrics.NewGauge("MSpanSys", 17.0),
				"Mallocs":       metrics.NewGauge("Mallocs", 18.0),
				"NextGC":        metrics.NewGauge("NextGC", 19.0),
				"NumForcedGC":   metrics.NewGauge("NumForcedGC", 20.0),
				"NumGC":         metrics.NewGauge("NumGC", 21.0),
				"OtherSys":      metrics.NewGauge("OtherSys", 22.0),
				"PauseTotalNs":  metrics.NewGauge("PauseTotalNs", 23.0),
				"StackInuse":    metrics.NewGauge("StackInuse", 24.0),
				"StackSys":      metrics.NewGauge("StackSys", 25.0),
				"Sys":           metrics.NewGauge("Sys", 26.0),
				"TotalAlloc":    metrics.NewGauge("TotalAlloc", 27.0),
			},
			counters: map[string]metrics.CounterMetric{
				"PollCount": metrics.NewCounter("PollCount", 27),
			},
		},
	},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := logger.Initialize("info")
			require.NoError(t, err, "Error init logger")
			storage := NewAgentMemCache(metrics.SupportedGaugeMetrics, metrics.SupportedCounterMetrics, logger)
			storage.UpdateMonitorMetrics(tt.monitor)

			gauges, counters := storage.GetMetrics()

			for key, value := range tt.expectedState.gauges {
				assert.Equal(t, value, gauges[key])
			}

			for key, value := range tt.expectedState.counters {
				assert.Equal(t, value, counters[key])
			}
		})
	}
}

func TestResetMetricValue(t *testing.T) {
	type args struct {
		metricType string
		metricName string
	}

	testCases := []struct {
		name string

		args              args
		gaugesInitState   map[string]metrics.GaugeMetric
		countersInitState map[string]metrics.CounterMetric

		gaugesUpdatedState   map[string]metrics.GaugeMetric
		countersUpdatedState map[string]metrics.CounterMetric
		storageError         error
	}{
		{
			name: "should return error when reset metric with invalid type",

			args: args{
				metricType: invalidMetricType,
				metricName: defaultMetricName,
			},
			gaugesInitState:   make(map[string]metrics.GaugeMetric),
			countersInitState: make(map[string]metrics.CounterMetric),

			gaugesUpdatedState:   make(map[string]metrics.GaugeMetric),
			countersUpdatedState: make(map[string]metrics.CounterMetric),
			storageError:         er.InvalidMetricType{},
		},
		{
			name: "should return error when reset non existing gauge metric",

			args: args{
				metricType: validGaugeMetricType,
				metricName: defaultMetricName,
			},
			gaugesInitState:   make(map[string]metrics.GaugeMetric),
			countersInitState: make(map[string]metrics.CounterMetric),

			gaugesUpdatedState:   make(map[string]metrics.GaugeMetric),
			countersUpdatedState: make(map[string]metrics.CounterMetric),
			storageError:         er.InvalidMetricName{},
		},

		{
			name: "should correct reset metric when reset existing gauge metric",

			args: args{
				metricType: validGaugeMetricType,
				metricName: defaultMetricName,
			},
			gaugesInitState: map[string]metrics.GaugeMetric{
				defaultMetricName: metrics.NewGauge(defaultMetricName, initGaugeMetricValue),
			},
			countersInitState: make(map[string]metrics.CounterMetric),

			gaugesUpdatedState: map[string]metrics.GaugeMetric{
				defaultMetricName: metrics.NewGauge(defaultMetricName, resetedGaugeMetricValue),
			},
			countersUpdatedState: make(map[string]metrics.CounterMetric),
			storageError:         nil,
		},

		{
			name: "should return error when reset non existing counter metric",

			args: args{
				metricType: validCounterMetricType,
				metricName: defaultMetricName,
			},
			gaugesInitState:   make(map[string]metrics.GaugeMetric),
			countersInitState: make(map[string]metrics.CounterMetric),

			gaugesUpdatedState:   make(map[string]metrics.GaugeMetric),
			countersUpdatedState: make(map[string]metrics.CounterMetric),
			storageError:         er.InvalidMetricName{},
		},

		{
			name: "should correct reset metric when reset existing counter metric",

			args: args{
				metricType: validCounterMetricType,
				metricName: defaultMetricName,
			},
			gaugesInitState: make(map[string]metrics.GaugeMetric),
			countersInitState: map[string]metrics.CounterMetric{
				defaultMetricName: metrics.NewCounter(defaultMetricName, initCounterMetricValue),
			},

			gaugesUpdatedState: make(map[string]metrics.GaugeMetric),
			countersUpdatedState: map[string]metrics.CounterMetric{
				defaultMetricName: metrics.NewCounter(defaultMetricName, resetedCounterMetricValue),
			},
			storageError: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := logger.Initialize("info")
			require.NoError(t, err, "Error init logger")
			metricsStorage := NewAgentMemCache(tt.gaugesInitState, tt.countersInitState, logger)

			err = metricsStorage.ResetMetricValue(tt.args.metricType, tt.args.metricName)
			if tt.storageError != nil {
				assert.IsType(t, tt.storageError, err)
			} else {
				require.NoError(t, err)
			}

			gauges, counters := metricsStorage.GetMetrics()

			assert.True(t, reflect.DeepEqual(tt.gaugesUpdatedState, gauges))
			assert.True(t, reflect.DeepEqual(tt.countersUpdatedState, counters))
		})
	}
}
