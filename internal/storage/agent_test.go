package storage

import (
	"testing"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMonitorMetrics(t *testing.T) {
	type state struct {
		gauges   map[string]model.GaugeMetric
		counters map[string]model.CounterMetric
	}

	testCases := []struct {
		name          string
		monitor       *model.Monitor
		expectedState state
	}{{
		name: "should correct update monitor metrics and related counters",
		monitor: &model.Monitor{
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
			gauges: map[string]model.GaugeMetric{
				"Alloc":         model.NewGauge("Alloc", 1.0),
				"BuckHashSys":   model.NewGauge("BuckHashSys", 2.0),
				"Frees":         model.NewGauge("Frees", 3.0),
				"GCCPUFraction": model.NewGauge("GCCPUFraction", 4.0),
				"GCSys":         model.NewGauge("GCSys", 5.0),
				"HeapAlloc":     model.NewGauge("HeapAlloc", 6.0),
				"HeapIdle":      model.NewGauge("HeapIdle", 7.0),
				"HeapInuse":     model.NewGauge("HeapInuse", 8.0),
				"HeapObjects":   model.NewGauge("HeapObjects", 9.0),
				"HeapReleased":  model.NewGauge("HeapReleased", 10.0),
				"HeapSys":       model.NewGauge("HeapSys", 11.0),
				"LastGC":        model.NewGauge("LastGC", 12.0),
				"Lookups":       model.NewGauge("Lookups", 13.0),
				"MCacheInuse":   model.NewGauge("MCacheInuse", 14.0),
				"MCacheSys":     model.NewGauge("MCacheSys", 15.0),
				"MSpanInuse":    model.NewGauge("MSpanInuse", 16.0),
				"MSpanSys":      model.NewGauge("MSpanSys", 17.0),
				"Mallocs":       model.NewGauge("Mallocs", 18.0),
				"NextGC":        model.NewGauge("NextGC", 19.0),
				"NumForcedGC":   model.NewGauge("NumForcedGC", 20.0),
				"NumGC":         model.NewGauge("NumGC", 21.0),
				"OtherSys":      model.NewGauge("OtherSys", 22.0),
				"PauseTotalNs":  model.NewGauge("PauseTotalNs", 23.0),
				"StackInuse":    model.NewGauge("StackInuse", 24.0),
				"StackSys":      model.NewGauge("StackSys", 25.0),
				"Sys":           model.NewGauge("Sys", 26.0),
				"TotalAlloc":    model.NewGauge("TotalAlloc", 27.0),
			},
			counters: map[string]model.CounterMetric{
				"PollCount": model.NewCounter("PollCount", 27),
			},
		},
	},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewAgentMemCache(model.SupportedGaugeMetrics, model.SupportedCounterMetrics)
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
