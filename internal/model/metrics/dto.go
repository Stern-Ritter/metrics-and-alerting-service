package metrics

var (
	ZeroCounterMetricValue int64   = 0
	ZeroGaugeMetricValue   float64 = 0.0
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewMetrics(id string, typeName string, delta *int64, value *float64) Metrics {
	return Metrics{id, typeName, delta, value}
}

func MetricsToGaugeMetric(m Metrics) GaugeMetric {
	return NewGauge(m.ID, *m.Value)
}

func MetricsToCounterMetric(m Metrics) CounterMetric {
	return NewCounter(m.ID, *m.Delta)
}

func GaugeMetricToMetrics(m GaugeMetric) Metrics {
	name := m.Name
	typeName := string(m.Type)
	value := m.GetValue()
	return NewMetrics(name, typeName, &ZeroCounterMetricValue, &value)
}

func CounterMetricToMetrics(m CounterMetric) Metrics {
	name := m.Name
	typeName := string(m.Type)
	delta := m.GetValue()
	return NewMetrics(name, typeName, &delta, &ZeroGaugeMetricValue)
}
