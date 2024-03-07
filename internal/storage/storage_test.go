package storage

import (
	"reflect"
	"testing"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	validGaugeMetricType   = string(model.Gauge)
	validCounterMetricType = string(model.Counter)
	invalidMetricType      = "Invalid"

	defaultMetricName = "metricName"

	initGaugeMetricValue               = 12.2
	invalidGaugeMetricValue            = "eleven point one"
	validGaugeMetricValue              = "11.1"
	parsedValidGaugeMetricValue        = 11.1
	updatedNonExistingGaugeMetricValue = 11.1
	updatedExistingGaugeMetricValue    = 11.1

	initCounterMetricValue               = int64(12)
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
		name string

		args            args
		gaugesInitState map[string]model.GaugeMetric

		gaugesUpdatedState map[string]model.GaugeMetric
	}{
		{
			name: "should correct create metric when update non existing gauge metric",

			args: args{
				metricName:  defaultMetricName,
				metricValue: parsedValidGaugeMetricValue,
			},
			gaugesInitState: make(map[string]model.GaugeMetric),

			gaugesUpdatedState: map[string]model.GaugeMetric{
				defaultMetricName: model.NewGauge(defaultMetricName, updatedNonExistingGaugeMetricValue),
			},
		},

		{
			name: "should correct update metric when update existing gauge metric",

			args: args{
				metricName:  defaultMetricName,
				metricValue: parsedValidGaugeMetricValue,
			},
			gaugesInitState: map[string]model.GaugeMetric{
				defaultMetricName: model.NewGauge(defaultMetricName, initGaugeMetricValue),
			},

			gaugesUpdatedState: map[string]model.GaugeMetric{
				defaultMetricName: model.NewGauge(defaultMetricName, updatedExistingGaugeMetricValue),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			metricsStorage := MemStorage{
				gauges: tt.gaugesInitState,
			}

			err := metricsStorage.UpdateGaugeMetric(model.NewGauge(tt.args.metricName, tt.args.metricValue))
			require.NoError(t, err)

			gauges, _ := metricsStorage.GetMetrics()
			assert.True(t, reflect.DeepEqual(tt.gaugesUpdatedState, gauges))
		})
	}
}

func TestUpdateCounterMetric(t *testing.T) {
	type args struct {
		metricName  string
		metricValue int64
	}

	testCases := []struct {
		name string

		args              args
		countersInitState map[string]model.CounterMetric

		countersUpdatedState map[string]model.CounterMetric
	}{
		{
			name: "should correct create metric when update non existing counter metric",

			args: args{
				metricName:  defaultMetricName,
				metricValue: parsedValidCounterMetricValue,
			},
			countersInitState: make(map[string]model.CounterMetric),

			countersUpdatedState: map[string]model.CounterMetric{
				defaultMetricName: model.NewCounter(defaultMetricName, updatedNonExistingCounterMetricValue),
			},
		},

		{
			name: "should correct update metric when update existing counter metric",

			args: args{
				metricName:  defaultMetricName,
				metricValue: parsedValidCounterMetricValue,
			},
			countersInitState: map[string]model.CounterMetric{
				defaultMetricName: model.NewCounter(defaultMetricName, initCounterMetricValue),
			},

			countersUpdatedState: map[string]model.CounterMetric{
				defaultMetricName: model.NewCounter(defaultMetricName, updatedExistingCounterMetricValue),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			metricsStorage := MemStorage{
				counters: tt.countersInitState,
			}

			err := metricsStorage.UpdateCounterMetric(model.NewCounter(tt.args.metricName, tt.args.metricValue))
			require.NoError(t, err)

			_, counters := metricsStorage.GetMetrics()
			assert.True(t, reflect.DeepEqual(tt.countersUpdatedState, counters))
		})
	}
}

func TestUpdateMetric(t *testing.T) {
	type args struct {
		metricType  string
		metricName  string
		metricValue string
	}

	testCases := []struct {
		name string

		args              args
		gaugesInitState   map[string]model.GaugeMetric
		countersInitState map[string]model.CounterMetric

		gaugesUpdatedState   map[string]model.GaugeMetric
		countersUpdatedState map[string]model.CounterMetric
		storageError         error
	}{
		{
			name: "should return error when update metric with invalid type",

			args: args{
				metricType:  invalidMetricType,
				metricName:  defaultMetricName,
				metricValue: validCounterMetricValue,
			},
			gaugesInitState:   make(map[string]model.GaugeMetric),
			countersInitState: make(map[string]model.CounterMetric),

			gaugesUpdatedState:   make(map[string]model.GaugeMetric),
			countersUpdatedState: make(map[string]model.CounterMetric),
			storageError:         errors.InvalidMetricType{},
		},
		{
			name: "should return error when update gauge metric with invalid value",

			args: args{
				metricType:  validGaugeMetricType,
				metricName:  defaultMetricName,
				metricValue: invalidGaugeMetricValue,
			},
			gaugesInitState:   make(map[string]model.GaugeMetric),
			countersInitState: make(map[string]model.CounterMetric),

			gaugesUpdatedState:   make(map[string]model.GaugeMetric),
			countersUpdatedState: make(map[string]model.CounterMetric),
			storageError:         errors.InvalidMetricValue{},
		},

		{
			name: "should return error when update counter metric with invalid value",
			args: args{
				metricType:  validCounterMetricType,
				metricName:  defaultMetricName,
				metricValue: invalidCounterMetricValue,
			},

			gaugesInitState:   make(map[string]model.GaugeMetric),
			countersInitState: make(map[string]model.CounterMetric),

			gaugesUpdatedState:   make(map[string]model.GaugeMetric),
			countersUpdatedState: make(map[string]model.CounterMetric),
			storageError:         errors.InvalidMetricValue{},
		},

		{
			name: "should correct create metric when update non existing gauge metric with valid value",

			args: args{
				metricType:  validGaugeMetricType,
				metricName:  defaultMetricName,
				metricValue: validGaugeMetricValue,
			},
			gaugesInitState:   make(map[string]model.GaugeMetric),
			countersInitState: make(map[string]model.CounterMetric),

			gaugesUpdatedState: map[string]model.GaugeMetric{
				defaultMetricName: model.NewGauge(defaultMetricName, updatedNonExistingGaugeMetricValue),
			},
			countersUpdatedState: make(map[string]model.CounterMetric),
			storageError:         nil,
		},

		{
			name: "should correct update metric when update existing gauge metric with valid value",

			args: args{
				metricType:  validGaugeMetricType,
				metricName:  defaultMetricName,
				metricValue: validGaugeMetricValue,
			},
			gaugesInitState: map[string]model.GaugeMetric{
				defaultMetricName: model.NewGauge(defaultMetricName, initGaugeMetricValue),
			},
			countersInitState: make(map[string]model.CounterMetric),

			gaugesUpdatedState: map[string]model.GaugeMetric{
				defaultMetricName: model.NewGauge(defaultMetricName, updatedExistingGaugeMetricValue),
			},
			countersUpdatedState: make(map[string]model.CounterMetric),
			storageError:         nil,
		},

		{
			name: "should correct create metric when update non existing counter metric with valid value",

			args: args{
				metricType:  validCounterMetricType,
				metricName:  defaultMetricName,
				metricValue: validCounterMetricValue,
			},
			gaugesInitState:   make(map[string]model.GaugeMetric),
			countersInitState: make(map[string]model.CounterMetric),

			gaugesUpdatedState: make(map[string]model.GaugeMetric),
			countersUpdatedState: map[string]model.CounterMetric{
				defaultMetricName: model.NewCounter(defaultMetricName, updatedNonExistingCounterMetricValue),
			},
			storageError: nil,
		},

		{
			name: "should correct update metric when update existing counter metric with valid value",

			args: args{
				metricType:  validCounterMetricType,
				metricName:  defaultMetricName,
				metricValue: validCounterMetricValue,
			},
			gaugesInitState: make(map[string]model.GaugeMetric),
			countersInitState: map[string]model.CounterMetric{
				defaultMetricName: model.NewCounter(defaultMetricName, initCounterMetricValue),
			},

			gaugesUpdatedState: make(map[string]model.GaugeMetric),
			countersUpdatedState: map[string]model.CounterMetric{
				defaultMetricName: model.NewCounter(defaultMetricName, updatedExistingCounterMetricValue),
			},
			storageError: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			metricsStorage := MemStorage{
				gauges:   tt.gaugesInitState,
				counters: tt.countersInitState,
			}

			err := metricsStorage.UpdateMetric(tt.args.metricType, tt.args.metricName, tt.args.metricValue)
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

func TestResetMetricValue(t *testing.T) {
	type args struct {
		metricType string
		metricName string
	}

	testCases := []struct {
		name string

		args              args
		gaugesInitState   map[string]model.GaugeMetric
		countersInitState map[string]model.CounterMetric

		gaugesUpdatedState   map[string]model.GaugeMetric
		countersUpdatedState map[string]model.CounterMetric
		storageError         error
	}{
		{
			name: "should return error when reset metric with invalid type",

			args: args{
				metricType: invalidMetricType,
				metricName: defaultMetricName,
			},
			gaugesInitState:   make(map[string]model.GaugeMetric),
			countersInitState: make(map[string]model.CounterMetric),

			gaugesUpdatedState:   make(map[string]model.GaugeMetric),
			countersUpdatedState: make(map[string]model.CounterMetric),
			storageError:         errors.InvalidMetricType{},
		},
		{
			name: "should return error when reset non existing gauge metric",

			args: args{
				metricType: validGaugeMetricType,
				metricName: defaultMetricName,
			},
			gaugesInitState:   make(map[string]model.GaugeMetric),
			countersInitState: make(map[string]model.CounterMetric),

			gaugesUpdatedState:   make(map[string]model.GaugeMetric),
			countersUpdatedState: make(map[string]model.CounterMetric),
			storageError:         errors.InvalidMetricName{},
		},

		{
			name: "should correct reset metric when reset existing gauge metric",

			args: args{
				metricType: validGaugeMetricType,
				metricName: defaultMetricName,
			},
			gaugesInitState: map[string]model.GaugeMetric{
				defaultMetricName: model.NewGauge(defaultMetricName, initGaugeMetricValue),
			},
			countersInitState: make(map[string]model.CounterMetric),

			gaugesUpdatedState: map[string]model.GaugeMetric{
				defaultMetricName: model.NewGauge(defaultMetricName, resetedGaugeMetricValue),
			},
			countersUpdatedState: make(map[string]model.CounterMetric),
			storageError:         nil,
		},

		{
			name: "should return error when reset non existing counter metric",

			args: args{
				metricType: validCounterMetricType,
				metricName: defaultMetricName,
			},
			gaugesInitState:   make(map[string]model.GaugeMetric),
			countersInitState: make(map[string]model.CounterMetric),

			gaugesUpdatedState:   make(map[string]model.GaugeMetric),
			countersUpdatedState: make(map[string]model.CounterMetric),
			storageError:         errors.InvalidMetricName{},
		},

		{
			name: "should correct reset metric when reset existing counter metric",

			args: args{
				metricType: validCounterMetricType,
				metricName: defaultMetricName,
			},
			gaugesInitState: make(map[string]model.GaugeMetric),
			countersInitState: map[string]model.CounterMetric{
				defaultMetricName: model.NewCounter(defaultMetricName, initCounterMetricValue),
			},

			gaugesUpdatedState: make(map[string]model.GaugeMetric),
			countersUpdatedState: map[string]model.CounterMetric{
				defaultMetricName: model.NewCounter(defaultMetricName, resetedCounterMetricValue),
			},
			storageError: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			metricsStorage := MemStorage{
				gauges:   tt.gaugesInitState,
				counters: tt.countersInitState,
			}

			err := metricsStorage.ResetMetricValue(tt.args.metricType, tt.args.metricName)
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
