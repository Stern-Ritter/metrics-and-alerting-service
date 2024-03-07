package transport

import (
	"net/http"
	"testing"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAgentMemCache struct {
	storage.AgentMemCache
	mock.Mock
}

func (c *MockAgentMemCache) UpdateMonitorMetrics(model *model.Monitor) {
	c.Called(model)
}

func (c *MockAgentMemCache) UpdateGaugeMetric(metric model.GaugeMetric) error {
	args := c.Called(metric)
	return args.Error(0)
}

func (c *MockAgentMemCache) ResetMetricValue(metricType, metricName string) error {
	args := c.Called(metricType, metricName)
	return args.Error(0)
}

func TestUpdateMetrics(t *testing.T) {
	t.Run("should update monitor metric and RandomValue gauge metric once", func(t *testing.T) {
		mockAgentMemCache := MockAgentMemCache{
			AgentMemCache: storage.NewAgentMemCache(make(map[string]model.GaugeMetric), make(map[string]model.CounterMetric)),
		}
		monitor := model.Monitor{}
		mockRandom := utils.NewRandom()

		mockAgentMemCache.On("UpdateMonitorMetrics", &monitor).Return(nil)
		mockAgentMemCache.On("UpdateGaugeMetric", mock.Anything).Return(nil)
		UpdateMetrics(&mockAgentMemCache, &monitor, &mockRandom)

		assert.True(t, mockAgentMemCache.AssertNumberOfCalls(t, "UpdateMonitorMetrics", 1), "should update monitor metrics once")
		assert.True(t, mockAgentMemCache.AssertNumberOfCalls(t, "UpdateGaugeMetric", 1), "should update RandomValue gauge metric once")
	})
}

func TestSendMetrics(t *testing.T) {
	t.Run("should update monitor metric and RandomValue gauge metric once", func(t *testing.T) {
		client := &http.Client{}
		url := "/test"
		mockAgentMemCache := MockAgentMemCache{
			AgentMemCache: storage.NewAgentMemCache(make(map[string]model.GaugeMetric), make(map[string]model.CounterMetric)),
		}

		mockAgentMemCache.On("ResetMetricValue", mock.Anything, mock.Anything).Return(nil)
		SendMetrics(client, url, &mockAgentMemCache)

		assert.True(t, mockAgentMemCache.AssertNumberOfCalls(t, "ResetMetricValue", 1), "should reset metrics counter once")
	})
}
