package storage

import (
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model"
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
			gauges:   make(map[string]model.GaugeMetric),
			counters: make(map[string]model.CounterMetric),
		},
	}
}
