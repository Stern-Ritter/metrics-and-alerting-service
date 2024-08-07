package metrics

import (
	pb "github.com/Stern-Ritter/metrics-and-alerting-service/proto/gen/metrics/metricsapi/v1"
)

// MetricDataToMetrics converts a pb.MetricData to a Metrics structure.
func MetricDataToMetrics(md *pb.MetricData) Metrics {
	var delta *int64
	var value *float64

	switch v := md.MetricValue.(type) {
	case *pb.MetricData_Delta:
		delta = &v.Delta
	case *pb.MetricData_Value:
		value = &v.Value
	}

	m := Metrics{
		ID:    md.Name,
		MType: md.Type,
		Delta: delta,
		Value: value,
	}

	return m
}

// RepeatedMetricDataToMetrics converts a slice of pb.MetricData to a slice of Metrics.
func RepeatedMetricDataToMetrics(md []*pb.MetricData) []Metrics {
	metrics := make([]Metrics, len(md))

	for i, m := range md {
		metrics[i] = MetricDataToMetrics(m)
	}

	return metrics
}

// MetricInfoToMetrics converts a pb.MetricInfo to a Metrics structure.
func MetricInfoToMetrics(mi *pb.MetricInfo) Metrics {
	m := Metrics{
		ID:    mi.Name,
		MType: mi.Type,
	}

	return m
}

// MetricsToMetricData converts a Metrics structure to a pb.MetricData.
func MetricsToMetricData(m Metrics) *pb.MetricData {
	md := &pb.MetricData{
		Name: m.ID,
		Type: m.MType,
	}

	switch MetricType(m.MType) {
	case Gauge:
		md.MetricValue = &pb.MetricData_Value{Value: *m.Value}
	case Counter:
		md.MetricValue = &pb.MetricData_Delta{Delta: *m.Delta}
	}

	return md
}

// MetricsToRepeatedMetricData converts a slice of Metrics to a slice of pb.MetricData.
func MetricsToRepeatedMetricData(metrics []Metrics) []*pb.MetricData {
	md := make([]*pb.MetricData, len(metrics))

	for i, m := range metrics {
		md[i] = MetricsToMetricData(m)
	}

	return md
}
