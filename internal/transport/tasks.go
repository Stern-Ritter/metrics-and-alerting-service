package transport

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

func UpdateMetrics(cache storage.AgentCache, monitor *model.Monitor, rand *utils.Random) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	monitor.Update(&ms)

	cache.UpdateMonitorMetrics(monitor)

	randomValue, err := rand.Float(0.1, 99.99)
	if err == nil {
		cache.UpdateGaugeMetric(model.NewGauge("RandomValue", randomValue))
	}
}

func SendMetrics(client *http.Client, url string, cache storage.AgentCache) {
	gauges, counters := cache.GetMetrics()
	cache.ResetMetricValue(string(model.Gauge), "PollCount")

	for _, metric := range gauges {
		_, err := sendPostRequest(client, url, "text/plain", string(metric.Type), metric.Name, formatGaugeMetricValue(metric.GetValue()))
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, metric := range counters {
		_, err := sendPostRequest(client, url, "text/plain", string(metric.Type), metric.Name, formatCounterMetricValue(metric.GetValue()))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func sendPostRequest(c *http.Client, url, contentType, metricType, metricName, metricValue string) (string, error) {
	endpoint := makeEndpoint(url, metricType, metricName, metricValue)

	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", contentType)
	res, err := c.Do(req)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func makeEndpoint(parts ...string) string {
	return strings.Join(parts, "/")
}

func formatGaugeMetricValue(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func formatCounterMetricValue(value int64) string {
	return strconv.Itoa(int(value))
}

func SetInterval(ctx context.Context, wg *sync.WaitGroup, task func(), interval time.Duration) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				wg.Done()
				return
			default:
				task()
				time.Sleep(interval)
			}
		}
	}()
}
