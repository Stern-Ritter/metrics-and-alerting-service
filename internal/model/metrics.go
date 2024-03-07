package model

type MetricType string

const (
	Gauge   = MetricType("gauge")
	Counter = MetricType("counter")
)

type Metric struct {
	Name string
	Type MetricType
}

type GaugeMetric struct {
	Metric
	value float64
}

func (g *GaugeMetric) SetValue(value float64) {
	g.value = value
}

func (g *GaugeMetric) GetValue() float64 {
	return g.value
}

func (g *GaugeMetric) ClearValue() {
	g.value = 0
}

func NewGauge(name string, value float64) GaugeMetric {
	return GaugeMetric{Metric: Metric{Name: name, Type: Gauge}, value: value}
}

type CounterMetric struct {
	Metric
	value int64
}

func (c *CounterMetric) SetValue(value int64) {
	c.value += value
}

func (c *CounterMetric) GetValue() int64 {
	return c.value
}

func (c *CounterMetric) ClearValue() {
	c.value = 0
}

func NewCounter(name string, value int64) CounterMetric {
	return CounterMetric{Metric: Metric{Name: name, Type: Counter}, value: value}
}

var SupportedGaugeMetrics = map[string]GaugeMetric{
	"Alloc":         NewGauge("Alloc", 0),
	"BuckHashSys":   NewGauge("BuckHashSys", 0),
	"Frees":         NewGauge("Frees", 0),
	"GCCPUFraction": NewGauge("GCCPUFraction", 0),
	"GCSys":         NewGauge("GCSys", 0),
	"HeapAlloc":     NewGauge("HeapAlloc", 0),
	"HeapIdle":      NewGauge("HeapIdle", 0),
	"HeapInuse":     NewGauge("HeapInuse", 0),
	"HeapObjects":   NewGauge("HeapObjects", 0),
	"HeapReleased":  NewGauge("HeapReleased", 0),
	"HeapSys":       NewGauge("HeapSys", 0),
	"LastGC":        NewGauge("LastGC", 0),
	"Lookups":       NewGauge("Lookups", 0),
	"MCacheInuse":   NewGauge("MCacheInuse", 0),
	"MCacheSys":     NewGauge("MCacheSys", 0),
	"MSpanInuse":    NewGauge("MSpanInuse", 0),
	"MSpanSys":      NewGauge("MSpanSys", 0),
	"Mallocs":       NewGauge("Mallocs", 0),
	"NextGC":        NewGauge("NextGC", 0),
	"NumForcedGC":   NewGauge("NumForcedGC", 0),
	"NumGC":         NewGauge("NumGC", 0),
	"OtherSys":      NewGauge("OtherSys", 0),
	"PauseTotalNs":  NewGauge("PauseTotalNs", 0),
	"StackInuse":    NewGauge("StackInuse", 0),
	"StackSys":      NewGauge("StackSys", 0),
	"Sys":           NewGauge("Sys", 0),
	"TotalAlloc":    NewGauge("TotalAlloc", 0),
	"RandomValue":   NewGauge("RandomValue", 0),
}

var SupportedCounterMetrics = map[string]CounterMetric{
	"PollCount": NewCounter("PollCount", 0),
}
