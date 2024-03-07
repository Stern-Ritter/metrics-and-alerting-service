package transport

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockStorage struct {
	storage.ServerMemStorage
	mock.Mock
}

func (m *MockStorage) UpdateMetric(metricType, metricName, metricValue string) error {
	args := m.Called(metricType, metricName, metricValue)
	return args.Error(0)
}

func TestIsRequestMethodAllowed(t *testing.T) {
	testCases := []struct {
		name           string
		allowedMethods []string
		method         string
		want           bool
	}{
		{
			name:           "should return true when method is allowed #1",
			allowedMethods: []string{http.MethodPost},
			method:         http.MethodPost,
			want:           true,
		},
		{
			name:           "should return true when method is allowed #2",
			allowedMethods: []string{http.MethodGet, http.MethodPost},
			method:         http.MethodPost,
			want:           true,
		},

		{
			name:           "should return false when method is not allowed #1",
			allowedMethods: []string{},
			method:         http.MethodPut,
			want:           false,
		},
		{
			name:           "should return false when method is not allowed #2",
			allowedMethods: []string{http.MethodGet, http.MethodPost},
			method:         http.MethodPut,
			want:           false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := isRequestMethodAllowed(tt.method, tt.allowedMethods)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParsePathVariables(t *testing.T) {
	type want struct {
		metricType  string
		metricName  string
		metricValue string
	}
	testCases := []struct {
		name   string
		path   string
		prefix string
		want   want
		err    bool
	}{
		{
			name:   "should return error when path doesn`t contain all required path variables",
			path:   "/update/type/name",
			prefix: "/update",
			want:   want{},
			err:    true,
		},
		{
			name:   "should return error when path is empty",
			path:   "",
			prefix: "/update",
			want:   want{},
			err:    true,
		},
		{
			name:   "should return error when one of required path variables is empty #1",
			path:   "/update/type/name/",
			prefix: "/update",
			want:   want{},
			err:    true,
		},

		{
			name:   "should return error when one of required path variables is empty #2",
			path:   "/update/type/ /value",
			prefix: "/update",
			want:   want{},
			err:    true,
		},
		{
			name:   "should return error when one of required path variables is empty #3",
			path:   "/update/ /name/value",
			prefix: "/update",
			want:   want{},
			err:    true,
		},
		{
			name:   "should return parsed path variables when path contains all required path variables #1",
			path:   "/update/type1/name1/value1",
			prefix: "/update",
			want: want{
				metricType:  "type1",
				metricName:  "name1",
				metricValue: "value1",
			},
			err: false,
		},
		{
			name:   "should return parsed path variables when path contains all required path variables #2",
			path:   "/update/type2/name2/value2",
			prefix: "/update",
			want: want{
				metricType:  "type2",
				metricName:  "name2",
				metricValue: "value2",
			},
			err: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			gotMetricType, gotMetricName, gotMetricValue, gotErr := parsePathVariables(tt.path, tt.prefix)
			if tt.err {
				assert.Error(t, gotErr)
			} else {
				assert.Equal(t, tt.want.metricType, gotMetricType)
				assert.Equal(t, tt.want.metricName, gotMetricName)
				assert.Equal(t, tt.want.metricValue, gotMetricValue)
			}
		})
	}
}

func TestUpdateMetricHandler(t *testing.T) {
	const (
		validPath                   = "/update/validType/name/validValue"
		invalidPath                 = "/update/validType//validValue"
		invalidMetricTypePath       = "/update/invalidType/name/validValue"
		invalidMetricValuePath      = "/update/validType/name/invalidValue"
		invalidMetricTypeErrorText  = "Storage error text for invalid metric type"
		invalidMetricValueErrorText = "Storage error text for invalid metric value"
	)
	type want struct {
		code int
		body string
	}

	testCases := []struct {
		name         string
		mathod       string
		path         string
		storageError error
		want         want
	}{
		{
			name:         "should retrun status code 405 when method is not allowed",
			mathod:       http.MethodGet,
			path:         validPath,
			storageError: nil,
			want: want{
				code: http.StatusMethodNotAllowed,
				body: "Only POST requests are allowed.\n",
			},
		},
		{
			name:         "should retrun status code 404 when path doesn`t contain all required path variables",
			mathod:       http.MethodPost,
			path:         invalidPath,
			storageError: nil,
			want: want{
				code: http.StatusNotFound,
				body: "The resource you requested has not been found at the specified address\n",
			},
		},
		{
			name:         "should retrun status code 400 when path contains invalid metric type",
			mathod:       http.MethodPost,
			path:         invalidMetricTypePath,
			storageError: errors.NewInvalidMetricType(invalidMetricTypeErrorText, nil),
			want: want{
				code: http.StatusBadRequest,
				body: fmt.Sprintf("%s\n", invalidMetricTypeErrorText),
			},
		},
		{
			name:         "should retrun status code 400 when path contains invalid metric value",
			mathod:       http.MethodPost,
			path:         invalidMetricValuePath,
			storageError: errors.NewInvalidMetricValue(invalidMetricValueErrorText, nil),
			want: want{
				code: http.StatusBadRequest,
				body: fmt.Sprintf("%s\n", invalidMetricValueErrorText),
			},
		},
		{
			name:         "should retrun status code 200 when path contains all required path variables #1",
			mathod:       http.MethodPost,
			path:         validPath,
			storageError: nil,
			want: want{
				code: http.StatusOK,
				body: "",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.mathod, tt.path, nil)
			response := httptest.NewRecorder()

			mockStorage := MockStorage{}
			mockStorage.On(
				"UpdateMetric",
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).Return(tt.storageError)
			handler := UpdateMetricHandler(&mockStorage)

			handler(response, request)

			res := response.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			body, err := io.ReadAll(res.Body)
			res.Body.Close()

			require.NoError(t, err)
			assert.Equal(t, tt.want.body, string(body))
		})
	}
}
