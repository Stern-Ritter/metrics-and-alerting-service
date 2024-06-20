package metrics

import (
	"fmt"
	"strconv"

	er "github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
)

// Metrics is a data transfer object for generic metric that can be gauge or counter.
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// NewMetricsWithStringValue is constructor for creating a new Metrics with the specified name, type, and value as a string.
// It returns an error if type is unsupported.
func NewMetricsWithStringValue(mName string, mTypeName string, value string) (Metrics, error) {
	switch MetricType(mTypeName) {
	case Gauge:
		v, err := parseGaugeMetricValue(value)
		if err != nil {
			return Metrics{}, err
		}
		return Metrics{ID: mName, MType: mTypeName, Value: &v}, nil

	case Counter:
		v, err := parseCounterMetricValue(value)
		if err != nil {
			return Metrics{}, err
		}
		return Metrics{ID: mName, MType: mTypeName, Delta: &v}, nil

	default:
		return Metrics{}, er.NewInvalidMetricType(fmt.Sprintf("Invalid metric type: %s", mTypeName), nil)
	}
}

// NewMetricsWithNumberValue is constructor for creating a new Metrics with the specified name, type, and value as a float64.
// It returns an error if type is unsupported.
func NewMetricsWithNumberValue(mName string, mTypeName string, value float64) (Metrics, error) {
	switch MetricType(mTypeName) {
	case Gauge:
		v := value
		return Metrics{ID: mName, MType: mTypeName, Value: &v}, nil

	case Counter:
		v := int64(value)
		return Metrics{ID: mName, MType: mTypeName, Delta: &v}, nil

	default:
		return Metrics{}, er.NewInvalidMetricType(fmt.Sprintf("Invalid metric type: %s", mTypeName), nil)
	}
}

// GetValue returns the value of the Metrics as a float64.
func (m Metrics) GetValue() (float64, error) {
	switch MetricType(m.MType) {
	case Gauge:
		return *m.Value, nil
	case Counter:
		return float64(*m.Delta), nil
	default:
		return 0, er.NewInvalidMetricType(fmt.Sprintf("Invalid metric type: %s", m.MType), nil)
	}
}

// MetricsToGaugeMetric maps a Metrics to a GaugeMetric.
func MetricsToGaugeMetric(m Metrics) GaugeMetric {
	return NewGauge(m.ID, *m.Value)
}

// MetricsToCounterMetric maps a Metrics to a CounterMetric.
func MetricsToCounterMetric(m Metrics) CounterMetric {
	return NewCounter(m.ID, *m.Delta)
}

// GaugeMetricToMetrics maps a GaugeMetric to a Metrics.
func GaugeMetricToMetrics(m GaugeMetric) Metrics {
	name := m.Name
	typeName := string(m.Type)
	value := m.GetValue()
	return Metrics{ID: name, MType: typeName, Value: &value}
}

// CounterMetricToMetrics maps a CounterMetric to a Metrics.
func CounterMetricToMetrics(m CounterMetric) Metrics {
	name := m.Name
	typeName := string(m.Type)
	delta := m.GetValue()
	return Metrics{ID: name, MType: typeName, Delta: &delta}
}

func parseGaugeMetricValue(v string) (float64, error) {
	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, er.NewInvalidMetricValue(
			fmt.Sprintf("The value for the %s metric should be of float64 type", Gauge), err)
	}

	return value, nil
}

func parseCounterMetricValue(v string) (int64, error) {
	value, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, er.NewInvalidMetricValue(
			fmt.Sprintf("The value for the %s metric should be of int64 type", Counter), err)
	}

	return value, nil
}
