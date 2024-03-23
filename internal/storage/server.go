package storage

import (
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
)

type ServerStorage interface {
	Storage
}

type ServerMemStorage struct {
	MemStorage
}

func NewServerMemStorage() ServerMemStorage {
	return ServerMemStorage{
		MemStorage: MemStorage{
			gauges:   make(map[string]metrics.GaugeMetric),
			counters: make(map[string]metrics.CounterMetric),
		},
	}
}
