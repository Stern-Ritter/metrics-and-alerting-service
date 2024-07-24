package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gentleman.v2"

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
		client := gentleman.New()
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
		agent := NewAgent(client, &mockAgentMemCache, &runtimeMonitor, &utilMonitor, &mockRandom, &cfg, nil,
			aLogger)

		mockAgentMemCache.On("UpdateRuntimeMonitorMetrics", &runtimeMonitor).Return(nil)
		mockAgentMemCache.On("UpdateGaugeMetric", mock.Anything).Return(metrics.GaugeMetric{}, nil)

		agent.UpdateRuntimeMetrics()

		assert.True(t, mockAgentMemCache.AssertNumberOfCalls(t, "UpdateRuntimeMonitorMetrics", 1),
			"should update monitor metrics once")
		assert.True(t, mockAgentMemCache.AssertNumberOfCalls(t, "UpdateGaugeMetric", 1),
			"should update 'RandomValue' gauge metric once")
	})
}
