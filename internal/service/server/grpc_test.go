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

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	pb "github.com/Stern-Ritter/metrics-and-alerting-service/proto/gen/metrics/metricsapi/v1"
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
		req                    *pb.MetricsV1ServiceUpdateMetricRequest
		resp                   *pb.MetricsV1ServiceUpdateMetricResponse
		useStorageUpdateMetric bool
		storageUpdateMetricErr error
		useStorageGetMetric    bool
		storageGetMetric       metrics.Metrics
		storageGetMetricErr    error
		expectedErr            error
	}{
		{
			name: "should return error when update metric request is invalid",
			req: &pb.MetricsV1ServiceUpdateMetricRequest{
				Metric: &pb.MetricData{
					Name:        gaugeMetricName,
					Type:        "unknown",
					MetricValue: &pb.MetricData_Value{Value: gaugeMetricValue},
				},
			},
			resp:        nil,
			expectedErr: status.Error(codes.InvalidArgument, "invalid MetricsV1ServiceUpdateMetricRequest.Metric: embedded message failed validation | caused by: invalid MetricData.Type: value must be in list [gauge counter]"),
		},
		{
			name: "should success update metric when update metric request is valid",
			req: &pb.MetricsV1ServiceUpdateMetricRequest{
				Metric: &pb.MetricData{
					Name:        gaugeMetricName,
					Type:        gaugeMetricType,
					MetricValue: &pb.MetricData_Value{Value: gaugeMetricValue},
				},
			},
			resp: &pb.MetricsV1ServiceUpdateMetricResponse{
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
		req                     *pb.MetricsV1ServiceUpdateMetricsBatchRequest
		useStorageUpdateMetrics bool
		storageUpdateMetricsErr error
		expectedErr             error
	}{
		{
			name: "should return error when update metrics request is invalid",
			req: &pb.MetricsV1ServiceUpdateMetricsBatchRequest{
				Metrics: []*pb.MetricData{
					{
						Name:        gaugeMetricName,
						Type:        "unknown",
						MetricValue: &pb.MetricData_Value{Value: gaugeInitValue},
					},
				},
			},
			expectedErr: status.Error(codes.InvalidArgument, "invalid MetricsV1ServiceUpdateMetricsBatchRequest.Metrics[0]: embedded message failed validation | caused by: invalid MetricData.Type: value must be in list [gauge counter]"),
		},
		{
			name: "should success update metrics when update metrics request is valid",
			req: &pb.MetricsV1ServiceUpdateMetricsBatchRequest{
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
		req                 *pb.MetricsV1ServiceGetMetricRequest
		resp                *pb.MetricsV1ServiceGetMetricResponse
		useStorageGetMetric bool
		storageGetMetric    metrics.Metrics
		storageGetMetricErr error
		expectedErr         error
	}{
		{
			name: "should return error when get metric request is invalid",
			req: &pb.MetricsV1ServiceGetMetricRequest{
				Metric: &pb.MetricInfo{
					Name: gaugeMetricName,
					Type: "unknown",
				},
			},
			resp:        nil,
			expectedErr: status.Error(codes.InvalidArgument, "invalid MetricsV1ServiceGetMetricRequest.Metric: embedded message failed validation | caused by: invalid MetricInfo.Type: value must be in list [gauge counter]"),
		},
		{
			name: "should return metric data when get metric request is valid",
			req: &pb.MetricsV1ServiceGetMetricRequest{
				Metric: &pb.MetricInfo{
					Name: gaugeMetricName,
					Type: gaugeMetricType,
				},
			},
			resp: &pb.MetricsV1ServiceGetMetricResponse{
				Metric: &pb.MetricData{
					Name:        gaugeMetricName,
					Type:        gaugeMetricType,
					MetricValue: &pb.MetricData_Value{Value: gaugeMetricValue},
				},
			},
			useStorageGetMetric: true,
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
			if tt.useStorageGetMetric {
				mockStorage.
					EXPECT().
					GetMetric(gomock.Any(), gomock.Any()).
					Return(tt.storageGetMetric, tt.storageGetMetricErr)
			}

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
		resp                     *pb.MetricsV1ServiceGetMetricsResponse
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
			resp: &pb.MetricsV1ServiceGetMetricsResponse{
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

			resp, err := s.GetMetrics(context.Background(), &pb.MetricsV1ServiceGetMetricsRequest{})
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

			_, err := s.Ping(context.Background(), &pb.MetricsV1ServicePingRequest{})
			if tt.expectedErr != nil {
				assert.ErrorIs(t, tt.expectedErr, err, "should return error: %s, got: %s", tt.expectedErr, err)
			} else {
				assert.NoError(t, err, "shouldn't return error, but got: %s", err)
			}
		})
	}
}
