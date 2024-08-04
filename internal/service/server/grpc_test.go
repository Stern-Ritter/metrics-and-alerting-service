package server

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	er "github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	pb "github.com/Stern-Ritter/metrics-and-alerting-service/proto/gen/metrics"
)

var (
	gaugeMetricName  = "Alloc"
	gaugeMetricType  = "gauge"
	gaugeMetricValue = 22.22
)

func TestUpdateMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := NewMockStorage(ctrl)
	config := &server.ServerConfig{}
	logger, err := logger.Initialize("info")
	require.NoError(t, err, "Error init logger")
	metricService := NewMetricService(mockStorage, logger)
	s := NewServer(metricService, config, nil, nil, logger)

	testCases := []struct {
		name                   string
		req                    *pb.UpdateMetricRequest
		resp                   *pb.UpdateMetricResponse
		useStorageUpdateMetric bool
		storageUpdateMetricErr error
		useStorageGetMetric    bool
		storageGetMetric       metrics.Metrics
		storageGetMetricErr    error
		expectedErr            error
	}{
		{
			name: "should return error when update metric request is invalid",
			req: &pb.UpdateMetricRequest{
				Metric: &pb.MetricData{
					Name:        gaugeMetricName,
					Type:        "unknown",
					MetricValue: &pb.MetricData_Value{Value: gaugeMetricValue},
				},
			},
			resp:                   nil,
			useStorageUpdateMetric: true,
			storageUpdateMetricErr: er.NewInvalidMetricType("Invalid metric type: unknown", nil),
			expectedErr:            status.Error(codes.InvalidArgument, "Invalid metric type: unknown"),
		},
		{
			name: "should success update metric when update metric request is valid",
			req: &pb.UpdateMetricRequest{
				Metric: &pb.MetricData{
					Name:        gaugeMetricName,
					Type:        gaugeMetricType,
					MetricValue: &pb.MetricData_Value{Value: gaugeMetricValue},
				},
			},
			resp: &pb.UpdateMetricResponse{
				Metric: &pb.MetricData{
					Name:        gaugeMetricName,
					Type:        gaugeMetricType,
					MetricValue: &pb.MetricData_Value{Value: gaugeMetricValue},
				},
			},
			useStorageUpdateMetric: true,
			storageUpdateMetricErr: nil,
			useStorageGetMetric:    true,
			storageGetMetric: metrics.Metrics{
				ID:    gaugeMetricName,
				MType: gaugeMetricType,
				Value: &gaugeMetricValue,
			},
			storageGetMetricErr: nil,
			expectedErr:         nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.useStorageUpdateMetric {
				mockStorage.
					EXPECT().
					UpdateMetric(gomock.Any(), gomock.Any()).
					Return(tt.storageUpdateMetricErr)
			}
			if tt.useStorageGetMetric {
				mockStorage.
					EXPECT().
					GetMetric(gomock.Any(), gomock.Any()).
					Return(tt.storageGetMetric, tt.storageGetMetricErr)
			}

			resp, err := s.UpdateMetric(context.Background(), tt.req)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, tt.expectedErr, err, "should return error: %s, got: %s", tt.expectedErr, err)
			} else {
				assert.NoError(t, err, "shouldn't return error, but got: %s", err)
				assert.Equal(t, tt.resp, resp, "should return response: %v, got: %v", tt.resp, resp)
			}
		})
	}
}

func TestUpdateMetricsBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := NewMockStorage(ctrl)
	config := &server.ServerConfig{}
	logger, err := logger.Initialize("info")
	require.NoError(t, err, "Error init logger")
	metricService := NewMetricService(mockStorage, logger)
	s := NewServer(metricService, config, nil, nil, logger)

	testCases := []struct {
		name                    string
		req                     *pb.UpdateMetricsBatchRequest
		useStorageUpdateMetrics bool
		storageUpdateMetricsErr error
		expectedErr             error
	}{
		{
			name: "should return error when update metrics request is invalid",
			req: &pb.UpdateMetricsBatchRequest{
				Metrics: []*pb.MetricData{
					{
						Name:        gaugeMetricName,
						Type:        "unknown",
						MetricValue: &pb.MetricData_Value{Value: gaugeInitValue},
					},
				},
			},
			useStorageUpdateMetrics: true,
			storageUpdateMetricsErr: er.NewInvalidMetricType("Invalid metric type: unknown", nil),
			expectedErr:             status.Error(codes.InvalidArgument, "Invalid metric type: unknown"),
		},
		{
			name: "should success update metrics when update metrics request is valid",
			req: &pb.UpdateMetricsBatchRequest{
				Metrics: []*pb.MetricData{
					{
						Name:        gaugeMetricName,
						Type:        gaugeMetricType,
						MetricValue: &pb.MetricData_Value{Value: gaugeMetricValue},
					},
				},
			},
			useStorageUpdateMetrics: true,
			storageUpdateMetricsErr: nil,
			expectedErr:             nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.useStorageUpdateMetrics {
				mockStorage.
					EXPECT().
					UpdateMetrics(gomock.Any(), gomock.Any()).
					Return(tt.storageUpdateMetricsErr)
			}

			_, err := s.UpdateMetricsBatch(context.Background(), tt.req)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, tt.expectedErr, err, "should return error: %s, got: %s", tt.expectedErr, err)
			} else {
				assert.NoError(t, err, "shouldn't return error, but got: %s", err)
			}
		})
	}
}

func TestGetMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := NewMockStorage(ctrl)
	config := &server.ServerConfig{}
	logger, err := logger.Initialize("info")
	require.NoError(t, err, "Error init logger")
	metricService := NewMetricService(mockStorage, logger)
	s := NewServer(metricService, config, nil, nil, logger)

	testCases := []struct {
		name                string
		req                 *pb.GetMetricRequest
		resp                *pb.GetMetricResponse
		storageGetMetric    metrics.Metrics
		storageGetMetricErr error
		expectedErr         error
	}{
		{
			name: "should return error when get metric request is invalid",
			req: &pb.GetMetricRequest{
				Metric: &pb.MetricInfo{
					Name: gaugeMetricName,
					Type: "unknown",
				},
			},
			resp:                nil,
			storageGetMetricErr: er.NewInvalidMetricType("Invalid metric type: unknown", nil),
			expectedErr:         status.Error(codes.NotFound, "Invalid metric type: unknown"),
		},
		{
			name: "should return metric data when get metric request is valid",
			req: &pb.GetMetricRequest{
				Metric: &pb.MetricInfo{
					Name: gaugeMetricName,
					Type: gaugeMetricType,
				},
			},
			resp: &pb.GetMetricResponse{
				Metric: &pb.MetricData{
					Name:        gaugeMetricName,
					Type:        gaugeMetricType,
					MetricValue: &pb.MetricData_Value{Value: gaugeMetricValue},
				},
			},
			storageGetMetric: metrics.Metrics{
				ID:    gaugeMetricName,
				MType: gaugeMetricType,
				Value: &gaugeMetricValue,
			},
			storageGetMetricErr: nil,
			expectedErr:         nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage.
				EXPECT().
				GetMetric(gomock.Any(), gomock.Any()).
				Return(tt.storageGetMetric, tt.storageGetMetricErr)

			resp, err := s.GetMetric(context.Background(), tt.req)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, tt.expectedErr, err, "should return error: %s, got: %s", tt.expectedErr, err)
			} else {
				assert.NoError(t, err, "shouldn't return error, but got: %s", err)
				assert.Equal(t, tt.resp, resp, "should return response: %v, got: %v", tt.resp, resp)
			}
		})
	}
}

func TestGetMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := NewMockStorage(ctrl)
	config := &server.ServerConfig{}
	logger, err := logger.Initialize("info")
	require.NoError(t, err, "Error init logger")
	metricService := NewMetricService(mockStorage, logger)
	s := NewServer(metricService, config, nil, nil, logger)

	testCases := []struct {
		name                     string
		resp                     *pb.GetMetricsResponse
		storageGetGaugeMetrics   map[string]metrics.GaugeMetric
		storageGetCounterMetrics map[string]metrics.CounterMetric
		storageGetMetricsErr     error
		expectedErr              error
	}{
		{
			name:                 "should return error when storage returns error",
			resp:                 nil,
			storageGetMetricsErr: errors.New("Internal server error"),
			expectedErr:          status.Error(codes.Internal, "Internal server error"),
		},
		{
			name: "should return metrics list when storage doesn`t return error",
			resp: &pb.GetMetricsResponse{
				Metrics: "Alloc",
			},
			storageGetGaugeMetrics: map[string]metrics.GaugeMetric{
				gaugeMetricName: {
					Metric: metrics.Metric{
						Name: gaugeMetricName,
						Type: metrics.Gauge,
					},
					Value: gaugeMetricValue,
				},
			},
			storageGetCounterMetrics: make(map[string]metrics.CounterMetric),
			storageGetMetricsErr:     nil,
			expectedErr:              nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage.
				EXPECT().
				GetMetrics(gomock.Any()).
				Return(tt.storageGetGaugeMetrics, tt.storageGetCounterMetrics, tt.storageGetMetricsErr)

			resp, err := s.GetMetrics(context.Background(), &emptypb.Empty{})
			if tt.expectedErr != nil {
				assert.ErrorIs(t, tt.expectedErr, err, "should return error: %s, got: %s", tt.expectedErr, err)
			} else {
				assert.NoError(t, err, "shouldn't return error, but got: %s", err)
				assert.Equal(t, tt.resp, resp, "should return response: %v, got: %v", tt.resp, resp)
			}
		})
	}
}

func TestPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := NewMockStorage(ctrl)
	config := &server.ServerConfig{}
	logger, err := logger.Initialize("info")
	require.NoError(t, err, "Error init logger")
	metricService := NewMetricService(mockStorage, logger)
	s := NewServer(metricService, config, nil, nil, logger)

	testCases := []struct {
		name           string
		storagePingErr error
		expectedErr    error
	}{
		{
			name:           "should return error when storage returns error",
			storagePingErr: errors.New("Internal server error"),
			expectedErr:    status.Error(codes.Internal, "Internal server error"),
		},
		{
			name:           "should succeed when storage doesn't return error",
			storagePingErr: nil,
			expectedErr:    nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage.
				EXPECT().
				Ping(gomock.Any()).
				Return(tt.storagePingErr)

			_, err := s.Ping(context.Background(), &emptypb.Empty{})
			if tt.expectedErr != nil {
				assert.ErrorIs(t, tt.expectedErr, err, "should return error: %s, got: %s", tt.expectedErr, err)
			} else {
				assert.NoError(t, err, "shouldn't return error, but got: %s", err)
			}
		})
	}
}
