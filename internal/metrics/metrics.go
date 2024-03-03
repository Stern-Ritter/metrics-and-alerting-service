package metrics

const (
	TypeGauge   = "gauge"
	TypeCounter = "counter"
)

type Gauge struct {
	value float64
}

func (g *Gauge) SetValue(value float64) {
	g.value = value
}

func (g *Gauge) GetValue() float64 {
	return g.value
}

func NewGauge(value float64) Gauge {
	return Gauge{value}
}

type Counter struct {
	value int64
}

func (c *Counter) SetValue(value int64) {
	c.value += value
}

func (c *Counter) GetValue() int64 {
	return c.value
}

func NewCounter(value int64) Counter {
	return Counter{value}
}
