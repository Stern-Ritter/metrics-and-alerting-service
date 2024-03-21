package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model"
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

func TestUpdateMetricHandler(t *testing.T) {
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
			mockStorage := NewMockServerStorage(ctrl)
			mockStorage.
				EXPECT().
				UpdateMetric(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(tt.storageError)
			s := NewServer(mockStorage)

			handler := http.HandlerFunc(s.UpdateMetricHandler)
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

func TestGetMetricHandler(t *testing.T) {
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
			mockStorage := NewMockServerStorage(ctrl)
			mockStorage.
				EXPECT().
				GetMetricValueByTypeAndName(gomock.Any(), gomock.Any()).
				Return(tt.storageReturnValue.value, tt.storageReturnValue.err)
			s := NewServer(mockStorage)

			handler := http.HandlerFunc(s.GetMetricHandler)
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

func TestGetMetricsHandler(t *testing.T) {
	type storageReturnValue struct {
		gauges   map[string]model.GaugeMetric
		counters map[string]model.CounterMetric
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
				gauges: map[string]model.GaugeMetric{
					"metric1": model.NewGauge("metric1", 1.0),
					"metric2": model.NewGauge("metric2", 2.0),
					"metric3": model.NewGauge("metric3", 3.0),
				},
				counters: map[string]model.CounterMetric{
					"metric4": model.NewCounter("metric4", 4),
					"metric5": model.NewCounter("metric5", 5),
					"metric6": model.NewCounter("metric6", 6),
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
			mockStorage := NewMockServerStorage(ctrl)
			mockStorage.
				EXPECT().
				GetMetrics().
				Return(tt.storageReturnValue.gauges, tt.storageReturnValue.counters)
			s := NewServer(mockStorage)

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
