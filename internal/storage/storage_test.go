package storage

import (
	"reflect"
	"testing"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		name string

		args            args
		gaugesInitState map[string]metrics.GaugeMetric

		gaugesUpdatedState map[string]metrics.GaugeMetric
	}{
		{
			name: "should correct create metric when update non existing gauge metric",

			args: args{
				metricName:  defaultMetricName,
				metricValue: parsedValidGaugeMetricValue,
			},
			gaugesInitState: make(map[string]metrics.GaugeMetric),

			gaugesUpdatedState: map[string]metrics.GaugeMetric{
				defaultMetricName: metrics.NewGauge(defaultMetricName, updatedNonExistingGaugeMetricValue),
			},
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
			metricsStorage := MemStorage{
				gauges: tt.gaugesInitState,
			}

			_, err := metricsStorage.UpdateGaugeMetric(metrics.NewGauge(tt.args.metricName, tt.args.metricValue))
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
		countersInitState map[string]metrics.CounterMetric

		countersUpdatedState map[string]metrics.CounterMetric
	}{
		{
			name: "should correct create metric when update non existing counter metric",

			args: args{
				metricName:  defaultMetricName,
				metricValue: parsedValidCounterMetricValue,
			},
			countersInitState: make(map[string]metrics.CounterMetric),

			countersUpdatedState: map[string]metrics.CounterMetric{
				defaultMetricName: metrics.NewCounter(defaultMetricName, updatedNonExistingCounterMetricValue),
			},
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
			metricsStorage := MemStorage{
				counters: tt.countersInitState,
			}

			_, err := metricsStorage.UpdateCounterMetric(metrics.NewCounter(tt.args.metricName, tt.args.metricValue))
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
		gaugesInitState   map[string]metrics.GaugeMetric
		countersInitState map[string]metrics.CounterMetric

		gaugesUpdatedState   map[string]metrics.GaugeMetric
		countersUpdatedState map[string]metrics.CounterMetric
		storageError         error
	}{
		{
			name: "should return error when update metric with invalid type",

			args: args{
				metricType:  invalidMetricType,
				metricName:  defaultMetricName,
				metricValue: validCounterMetricValue,
			},
			gaugesInitState:   make(map[string]metrics.GaugeMetric),
			countersInitState: make(map[string]metrics.CounterMetric),

			gaugesUpdatedState:   make(map[string]metrics.GaugeMetric),
			countersUpdatedState: make(map[string]metrics.CounterMetric),
			storageError:         errors.InvalidMetricType{},
		},
		{
			name: "should return error when update gauge metric with invalid value",

			args: args{
				metricType:  validGaugeMetricType,
				metricName:  defaultMetricName,
				metricValue: invalidGaugeMetricValue,
			},
			gaugesInitState:   make(map[string]metrics.GaugeMetric),
			countersInitState: make(map[string]metrics.CounterMetric),

			gaugesUpdatedState:   make(map[string]metrics.GaugeMetric),
			countersUpdatedState: make(map[string]metrics.CounterMetric),
			storageError:         errors.InvalidMetricValue{},
		},

		{
			name: "should return error when update counter metric with invalid value",
			args: args{
				metricType:  validCounterMetricType,
				metricName:  defaultMetricName,
				metricValue: invalidCounterMetricValue,
			},

			gaugesInitState:   make(map[string]metrics.GaugeMetric),
			countersInitState: make(map[string]metrics.CounterMetric),

			gaugesUpdatedState:   make(map[string]metrics.GaugeMetric),
			countersUpdatedState: make(map[string]metrics.CounterMetric),
			storageError:         errors.InvalidMetricValue{},
		},

		{
			name: "should correct create metric when update non existing gauge metric with valid value",

			args: args{
				metricType:  validGaugeMetricType,
				metricName:  defaultMetricName,
				metricValue: validGaugeMetricValue,
			},
			gaugesInitState:   make(map[string]metrics.GaugeMetric),
			countersInitState: make(map[string]metrics.CounterMetric),

			gaugesUpdatedState: map[string]metrics.GaugeMetric{
				defaultMetricName: metrics.NewGauge(defaultMetricName, updatedNonExistingGaugeMetricValue),
			},
			countersUpdatedState: make(map[string]metrics.CounterMetric),
			storageError:         nil,
		},

		{
			name: "should correct update metric when update existing gauge metric with valid value",

			args: args{
				metricType:  validGaugeMetricType,
				metricName:  defaultMetricName,
				metricValue: validGaugeMetricValue,
			},
			gaugesInitState: map[string]metrics.GaugeMetric{
				defaultMetricName: metrics.NewGauge(defaultMetricName, initGaugeMetricValue),
			},
			countersInitState: make(map[string]metrics.CounterMetric),

			gaugesUpdatedState: map[string]metrics.GaugeMetric{
				defaultMetricName: metrics.NewGauge(defaultMetricName, updatedExistingGaugeMetricValue),
			},
			countersUpdatedState: make(map[string]metrics.CounterMetric),
			storageError:         nil,
		},

		{
			name: "should correct create metric when update non existing counter metric with valid value",

			args: args{
				metricType:  validCounterMetricType,
				metricName:  defaultMetricName,
				metricValue: validCounterMetricValue,
			},
			gaugesInitState:   make(map[string]metrics.GaugeMetric),
			countersInitState: make(map[string]metrics.CounterMetric),

			gaugesUpdatedState: make(map[string]metrics.GaugeMetric),
			countersUpdatedState: map[string]metrics.CounterMetric{
				defaultMetricName: metrics.NewCounter(defaultMetricName, updatedNonExistingCounterMetricValue),
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
			gaugesInitState: make(map[string]metrics.GaugeMetric),
			countersInitState: map[string]metrics.CounterMetric{
				defaultMetricName: metrics.NewCounter(defaultMetricName, initCounterMetricValue),
			},

			gaugesUpdatedState: make(map[string]metrics.GaugeMetric),
			countersUpdatedState: map[string]metrics.CounterMetric{
				defaultMetricName: metrics.NewCounter(defaultMetricName, updatedExistingCounterMetricValue),
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
			storageError:         errors.InvalidMetricType{},
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
			storageError:         errors.InvalidMetricName{},
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
			storageError:         errors.InvalidMetricName{},
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

func TestGetMetricValueByTypeAndName(t *testing.T) {
	type args struct {
		metricType string
		metricName string
	}

	type want struct {
		value string
		err   error
	}

	testCases := []struct {
		name string

		args              args
		gaugesInitState   map[string]metrics.GaugeMetric
		countersInitState map[string]metrics.CounterMetric

		want want
	}{
		{
			name: "should return error when get metric value with invalid type",

			args: args{
				metricType: invalidMetricType,
				metricName: defaultMetricName,
			},
			gaugesInitState:   make(map[string]metrics.GaugeMetric),
			countersInitState: make(map[string]metrics.CounterMetric),

			want: want{
				value: "",
				err:   errors.InvalidMetricType{},
			},
		},
		{
			name: "should return error when get non existing gauge metric value",

			args: args{
				metricType: validGaugeMetricType,
				metricName: defaultMetricName,
			},
			gaugesInitState:   make(map[string]metrics.GaugeMetric),
			countersInitState: make(map[string]metrics.CounterMetric),

			want: want{
				value: "",
				err:   errors.InvalidMetricName{},
			},
		},

		{
			name: "should correct return metric value when get existing gauge metric value",

			args: args{
				metricType: validGaugeMetricType,
				metricName: defaultMetricName,
			},
			gaugesInitState: map[string]metrics.GaugeMetric{
				defaultMetricName: metrics.NewGauge(defaultMetricName, initGaugeMetricValue),
			},
			countersInitState: make(map[string]metrics.CounterMetric),

			want: want{
				value: parsedInitGaugeMetricValue,
				err:   nil,
			},
		},

		{
			name: "should return error when get non existing counter metric value",

			args: args{
				metricType: validCounterMetricType,
				metricName: defaultMetricName,
			},
			gaugesInitState:   make(map[string]metrics.GaugeMetric),
			countersInitState: make(map[string]metrics.CounterMetric),

			want: want{
				value: "",
				err:   errors.InvalidMetricName{},
			},
		},

		{
			name: "should correct return metric value when get existing gauge metric value",

			args: args{
				metricType: validCounterMetricType,
				metricName: defaultMetricName,
			},
			gaugesInitState: make(map[string]metrics.GaugeMetric),
			countersInitState: map[string]metrics.CounterMetric{
				defaultMetricName: metrics.NewCounter(defaultMetricName, initCounterMetricValue),
			},

			want: want{
				value: parsedInitCounterMetricValue,
				err:   nil,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			metricsStorage := MemStorage{
				gauges:   tt.gaugesInitState,
				counters: tt.countersInitState,
			}

			value, err := metricsStorage.GetMetricValueByTypeAndName(tt.args.metricType, tt.args.metricName)
			if tt.want.err != nil {
				assert.IsType(t, tt.want.err, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, value, tt.want.value)
		})
	}
}

func TestGetGaugeMetric(t *testing.T) {
	testCases := []struct {
		name            string
		gaugesInitState map[string]metrics.GaugeMetric
		metricName      string
		wantError       bool
	}{
		{
			name: "should correct return gauge metric when get existing gauge metric",
			gaugesInitState: map[string]metrics.GaugeMetric{
				"first": metrics.NewGauge("first", 64),
			},
			metricName: "first",
			wantError:  false,
		},
		{
			name: "should return error when get non existing gauge metric",
			gaugesInitState: map[string]metrics.GaugeMetric{
				"first": metrics.NewGauge("first", 64),
			},
			metricName: "second",
			wantError:  true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			metricsStorage := MemStorage{
				gauges: tt.gaugesInitState,
			}

			metric, err := metricsStorage.GetGaugeMetric(tt.metricName)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.metricName, metric.Name)
				assert.Equal(t, metrics.Gauge, metric.Type)
			}
		})
	}
}

func TestGetCounterMetric(t *testing.T) {
	testCases := []struct {
		name              string
		countersInitState map[string]metrics.CounterMetric
		metricName        string
		wantError         bool
	}{
		{
			name: "should correct return gauge metric when get existing gauge metric",
			countersInitState: map[string]metrics.CounterMetric{
				"first": metrics.NewCounter("first", 64),
			},
			metricName: "first",
			wantError:  false,
		},
		{
			name: "should return error when get non existing gauge metric",
			countersInitState: map[string]metrics.CounterMetric{
				"first": metrics.NewCounter("first", 64),
			},
			metricName: "second",
			wantError:  true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			metricsStorage := MemStorage{
				counters: tt.countersInitState,
			}

			metric, err := metricsStorage.GetCounterMetric(tt.metricName)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.metricName, metric.Name)
				assert.Equal(t, metrics.Counter, metric.Type)
			}
		})
	}
}
