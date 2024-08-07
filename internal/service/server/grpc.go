package server

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	er "github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	pb "github.com/Stern-Ritter/metrics-and-alerting-service/proto/gen/metrics/metricsapi/v1"
)

// UpdateMetric updates a single metric based on request.
func (s *Server) UpdateMetric(ctx context.Context, in *pb.MetricsV1ServiceUpdateMetricRequest) (*pb.MetricsV1ServiceUpdateMetricResponse, error) {
	err := in.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	metric := metrics.MetricDataToMetrics(in.Metric)
	updatedMetric, err := s.MetricService.UpdateMetricWithBody(ctx, metric, s.isSyncSaveStorageState(),
		s.Config.FileStoragePath)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	metricData := metrics.MetricsToMetricData(updatedMetric)
	resp := pb.MetricsV1ServiceUpdateMetricResponse{
		Metric: metricData,
	}

	return &resp, nil
}

// UpdateMetricsBatch updates multiple metrics based on request.
func (s *Server) UpdateMetricsBatch(ctx context.Context, in *pb.MetricsV1ServiceUpdateMetricsBatchRequest) (*pb.MetricsV1ServiceUpdateMetricsBatchResponse, error) {
	err := in.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	metrics := metrics.RepeatedMetricDataToMetrics(in.Metrics)

	err = s.MetricService.UpdateMetricsBatchWithBody(ctx, metrics, s.isSyncSaveStorageState(),
		s.Config.FileStoragePath)

	if err != nil {
		var invalidMetricType er.InvalidMetricType
		var invalidMetricValue er.InvalidMetricValue
		if errors.As(err, &invalidMetricType) || errors.As(err, &invalidMetricValue) {
			return &pb.MetricsV1ServiceUpdateMetricsBatchResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}

		return &pb.MetricsV1ServiceUpdateMetricsBatchResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &pb.MetricsV1ServiceUpdateMetricsBatchResponse{}, nil
}

// GetMetric retrieves the data of a single metric based on request.
func (s *Server) GetMetric(ctx context.Context, in *pb.MetricsV1ServiceGetMetricRequest) (*pb.MetricsV1ServiceGetMetricResponse, error) {
	err := in.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	metric := metrics.MetricInfoToMetrics(in.Metric)

	savedMetric, err := s.MetricService.GetMetricValueWithBody(ctx, metric)
	if err != nil {
		var invalidMetricType er.InvalidMetricType
		var invalidMetricName er.InvalidMetricName
		if errors.As(err, &invalidMetricType) || errors.As(err, &invalidMetricName) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	metricData := metrics.MetricsToMetricData(savedMetric)
	resp := pb.MetricsV1ServiceGetMetricResponse{
		Metric: metricData,
	}

	return &resp, nil
}

// GetMetrics retrieves all metrics list.
func (s *Server) GetMetrics(ctx context.Context, in *pb.MetricsV1ServiceGetMetricsRequest) (*pb.MetricsV1ServiceGetMetricsResponse, error) {
	gauges, counters, err := s.MetricService.GetMetrics(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	metrics := getMetricsString(gauges, counters)
	resp := pb.MetricsV1ServiceGetMetricsResponse{
		Metrics: metrics,
	}

	return &resp, nil
}

// Ping checks the availability of the database.
func (s *Server) Ping(ctx context.Context, in *pb.MetricsV1ServicePingRequest) (*pb.MetricsV1ServicePingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	err := s.MetricService.PingDatabase(ctx)
	if err != nil {
		return &pb.MetricsV1ServicePingResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &pb.MetricsV1ServicePingResponse{}, nil
}
