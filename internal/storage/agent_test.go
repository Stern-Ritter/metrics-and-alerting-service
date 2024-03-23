package storage

import (
	"testing"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	"github.com/stretchr/testify/assert"
)

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
			storage := NewAgentMemCache(metrics.SupportedGaugeMetrics, metrics.SupportedCounterMetrics)
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
