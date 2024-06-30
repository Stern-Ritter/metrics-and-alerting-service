package metrics

// MetricType is the type of a metric.
type MetricType string

const (
	// Gauge is the gauge metric type.
	Gauge = MetricType("gauge")
	// Counter is the counter metric type.
	Counter = MetricType("counter")
)

// Metric contains common attributes of a metric.
type Metric struct {
	Name string     `json:"name"` // The name of the metric
	Type MetricType `json:"type"` // The type of the metric (gauge or counter)
}

// GaugeMetric is a metric with a float value.
type GaugeMetric struct {
	Metric `json:"metric"` // The common attributes of a metric
	Value  float64         `json:"value"` // The value of the gauge metric
}

// SetValue sets the value of the gauge metric.
func (g *GaugeMetric) SetValue(value float64) {
	g.Value = value
}

// GetValue returns the value of the gauge metric.
func (g *GaugeMetric) GetValue() float64 {
	return g.Value
}

// ClearValue resets the value of the gauge metric to 0.
func (g *GaugeMetric) ClearValue() {
	g.Value = 0
}

// NewGauge is constructor for creating a new GaugeMetric with the specified name and value.
func NewGauge(name string, value float64) GaugeMetric {
	return GaugeMetric{Metric: Metric{Name: name, Type: Gauge}, Value: value}
}

// CounterMetric is a metric with an int value.
type CounterMetric struct {
	Metric `json:"metric"`
	Value  int64 `json:"value"`
}

// SetValue increments the value of the counter metric by the specified value.
func (c *CounterMetric) SetValue(value int64) {
	c.Value += value
}

// GetValue returns the value of the counter metric.
func (c *CounterMetric) GetValue() int64 {
	return c.Value
}

// ClearValue resets the value of the counter metric to 0.
func (c *CounterMetric) ClearValue() {
	c.Value = 0
}

// NewCounter is constructor for creating a new CounterMetric with the specified name and initial value.
func NewCounter(name string, value int64) CounterMetric {
	return CounterMetric{Metric: Metric{Name: name, Type: Counter}, Value: value}
}

// SupportedGaugeMetrics is a predefined map of supported gauge metrics.
var SupportedGaugeMetrics = map[string]GaugeMetric{
	"Alloc":           NewGauge("Alloc", 0),
	"BuckHashSys":     NewGauge("BuckHashSys", 0),
	"Frees":           NewGauge("Frees", 0),
	"GCCPUFraction":   NewGauge("GCCPUFraction", 0),
	"GCSys":           NewGauge("GCSys", 0),
	"HeapAlloc":       NewGauge("HeapAlloc", 0),
	"HeapIdle":        NewGauge("HeapIdle", 0),
	"HeapInuse":       NewGauge("HeapInuse", 0),
	"HeapObjects":     NewGauge("HeapObjects", 0),
	"HeapReleased":    NewGauge("HeapReleased", 0),
	"HeapSys":         NewGauge("HeapSys", 0),
	"LastGC":          NewGauge("LastGC", 0),
	"Lookups":         NewGauge("Lookups", 0),
	"MCacheInuse":     NewGauge("MCacheInuse", 0),
	"MCacheSys":       NewGauge("MCacheSys", 0),
	"MSpanInuse":      NewGauge("MSpanInuse", 0),
	"MSpanSys":        NewGauge("MSpanSys", 0),
	"Mallocs":         NewGauge("Mallocs", 0),
	"NextGC":          NewGauge("NextGC", 0),
	"NumForcedGC":     NewGauge("NumForcedGC", 0),
	"NumGC":           NewGauge("NumGC", 0),
	"OtherSys":        NewGauge("OtherSys", 0),
	"PauseTotalNs":    NewGauge("PauseTotalNs", 0),
	"StackInuse":      NewGauge("StackInuse", 0),
	"StackSys":        NewGauge("StackSys", 0),
	"Sys":             NewGauge("Sys", 0),
	"TotalAlloc":      NewGauge("TotalAlloc", 0),
	"RandomValue":     NewGauge("RandomValue", 0),
	"TotalMemory":     NewGauge("TotalMemory", 0),
	"FreeMemory":      NewGauge("FreeMemory", 0),
	"CPUutilization1": NewGauge("CPUutilization1", 0),
}

// SupportedCounterMetrics is a predefined map of supported counter metrics.
var SupportedCounterMetrics = map[string]CounterMetric{
	"PollCount": NewCounter("PollCount", 0),
}
