package server

import (
	"context"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
)

// Storage defines an interface for a metrics storage.
type Storage interface {
	// UpdateMetric updates a single metric in the storage.
	UpdateMetric(ctx context.Context, metric metrics.Metrics) error
	// UpdateMetrics updates multiple metrics in the storage.
	UpdateMetrics(ctx context.Context, metrics []metrics.Metrics) error
	// GetMetric gets a single metric from the storage.
	GetMetric(ctx context.Context, metric metrics.Metrics) (metrics.Metrics, error)
	// GetMetrics gets all metrics from the storage.
	GetMetrics(ctx context.Context) (map[string]metrics.GaugeMetric, map[string]metrics.CounterMetric, error)
	// Restore restores the storage state from a file.
	Restore(fName string) error
	// Save saves the storage state to a file.
	Save(fName string) error
	// Ping checks the connection to the storage.
	Ping(ctx context.Context) error
	// Close closes the connection to the storage.
	Close() error
}
