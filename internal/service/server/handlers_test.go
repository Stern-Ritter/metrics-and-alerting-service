package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	validURL                    = "/update/counter/name/2.0"
	invalidMetricTypeURL        = "/update/invalidType/name/2.0"
	invalidMetricValueURL       = "/update/counter/name/two"
	invalidMetricTypeErrorText  = "Storage error text for invalid metric type"
	invalidMetricValueErrorText = "Storage error text for invalid metric value"
)

func TestUpdateMetricHandlerWithPathVars(t *testing.T) {
	type want struct {
		code int
		body string
	}

	testCases := []struct {
		name         string
		method       string
		url          string
		storageError error
		want         want
	}{
		{
			name:         "should retrun status code 400 when url contains invalid metric type",
			method:       http.MethodPost,
			url:          invalidMetricTypeURL,
			storageError: errors.NewInvalidMetricType(invalidMetricTypeErrorText, nil),
			want: want{
				code: http.StatusBadRequest,
				body: fmt.Sprintf("%s\n", invalidMetricTypeErrorText),
			},
		},
		{
			name:         "should retrun status code 400 when url contains invalid metric value",
			method:       http.MethodPost,
			url:          invalidMetricValueURL,
			storageError: errors.NewInvalidMetricValue(invalidMetricValueErrorText, nil),
			want: want{
				code: http.StatusBadRequest,
				body: fmt.Sprintf("%s\n", invalidMetricValueErrorText),
			},
		},
		{
			name:         "should retrun status code 200 when url contains valid metric type and value",
			method:       http.MethodPost,
			url:          validURL,
			storageError: nil,
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

			mockStorage := NewMockServerStorage(ctrl)
			mockStorage.
				EXPECT().
				UpdateMetric(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(tt.storageError)
			config := &server.ServerConfig{}
			logger, err := logger.Initialize("info")
			require.NoError(t, err, "Error init logger")
			s := NewServer(mockStorage, config, logger)

			handler := http.HandlerFunc(s.UpdateMetricHandlerWithPathVars)
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

func TestUpdateMetricHandlerWithBody(t *testing.T) {
	type want struct {
		code int
		body string
	}

	testCases := []struct {
		name           string
		method         string
		url            string
		body           string
		useStorage     bool
		returnedMetric metrics.GaugeMetric
		returnedError  error
		want           want
	}{
		{
			name:       "should return status code 500 when body contains incorrect json",
			method:     http.MethodPost,
			url:        "/update",
			body:       `{ id: "Alloc" type: "gauge" value: 22.2 }`,
			useStorage: false,
			want: want{
				code: http.StatusBadRequest,
				body: "Error decode request JSON body\n",
			},
		},
		{
			name:           "should return status code 400 when body contains invalid metric type",
			method:         http.MethodPost,
			url:            "/update",
			body:           `{ "id": "Alloc", "type": "unknown", "value": 22.2 }`,
			useStorage:     false,
			returnedMetric: metrics.GaugeMetric{},
			returnedError:  nil,
			want: want{
				code: http.StatusBadRequest,
				body: "Invalid metric type: unknown\n",
			},
		},
		{
			name:           "should return status code 200 when body contains valid metric type",
			method:         http.MethodPost,
			url:            "/update",
			body:           `{ "id": "Alloc", "type": "gauge", "value": 22.2}`,
			useStorage:     true,
			returnedMetric: metrics.NewGauge("Alloc", 22.2),
			returnedError:  nil,
			want: want{
				code: http.StatusOK,
				body: `{"id":"Alloc","type":"gauge","value":22.2}`,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockServerStorage(ctrl)
			if tt.useStorage {
				mockStorage.
					EXPECT().
					UpdateGaugeMetric(gomock.Any()).
					Return(tt.returnedMetric, tt.returnedError)
			}
			config := &server.ServerConfig{}
			logger, err := logger.Initialize("info")
			require.NoError(t, err, "Error init logger")
			s := NewServer(mockStorage, config, logger)

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

func TestGetMetricHandlerWithPathVars(t *testing.T) {
	type storageReturnValue struct {
		value string
		err   error
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
			name:   "should retrun status code 404 when url contains not existing metric type",
			method: http.MethodPost,
			url:    invalidMetricTypeURL,
			storageReturnValue: storageReturnValue{
				value: "",
				err:   errors.NewInvalidMetricType(invalidMetricTypeErrorText, nil),
			},
			want: want{
				code: http.StatusNotFound,
				body: fmt.Sprintf("%s\n", invalidMetricTypeErrorText),
			},
		},
		{
			name:   "should retrun status code 404 when url contains not existing metric name",
			method: http.MethodPost,
			url:    invalidMetricValueURL,
			storageReturnValue: storageReturnValue{
				value: "",
				err:   errors.NewInvalidMetricName(invalidMetricValueErrorText, nil),
			},
			want: want{
				code: http.StatusNotFound,
				body: fmt.Sprintf("%s\n", invalidMetricValueErrorText),
			},
		},
		{
			name:   "should retrun status code 200 when contains existing metric type and value",
			method: http.MethodPost,
			url:    validURL,
			storageReturnValue: storageReturnValue{
				value: "1",
				err:   nil,
			},
			want: want{
				code: http.StatusOK,
				body: "1",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockServerStorage(ctrl)
			mockStorage.
				EXPECT().
				GetMetricValueByTypeAndName(gomock.Any(), gomock.Any()).
				Return(tt.storageReturnValue.value, tt.storageReturnValue.err)
			config := &server.ServerConfig{}
			logger, err := logger.Initialize("info")
			require.NoError(t, err, "Error init logger")
			s := NewServer(mockStorage, config, logger)

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
		name           string
		method         string
		url            string
		body           string
		useStorage     bool
		returnedMetric metrics.GaugeMetric
		returnedError  error
		want           want
	}{
		{
			name:       "should return status code 500 when body contains incorrect json",
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
			name:           "should return status code 400 when body contains invalid metric type",
			method:         http.MethodPost,
			url:            "/value",
			body:           `{ "id": "Alloc", "type": "unknown" }`,
			useStorage:     false,
			returnedMetric: metrics.GaugeMetric{},
			returnedError:  nil,
			want: want{
				code: http.StatusBadRequest,
				body: "Invalid metric type: unknown\n",
			},
		},
		{
			name:           "should return status code 200 when body contains valid metric type",
			method:         http.MethodPost,
			url:            "/value",
			body:           `{"id": "Alloc", "type": "gauge"}`,
			useStorage:     true,
			returnedMetric: metrics.NewGauge("Alloc", 22.2),
			returnedError:  nil,
			want: want{
				code: http.StatusOK,
				body: `{"id":"Alloc","type":"gauge","value":22.2}`,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := NewMockServerStorage(ctrl)
			if tt.useStorage {
				mockStorage.
					EXPECT().
					GetGaugeMetric(gomock.Any()).
					Return(tt.returnedMetric, tt.returnedError)
			}
			config := &server.ServerConfig{}
			logger, err := logger.Initialize("info")
			require.NoError(t, err, "Error init logger")
			s := NewServer(mockStorage, config, logger)

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
			url:    validURL,
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

			mockStorage := NewMockServerStorage(ctrl)
			mockStorage.
				EXPECT().
				GetMetrics().
				Return(tt.storageReturnValue.gauges, tt.storageReturnValue.counters)
			config := &server.ServerConfig{}
			logger, err := logger.Initialize("info")
			require.NoError(t, err, "Error init logger")
			s := NewServer(mockStorage, config, logger)

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
