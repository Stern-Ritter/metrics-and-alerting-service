package server

import (
	"context"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
)

type Storage interface {
	UpdateMetric(ctx context.Context, metric metrics.Metrics) error
	UpdateMetrics(ctx context.Context, metrics []metrics.Metrics) error
	GetMetric(ctx context.Context, metric metrics.Metrics) (metrics.Metrics, error)
	GetMetrics(ctx context.Context) (map[string]metrics.GaugeMetric, map[string]metrics.CounterMetric, error)

	Restore(fname string) error
	Save(fname string) error

	Ping(ctx context.Context) error
}
