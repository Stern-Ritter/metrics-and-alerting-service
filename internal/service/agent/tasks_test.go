package agent

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/agent"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	cache "github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
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
		client := resty.New()
		aLogger, err := logger.Initialize("info")
		require.NoError(t, err, "Error init logger")
		mockAgentMemCache := MockAgentMemCache{
			AgentMemCache: cache.NewAgentMemCache(make(map[string]metrics.GaugeMetric),
				make(map[string]metrics.CounterMetric), aLogger),
		}
		monitor := monitors.Monitor{}
		mockRandom := utils.NewRandom()
		cfg := config.AgentConfig{}
		agent := NewAgent(client, &mockAgentMemCache, &monitor, &mockRandom, &cfg, aLogger)

		mockAgentMemCache.On("UpdateMonitorMetrics", &monitor).Return(nil)
		mockAgentMemCache.On("UpdateGaugeMetric", mock.Anything).Return(metrics.GaugeMetric{}, nil)

		agent.UpdateMetrics()

		assert.True(t, mockAgentMemCache.AssertNumberOfCalls(t, "UpdateMonitorMetrics", 1),
			"should update monitor metrics once")
		assert.True(t, mockAgentMemCache.AssertNumberOfCalls(t, "UpdateGaugeMetric", 1),
			"should update 'RandomValue' gauge metric once")
	})
}

func TestSendMetrics(t *testing.T) {
	t.Run("should reset 'PollCount' counter metric once", func(t *testing.T) {
		client := resty.New()
		aLogger, err := logger.Initialize("info")
		require.NoError(t, err, "Error init logger")
		mockAgentMemCache := MockAgentMemCache{
			AgentMemCache: cache.NewAgentMemCache(
				map[string]metrics.GaugeMetric{
					"first": metrics.NewGauge("first", 1.1),
				},
				map[string]metrics.CounterMetric{
					"second": metrics.NewCounter("second", 2),
				},
				aLogger),
		}
		monitor := monitors.Monitor{}
		mockRandom := utils.NewRandom()
		cfg := config.AgentConfig{
			SendMetricsURL:      "localhost:8080",
			SendMetricsEndPoint: "/test",
		}
		agent := NewAgent(client, &mockAgentMemCache, &monitor, &mockRandom, &cfg, aLogger)

		mockAgentMemCache.On("ResetMetricValue", mock.Anything, mock.Anything).Return(nil)

		httpmock.ActivateNonDefault(client.GetClient())
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", "http://localhost:8080/test",
			httpmock.NewStringResponder(200, "{}"))

		agent.SendMetrics()

		assert.True(t, mockAgentMemCache.AssertNumberOfCalls(t, "ResetMetricValue", 1),
			"should reset 'PollCount' counter metric once")

		metricsBatchCount := 1

		httpmock.GetTotalCallCount()
		info := httpmock.GetCallCountInfo()
		callCount := info["POST http://localhost:8080/test"]
		assert.Equal(t, metricsBatchCount, callCount)
	})
}
