package agent

import (
	"testing"

	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/agent"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/monitors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type MockAgentMemCache struct {
	storage.AgentMemCache
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
			AgentMemCache: storage.NewAgentMemCache(make(map[string]metrics.GaugeMetric), make(map[string]metrics.CounterMetric), logger),
		}
		monitor := monitors.Monitor{}
		mockRandom := utils.NewRandom()

		mockAgentMemCache.On("UpdateMonitorMetrics", &monitor).Return(nil)
		mockAgentMemCache.On("UpdateGaugeMetric", mock.Anything).Return(metrics.GaugeMetric{}, nil)
		UpdateMetrics(&mockAgentMemCache, &monitor, &mockRandom, zap.NewNop())

		assert.True(t, mockAgentMemCache.AssertNumberOfCalls(t, "UpdateMonitorMetrics", 1), "should update monitor metrics once")
		assert.True(t, mockAgentMemCache.AssertNumberOfCalls(t, "UpdateGaugeMetric", 1), "should update 'RandomValue' gauge metric once")
	})
}

func TestSendMetrics(t *testing.T) {
	t.Run("should reset 'PollCount' counter metric once", func(t *testing.T) {
		client := resty.New()
		url := ":8080"
		endpoint := "/test"
		logger, err := logger.Initialize("info")
		require.NoError(t, err, "Error init logger")
		mockAgentMemCache := MockAgentMemCache{
			AgentMemCache: storage.NewAgentMemCache(make(map[string]metrics.GaugeMetric), make(map[string]metrics.CounterMetric), logger),
		}

		mockAgentMemCache.On("ResetMetricValue", mock.Anything, mock.Anything).Return(nil)
		SendMetrics(client, url, endpoint, &mockAgentMemCache, zap.NewNop())

		assert.True(t, mockAgentMemCache.AssertNumberOfCalls(t, "ResetMetricValue", 1), "should reset 'PollCount' counter metric once")
	})
}
