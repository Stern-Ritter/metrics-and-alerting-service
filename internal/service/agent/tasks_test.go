package agent

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	gmock "gopkg.in/h2non/gentleman-mock.v2"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/context"

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
		client := gentleman.New()
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
		var callCount int64

		client := gentleman.New()
		client.URL("localhost:8080")
		client.Use(gmock.Plugin)
		client.UseHandler("after dial", func(ctx *context.Context, handler context.Handler) {
			atomic.AddInt64(&callCount, 1)
			handler.Next(ctx)
		})

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

		gmock.New("http://localhost:8080").
			Post("/test").
			Reply(200)
		defer gmock.Disable()

		agent.SendMetrics()

		assert.True(t, mockAgentMemCache.AssertNumberOfCalls(t, "ResetMetricValue", 1),
			"should reset 'PollCount' counter metric once")

		var metricsBatchCount int64 = 1
		assert.Equal(t, metricsBatchCount, callCount)
	})
}
