package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	er "github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
)

func addURLParams(req *http.Request, params map[string]string) *http.Request {
	ctx := chi.NewRouteContext()
	for k, v := range params {
		ctx.URLParams.Add(k, v)
	}
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
}

var (
	gaugeInitValue   float64 = 22.2
	counterInitValue int64   = 10
)

func TestUpdateMetricHandlerWithPathVars(t *testing.T) {
	type want struct {
		code int
		body string
	}

	testCases := []struct {
		name       string
		method     string
		url        string
		pathParams map[string]string
		useStorage bool
		want       want
	}{
		{
			name:   "should retrun status code 400 when url contains invalid metric type",
			method: http.MethodPost,
			url:    "/update",
			pathParams: map[string]string{
				"type":  "invalid",
				"name":  "simple",
				"value": "2.0",
			},
			useStorage: false,
			want: want{
				code: http.StatusBadRequest,
				body: "Invalid metric type: invalid\n",
			},
		},
		{
			name:   "should retrun status code 400 when url contains invalid metric value",
			method: http.MethodPost,
			url:    "/update",
			pathParams: map[string]string{
				"type":  "gauge",
				"name":  "simple",
				"value": "two",
			},
			useStorage: false,
			want: want{
				code: http.StatusBadRequest,
				body: "The value for the gauge metric should be of float64 type\n",
			},
		},
		{
			name:   "should retrun status code 200 when url contains valid metric type and value",
			method: http.MethodPost,
			url:    "/update",
			pathParams: map[string]string{
				"type":  "gauge",
				"name":  "simple",
				"value": "2.0",
			},
			useStorage: true,
			want: want{
				code: http.StatusOK,
				body: "",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockStorage(ctrl)
			if tt.useStorage {
				mockStorage.
					EXPECT().
					UpdateMetric(gomock.Any(), gomock.Any()).
					Return(nil)
			}
			config := &server.ServerConfig{}
			logger, err := logger.Initialize("info")
			require.NoError(t, err, "Error init logger")
			metricService := NewMetricService(mockStorage, logger)
			s := NewServer(metricService, config, logger)

			handler := http.HandlerFunc(s.UpdateMetricHandlerWithPathVars)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, nil)
			req = addURLParams(req, tt.pathParams)
			handler.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			require.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, tt.want.code, resp.StatusCode, "Response code didn't match expected")
			if tt.want.body != "" {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				require.Equal(t, tt.want.body, string(body))
			}
		})
	}
}

func TestUpdateMetricHandlerWithBody(t *testing.T) {
	type want struct {
		code int
		body string
	}

	testCases := []struct {
		name                     string
		method                   string
		url                      string
		body                     string
		useStorageUpdateMetric   bool
		storageUpdateMetricError error
		useStorageGetMetric      bool
		metric                   metrics.Metrics
		storageGetMetricError    error
		want                     want
	}{
		{
			name:                   "should return status code 400 when body contains incorrect json",
			method:                 http.MethodPost,
			url:                    "/update",
			body:                   `{ id: "Alloc" type: "gauge" value: 22.2 }`,
			useStorageUpdateMetric: false,
			useStorageGetMetric:    false,
			want: want{
				code: http.StatusBadRequest,
				body: "Error decode request JSON body\n",
			},
		},
		{
			name:                     "should return status code 400 when body contains invalid metric type",
			method:                   http.MethodPost,
			url:                      "/update",
			body:                     `{ "id": "Alloc", "type": "unknown", "value": 22.2 }`,
			useStorageUpdateMetric:   true,
			storageUpdateMetricError: er.NewInvalidMetricType("Invalid metric type: unknown", nil),
			useStorageGetMetric:      false,
			want: want{
				code: http.StatusBadRequest,
				body: "Invalid metric type: unknown\n",
			},
		},
		{
			name:                   "should return status code 200 when body contains valid gauge metric",
			method:                 http.MethodPost,
			url:                    "/update",
			body:                   `{ "id": "Alloc", "type": "gauge", "value": 22.2}`,
			useStorageUpdateMetric: true,
			useStorageGetMetric:    true,
			metric:                 metrics.Metrics{ID: "Alloc", MType: "gauge", Value: &gaugeInitValue},
			want: want{
				code: http.StatusOK,
				body: `{"id":"Alloc","type":"gauge","value":22.2}`,
			},
		},
		{
			name:                   "should return status code 200 when body contains valid counter metric",
			method:                 http.MethodPost,
			url:                    "/update",
			body:                   `{ "id": "PoolCount", "type": "counter", "delta": 10}`,
			useStorageUpdateMetric: true,
			useStorageGetMetric:    true,
			metric:                 metrics.Metrics{ID: "PoolCount", MType: "counter", Delta: &counterInitValue},
			want: want{
				code: http.StatusOK,
				body: `{"id":"PoolCount","type":"counter","delta":10}`,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockStorage(ctrl)
			if tt.useStorageUpdateMetric {
				mockStorage.
					EXPECT().
					UpdateMetric(gomock.Any(), gomock.Any()).
					Return(tt.storageUpdateMetricError)
			}
			if tt.useStorageGetMetric {
				mockStorage.
					EXPECT().
					GetMetric(gomock.Any(), gomock.Any()).
					Return(tt.metric, tt.storageGetMetricError)
			}

			config := &server.ServerConfig{}
			logger, err := logger.Initialize("info")
			require.NoError(t, err, "Error init logger")
			metricService := NewMetricService(mockStorage, logger)
			s := NewServer(metricService, config, logger)

			handler := http.HandlerFunc(s.UpdateMetricHandlerWithBody)
			server := httptest.NewServer(handler)
			defer server.Close()

			req := resty.New().R()
			req.Method = tt.method
			req.Body = tt.body
			req.URL = fmt.Sprintf("%s%s", server.URL, tt.url)

			resp, err := req.Send()
			require.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, tt.want.code, resp.StatusCode(), "Response code didn't match expected")
			assert.Equal(t, tt.want.body, string(resp.Body()))
		})
	}
}

func TestUpdateMetricsBatchHandlerWithBody(t *testing.T) {
	type want struct {
		code int
		body string
	}

	testCases := []struct {
		name                      string
		method                    string
		url                       string
		body                      string
		useStorageUpdateMetrics   bool
		storageUpdateMetricsError error
		want                      want
	}{
		{
			name:                    "should return status code 400 when body contains incorrect json",
			method:                  http.MethodPost,
			url:                     "/updates",
			body:                    `[{ id: "Alloc" type: "gauge" value: 22.2 }]`,
			useStorageUpdateMetrics: false,
			want: want{
				code: http.StatusBadRequest,
				body: "Error decode request JSON body\n",
			},
		},
		{
			name:                      "should return status code 400 when body contains metrics batch with invalid metric type",
			method:                    http.MethodPost,
			url:                       "/updates",
			body:                      `[{ "id": "Alloc", "type": "unknown", "value": 22.2 }]`,
			useStorageUpdateMetrics:   true,
			storageUpdateMetricsError: er.NewInvalidMetricType("Invalid metric type: unknown", nil),
			want: want{
				code: http.StatusBadRequest,
				body: "Invalid metric type: unknown\n",
			},
		},
		{
			name:                    "should return status code 400 when body contains metrics batch with gauge metric with invalid value",
			method:                  http.MethodPost,
			url:                     "/updates",
			body:                    `[{ "id": "Alloc", "type": "gauge", "value": "twenty"}]`,
			useStorageUpdateMetrics: false,
			want: want{
				code: http.StatusBadRequest,
				body: "Error decode request JSON body\n",
			},
		},
		{
			name:                    "should return status code 200 when body contains metrics batch with valid gauge metric",
			method:                  http.MethodPost,
			url:                     "/updates",
			body:                    `[{ "id": "Alloc", "type": "gauge", "value": 22.2}]`,
			useStorageUpdateMetrics: true,
			want: want{
				code: http.StatusOK,
				body: "",
			},
		},
		{
			name:                    "should return status code 400 when body contains metrics batch with counter metric with invalid value",
			method:                  http.MethodPost,
			url:                     "/updates",
			body:                    `[{ "id": "PoolCount", "type": "counter", "delta": 10.11}]`,
			useStorageUpdateMetrics: false,
			want: want{
				code: http.StatusBadRequest,
				body: "Error decode request JSON body\n",
			},
		},
		{
			name:                    "should return status code 200 when body contains metrics batch with valid counter metric",
			method:                  http.MethodPost,
			url:                     "/update",
			body:                    `[{ "id": "PoolCount", "type": "counter", "delta": 10}]`,
			useStorageUpdateMetrics: true,
			want: want{
				code: http.StatusOK,
				body: "",
			},
		},
		{
			name:   "should return status code 200 when body contains metrics batch with multiple valid gauge and counter metrics",
			method: http.MethodPost,
			url:    "/update",
			body: `
			[{ "id": "Alloc", "type": "gauge", "value": 22.2},{ "id": "PoolCount", "type": "counter", "delta": 10}]`,
			useStorageUpdateMetrics: true,
			want: want{
				code: http.StatusOK,
				body: "",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockStorage(ctrl)
			if tt.useStorageUpdateMetrics {
				mockStorage.
					EXPECT().
					UpdateMetrics(gomock.Any(), gomock.Any()).
					Return(tt.storageUpdateMetricsError)
			}

			config := &server.ServerConfig{}
			logger, err := logger.Initialize("info")
			require.NoError(t, err, "Error init logger")
			metricService := NewMetricService(mockStorage, logger)
			s := NewServer(metricService, config, logger)

			handler := http.HandlerFunc(s.UpdateMetricsBatchHandlerWithBody)
			server := httptest.NewServer(handler)
			defer server.Close()

			req := resty.New().R()
			req.Method = tt.method
			req.Body = tt.body
			req.URL = fmt.Sprintf("%s%s", server.URL, tt.url)

			resp, err := req.Send()
			require.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, tt.want.code, resp.StatusCode(), "Response code didn't match expected")
			assert.Equal(t, tt.want.body, string(resp.Body()))
		})
	}
}

func TestGetMetricHandlerWithPathVars(t *testing.T) {
	type want struct {
		code int
		body string
	}

	testCases := []struct {
		name                  string
		method                string
		url                   string
		metric                metrics.Metrics
		storageGetMetricError error
		want                  want
	}{
		{
			name:                  "should retrun status code 404 when url contains not existing metric type",
			method:                http.MethodPost,
			url:                   "/update",
			metric:                metrics.Metrics{},
			storageGetMetricError: er.NewInvalidMetricType("Invalid metric type: unknown", nil),
			want: want{
				code: http.StatusNotFound,
				body: "Invalid metric type: unknown\n",
			},
		},
		{
			name:                  "should retrun status code 404 when url contains not existing metric name",
			method:                http.MethodPost,
			url:                   "/update",
			metric:                metrics.Metrics{},
			storageGetMetricError: er.NewInvalidMetricName("Gauge metric with name: unknown not exists", nil),
			want: want{
				code: http.StatusNotFound,
				body: "Gauge metric with name: unknown not exists\n",
			},
		},
		{
			name:                  "should retrun status code 200 when contains existing metric type and value",
			method:                http.MethodPost,
			url:                   "/update",
			metric:                metrics.Metrics{ID: "PoolCounter", MType: "counter", Delta: &counterInitValue},
			storageGetMetricError: nil,
			want: want{
				code: http.StatusOK,
				body: "10",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockStorage(ctrl)
			mockStorage.
				EXPECT().
				GetMetric(gomock.Any(), gomock.Any()).
				Return(tt.metric, tt.storageGetMetricError)
			config := &server.ServerConfig{}
			logger, err := logger.Initialize("info")
			require.NoError(t, err, "Error init logger")
			metricService := NewMetricService(mockStorage, logger)
			s := NewServer(metricService, config, logger)

			handler := http.HandlerFunc(s.GetMetricHandlerWithPathVars)
			server := httptest.NewServer(handler)
			defer server.Close()

			req := resty.New().R()
			req.Method = tt.method
			req.URL = fmt.Sprintf("%s%s", server.URL, tt.url)

			resp, err := req.Send()

			require.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, tt.want.code, resp.StatusCode(), "Response code didn't match expected")
			assert.Equal(t, tt.want.body, string(resp.Body()))
		})
	}
}

func TestGetMetricHandlerWithBody(t *testing.T) {
	type want struct {
		code int
		body string
	}

	testCases := []struct {
		name                  string
		method                string
		url                   string
		body                  string
		useStorage            bool
		metric                metrics.Metrics
		storageGetMetricError error
		want                  want
	}{
		{
			name:       "should return status code 400 when body contains incorrect json",
			method:     http.MethodPost,
			url:        "/value",
			body:       `{ id: "Alloc" type: "gauge" }`,
			useStorage: false,
			want: want{
				code: http.StatusBadRequest,
				body: "Error decode request JSON body\n",
			},
		},
		{
			name:                  "should return status code 400 when body contains not existing metric type",
			method:                http.MethodPost,
			url:                   "/value",
			body:                  `{ "id": "Alloc", "type": "unknown" }`,
			useStorage:            true,
			metric:                metrics.Metrics{},
			storageGetMetricError: er.NewInvalidMetricType("Invalid metric type: unknown", nil),
			want: want{
				code: http.StatusNotFound,
				body: "Invalid metric type: unknown\n",
			},
		},
		{
			name:       "should return status code 200 when body contains valid gauge metric",
			method:     http.MethodPost,
			url:        "/value",
			body:       `{"id": "Alloc", "type": "gauge"}`,
			useStorage: true,
			metric:     metrics.Metrics{ID: "Alloc", MType: "gauge", Value: &gaugeInitValue},
			want: want{
				code: http.StatusOK,
				body: `{"id":"Alloc","type":"gauge","value":22.2}`,
			},
		},
		{
			name:       "should return status code 200 when body contains valid counter metric",
			method:     http.MethodPost,
			url:        "/value",
			body:       `{"id": "PoolCount", "type": "counter"}`,
			useStorage: true,
			metric:     metrics.Metrics{ID: "PoolCount", MType: "counter", Delta: &counterInitValue},
			want: want{
				code: http.StatusOK,
				body: `{"id":"PoolCount","type":"counter","delta":10}`,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockStorage(ctrl)
			if tt.useStorage {
				mockStorage.
					EXPECT().
					GetMetric(gomock.Any(), gomock.Any()).
					Return(tt.metric, tt.storageGetMetricError)
			}
			config := &server.ServerConfig{}
			logger, err := logger.Initialize("info")
			require.NoError(t, err, "Error init logger")
			metricService := NewMetricService(mockStorage, logger)
			s := NewServer(metricService, config, logger)

			handler := http.HandlerFunc(s.GetMetricHandlerWithBody)
			server := httptest.NewServer(handler)
			defer server.Close()

			req := resty.New().R()
			req.Method = tt.method
			req.Header.Set("Content-type", "application/json")
			req.Body = tt.body
			req.URL = fmt.Sprintf("%s%s", server.URL, tt.url)

			resp, err := req.Send()
			require.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, tt.want.code, resp.StatusCode(), "Response code didn't match expected")
			assert.Equal(t, tt.want.body, string(resp.Body()))
		})
	}
}

func TestGetMetricsHandler(t *testing.T) {
	type storageReturnValue struct {
		gauges   map[string]metrics.GaugeMetric
		counters map[string]metrics.CounterMetric
	}

	type want struct {
		code int
		body string
	}

	testCases := []struct {
		name               string
		method             string
		url                string
		storageReturnValue storageReturnValue
		want               want
	}{
		{
			name:   "should retrun status code 200 when contains existing metric type and value",
			method: http.MethodPost,
			url:    "/",
			storageReturnValue: storageReturnValue{
				gauges: map[string]metrics.GaugeMetric{
					"metric1": metrics.NewGauge("metric1", 1.0),
					"metric2": metrics.NewGauge("metric2", 2.0),
					"metric3": metrics.NewGauge("metric3", 3.0),
				},
				counters: map[string]metrics.CounterMetric{
					"metric4": metrics.NewCounter("metric4", 4),
					"metric5": metrics.NewCounter("metric5", 5),
					"metric6": metrics.NewCounter("metric6", 6),
				},
			},
			want: want{
				code: http.StatusOK,
				body: "metric1,\nmetric2,\nmetric3,\nmetric4,\nmetric5,\nmetric6",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockStorage(ctrl)
			mockStorage.
				EXPECT().
				GetMetrics(gomock.Any()).
				Return(tt.storageReturnValue.gauges, tt.storageReturnValue.counters, nil)
			config := &server.ServerConfig{}
			logger, err := logger.Initialize("info")
			require.NoError(t, err, "Error init logger")
			metricService := NewMetricService(mockStorage, logger)
			s := NewServer(metricService, config, logger)

			handler := http.HandlerFunc(s.GetMetricsHandler)
			server := httptest.NewServer(handler)
			defer server.Close()

			req := resty.New().R()
			req.Method = tt.method
			req.URL = fmt.Sprintf("%s%s", server.URL, tt.url)

			resp, err := req.Send()

			require.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, tt.want.code, resp.StatusCode(), "Response code didn't match expected")
			assert.Equal(t, tt.want.body, string(resp.Body()))
		})
	}
}
