package agent

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gopkg.in/h2non/gentleman.v2"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/agent"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	cache "github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	mocks "github.com/Stern-Ritter/metrics-and-alerting-service/mocks"
	pb "github.com/Stern-Ritter/metrics-and-alerting-service/proto/gen/metrics/metricsapi/v1"
)

type MockAgentMemCache struct {
	cache.AgentMemCache
	mock.Mock
}

func (c *MockAgentMemCache) UpdateRuntimeMonitorMetrics(model *monitors.RuntimeMonitor) {
	c.Called(model)
}

func (c *MockAgentMemCache) UpdateGaugeMetric(metric metrics.GaugeMetric) (metrics.GaugeMetric, error) {
	args := c.Called(metric)
	return args.Get(0).(metrics.GaugeMetric), args.Error(1)
}

func (c *MockAgentMemCache) ResetMetricValue(metricType, metricName string) error {
	args := c.Called(metricType, metricName)
	return args.Error(0)
}

func TestUpdateRuntimeMetrics(t *testing.T) {
	t.Run("should update monitor metrics and 'RandomValue' gauge metric once", func(t *testing.T) {
		aLogger, err := logger.Initialize("info")
		require.NoError(t, err, "Error init logger")
		mockAgentMemCache := MockAgentMemCache{
			AgentMemCache: cache.NewAgentMemCache(make(map[string]metrics.GaugeMetric),
				make(map[string]metrics.CounterMetric), aLogger),
		}
		runtimeMonitor := monitors.RuntimeMonitor{}
		utilMonitor := monitors.UtilMonitor{}
		mockRandom := utils.NewRandom()
		cfg := config.AgentConfig{}
		agent := NewAgent(&mockAgentMemCache, &runtimeMonitor, &utilMonitor, &mockRandom, &cfg, nil, aLogger)

		mockAgentMemCache.On("UpdateRuntimeMonitorMetrics", &runtimeMonitor).Return(nil)
		mockAgentMemCache.On("UpdateGaugeMetric", mock.Anything).Return(metrics.GaugeMetric{}, nil)

		agent.UpdateRuntimeMetrics()

		assert.True(t, mockAgentMemCache.AssertNumberOfCalls(t, "UpdateRuntimeMonitorMetrics", 1),
			"should update monitor metrics once")
		assert.True(t, mockAgentMemCache.AssertNumberOfCalls(t, "UpdateGaugeMetric", 1),
			"should update 'RandomValue' gauge metric once")
	})
}

func TestStartSendMetricsWorkerPool(t *testing.T) {
	workersCount := 3
	workersCounter := atomic.Int64{}

	cfg := config.AgentConfig{
		RateLimit: workersCount,
	}
	core, recorded := observer.New(zapcore.DebugLevel)
	observerLogger := zap.New(core)
	aLogger := &logger.AgentLogger{Logger: observerLogger}
	metricsCh := make(chan []metrics.Metrics, 10)

	agent := &Agent{
		Config:    &cfg,
		Logger:    aLogger,
		metricsCh: metricsCh,
	}

	wg := &sync.WaitGroup{}
	worker := func(id int, metricsCh <-chan []metrics.Metrics, wg *sync.WaitGroup) {
		workersCounter.Add(1)
		wg.Done()
	}

	agent.StartSendMetricsWorkerPool(wg, worker)

	wg.Wait()

	logs := recorded.FilterMessage("Worker pool started").All()
	assert.Len(t, logs, 1, "should be logged info about the starting one worker pool")
	assert.Equal(t, "starting send metrics worker pool", logs[0].ContextMap()["event"],
		"should be logged info about the starting one worker pool")
	assert.Equal(t, int64(workersCount), workersCounter.Load(),
		"counter should be equal: %d, after %d workers have increased it by 1", workersCount, workersCounter.Load())
}

func TestSendMetricsWithHTTPWorker_OkResponse(t *testing.T) {
	cfg := config.AgentConfig{
		SendMetricsEndPoint: "/metrics",
	}

	core, recorded := observer.New(zapcore.DebugLevel)
	observerLogger := zap.New(core)
	aLogger := &logger.AgentLogger{Logger: observerLogger}

	metricsCh := make(chan []metrics.Metrics, 10)

	metricValue := 22.22
	sentMetrics := []metrics.Metrics{
		{
			ID:    "Alloc",
			MType: "gauge",
			Value: &metricValue,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("content-type")
		assert.Equal(t, "application/json", contentType, "Content-Type header should be application/json")

		ip := r.Header.Get(ipKey)
		assert.NotEmpty(t, ip, "IP address header should not be empty")

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err, "unexpected error reading body")

		gotMetrics := make([]metrics.Metrics, 0)
		err = json.Unmarshal(body, &gotMetrics)
		assert.NoError(t, err, "unexpected error unmarshalling body")
		assert.Equal(t, sentMetrics, gotMetrics, "expected to receive %v metrics, but got %v",
			sentMetrics, gotMetrics)

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := gentleman.New()
	client.URL(ts.URL)

	agent := &Agent{
		Config:                         &cfg,
		Logger:                         aLogger,
		metricsCh:                      metricsCh,
		HTTPClient:                     client,
		sendMetricsBatchRetryIntervals: backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(1 * time.Millisecond)),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	worker := func() {
		agent.SendMetricsWithHTTPWorker(1, metricsCh, &wg)
	}

	go worker()

	metricsCh <- sentMetrics
	close(metricsCh)

	wg.Wait()

	logs := recorded.FilterMessage("Worker started").All()
	assert.Len(t, logs, 1, "should log info about the starting worker")
	assert.Equal(t, "starting send metrics worker", logs[0].ContextMap()["event"],
		"should log info about the starting worker")

	logs = recorded.FilterMessage("Success sent metrics update").All()
	assert.Len(t, logs, 1, "should log info about the successful metrics update")
	assert.Equal(t, "sending metrics update", logs[0].ContextMap()["event"],
		"should log info about the successful metrics update")

	logs = recorded.FilterMessage("Worker stopped").All()
	assert.Len(t, logs, 1, "should log info about the stopping worker")
	assert.Equal(t, "stopping send metrics worker", logs[0].ContextMap()["event"],
		"should log info about the stopping worker")
}

func TestSendMetricsWithHTTPWorker_ErrorResponse(t *testing.T) {
	cfg := config.AgentConfig{
		SendMetricsEndPoint: "/metrics",
	}

	core, recorded := observer.New(zapcore.DebugLevel)
	observerLogger := zap.New(core)
	aLogger := &logger.AgentLogger{Logger: observerLogger}

	metricsCh := make(chan []metrics.Metrics, 10)

	metricValue := 22.22
	sentMetrics := []metrics.Metrics{
		{
			ID:    "Alloc",
			MType: "gauge",
			Value: &metricValue,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		assert.Equal(t, "application/json", contentType, "Content-Type header should be application/json")

		ip := r.Header.Get(ipKey)
		assert.NotEmpty(t, ip, "IP address header should not be empty")

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err, "unexpected error reading body")

		gotMetrics := make([]metrics.Metrics, 0)
		err = json.Unmarshal(body, &gotMetrics)
		assert.NoError(t, err, "unexpected error unmarshalling body")
		assert.Equal(t, sentMetrics, gotMetrics, "expected to receive %v metrics, but got %v",
			sentMetrics, gotMetrics)

		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := gentleman.New()
	client.URL(ts.URL)

	agent := &Agent{
		Config:                         &cfg,
		Logger:                         aLogger,
		metricsCh:                      metricsCh,
		HTTPClient:                     client,
		sendMetricsBatchRetryIntervals: backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(1 * time.Millisecond)),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	worker := func() {
		agent.SendMetricsWithHTTPWorker(1, metricsCh, &wg)
	}

	go worker()

	metricsCh <- sentMetrics
	close(metricsCh)

	wg.Wait()

	logs := recorded.FilterMessage("Worker started").All()
	assert.Len(t, logs, 1, "should log info about the starting worker")
	assert.Equal(t, "starting send metrics worker", logs[0].ContextMap()["event"],
		"should log info about the starting worker")

	logs = recorded.FilterMessage("unsuccess request sent on url: /metrics, status code: 500").All()
	assert.Len(t, logs, 1, "should log error about sending metrics update")
	assert.Equal(t, "sending metrics update", logs[0].ContextMap()["event"],
		"should log info about the sending metrics update")

	logs = recorded.FilterMessage("Worker stopped").All()
	assert.Len(t, logs, 1, "should log info about the stopping worker")
	assert.Equal(t, "stopping send metrics worker", logs[0].ContextMap()["event"],
		"should log info about the stopping worker")
}

func TestSendMetricsWithGrpcWorker_OkResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	core, recorded := observer.New(zapcore.DebugLevel)
	observerLogger := zap.New(core)
	aLogger := &logger.AgentLogger{Logger: observerLogger}

	mockGRPCClient := mocks.NewMockMetricsV1ServiceClient(ctrl)

	metricValue := 22.22
	metricsData := &pb.MetricData{
		Name:        "Alloc",
		Type:        "gauge",
		MetricValue: &pb.MetricData_Value{Value: metricValue},
	}
	updateMetricsBatchRequest := &pb.MetricsV1ServiceUpdateMetricsBatchRequest{
		Metrics: []*pb.MetricData{metricsData},
	}

	mockGRPCClient.EXPECT().
		UpdateMetricsBatch(gomock.Any(), updateMetricsBatchRequest).
		Return(&emptypb.Empty{}, nil).
		Times(1)

	cfg := &config.AgentConfig{}
	agent := &Agent{
		Config:                         cfg,
		Logger:                         aLogger,
		GRPCClient:                     mockGRPCClient,
		sendMetricsBatchRetryIntervals: backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(1 * time.Millisecond)),
	}

	metricsCh := make(chan []metrics.Metrics, 10)

	var wg sync.WaitGroup
	wg.Add(1)
	worker := func() {
		agent.SendMetricsWithGrpcWorker(1, metricsCh, &wg)
	}

	go worker()

	sentMetrics := []metrics.Metrics{
		{
			ID:    "Alloc",
			MType: "gauge",
			Value: &metricValue,
		},
	}
	metricsCh <- sentMetrics
	close(metricsCh)

	wg.Wait()

	logs := recorded.FilterMessage("Worker started").All()
	assert.Len(t, logs, 1, "should log info about the starting worker")
	assert.Equal(t, "starting send metrics worker", logs[0].ContextMap()["event"],
		"should log info about the starting worker")

	logs = recorded.FilterMessage("Success sent metrics update").All()
	assert.Len(t, logs, 1, "should log info about the successful metrics update")
	assert.Equal(t, "sending metrics update", logs[0].ContextMap()["event"],
		"should log info about the successful metrics update")

	logs = recorded.FilterMessage("Worker stopped").All()
	assert.Len(t, logs, 1, "should log info about the stopping worker")
	assert.Equal(t, "stopping send metrics worker", logs[0].ContextMap()["event"],
		"should log info about the stopping worker")
}

func TestSendMetricsWithGrpcWorker_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	core, recorded := observer.New(zapcore.DebugLevel)
	observerLogger := zap.New(core)
	aLogger := &logger.AgentLogger{Logger: observerLogger}

	mockGRPCClient := mocks.NewMockMetricsV1ServiceClient(ctrl)

	metricValue := 22.22
	metricsData := &pb.MetricData{
		Name:        "Alloc",
		Type:        "gauge",
		MetricValue: &pb.MetricData_Value{Value: metricValue},
	}
	updateMetricsBatchRequest := &pb.MetricsV1ServiceUpdateMetricsBatchRequest{
		Metrics: []*pb.MetricData{metricsData},
	}

	mockGRPCClient.EXPECT().
		UpdateMetricsBatch(gomock.Any(), updateMetricsBatchRequest).
		Return(nil, status.Error(codes.Unavailable, "service unavailable")).
		Times(1)

	cfg := &config.AgentConfig{}
	agent := &Agent{
		Config:                         cfg,
		Logger:                         aLogger,
		GRPCClient:                     mockGRPCClient,
		sendMetricsBatchRetryIntervals: backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(1 * time.Millisecond)),
	}

	metricsCh := make(chan []metrics.Metrics, 10)

	var wg sync.WaitGroup
	wg.Add(1)
	worker := func() {
		agent.SendMetricsWithGrpcWorker(1, metricsCh, &wg)
	}

	go worker()

	sentMetrics := []metrics.Metrics{
		{
			ID:    "Alloc",
			MType: "gauge",
			Value: &metricValue,
		},
	}
	metricsCh <- sentMetrics
	close(metricsCh)

	wg.Wait()

	logs := recorded.FilterMessage("Worker started").All()
	assert.Len(t, logs, 1, "should log info about the starting worker")
	assert.Equal(t, "starting send metrics worker", logs[0].ContextMap()["event"],
		"should log info about the starting worker")

	logs = recorded.FilterMessage("unsuccess request sent on url: , status code: 14").All()
	assert.Len(t, logs, 1, "should log info about the error in sending metrics update")
	assert.Equal(t, "sending metrics update", logs[0].ContextMap()["event"],
		"should log info about the sending metrics update with error")

	logs = recorded.FilterMessage("Worker stopped").All()
	assert.Len(t, logs, 1, "should log info about the stopping worker")
	assert.Equal(t, "stopping send metrics worker", logs[0].ContextMap()["event"],
		"should log info about the stopping worker")
}
