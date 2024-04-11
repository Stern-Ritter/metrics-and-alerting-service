package agent

import (
	"testing"

	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	cache "github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockAgentMemCache struct {
	cache.AgentMemCache
	mock.Mock
}

func (c *MockAgentMemCache) UpdateMonitorMetrics(model *monitors.Monitor) {
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

func TestUpdateMetrics(t *testing.T) {
	t.Run("should update monitor metrics and 'RandomValue' gauge metric once", func(t *testing.T) {
		logger, err := logger.Initialize("info")
		require.NoError(t, err, "Error init logger")
		mockAgentMemCache := MockAgentMemCache{
			AgentMemCache: cache.NewAgentMemCache(make(map[string]metrics.GaugeMetric), make(map[string]metrics.CounterMetric), logger),
		}
		monitor := monitors.Monitor{}
		mockRandom := utils.NewRandom()

		mockAgentMemCache.On("UpdateMonitorMetrics", &monitor).Return(nil)
		mockAgentMemCache.On("UpdateGaugeMetric", mock.Anything).Return(metrics.GaugeMetric{}, nil)
		UpdateMetrics(&mockAgentMemCache, &monitor, &mockRandom, logger)

		assert.True(t, mockAgentMemCache.AssertNumberOfCalls(t, "UpdateMonitorMetrics", 1), "should update monitor metrics once")
		assert.True(t, mockAgentMemCache.AssertNumberOfCalls(t, "UpdateGaugeMetric", 1), "should update 'RandomValue' gauge metric once")
	})
}

func TestSendMetrics(t *testing.T) {
	t.Run("should reset 'PollCount' counter metric once", func(t *testing.T) {
		client := resty.New()
		url := "localhost:8080"
		endpoint := "/test"
		logger, err := logger.Initialize("info")
		require.NoError(t, err, "Error init logger")

		gaugeMetric := metrics.NewGauge("first", 1.1)
		initGauges := map[string]metrics.GaugeMetric{
			"first": gaugeMetric,
		}

		counterMetic := metrics.NewCounter("second", 2)
		initCounters := map[string]metrics.CounterMetric{
			"second": counterMetic,
		}

		metricsBatchCount := 1

		mockAgentMemCache := MockAgentMemCache{
			AgentMemCache: cache.NewAgentMemCache(initGauges, initCounters, logger),
		}

		mockAgentMemCache.On("ResetMetricValue", mock.Anything, mock.Anything).Return(nil)

		httpmock.ActivateNonDefault(client.GetClient())
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", "http://localhost:8080/test",
			httpmock.NewStringResponder(200, "{}"))

		SendMetrics(client, url, endpoint, &mockAgentMemCache, logger)

		assert.True(t, mockAgentMemCache.AssertNumberOfCalls(t, "ResetMetricValue", 1), "should reset 'PollCount' counter metric once")

		httpmock.GetTotalCallCount()
		info := httpmock.GetCallCountInfo()
		callCount := info["POST http://localhost:8080/test"]
		assert.Equal(t, metricsBatchCount, callCount)
	})
}
