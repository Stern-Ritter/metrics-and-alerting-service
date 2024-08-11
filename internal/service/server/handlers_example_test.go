package server

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	er "github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
)

type ExampleMockStorage struct {
	mock.Mock
}

func (m *ExampleMockStorage) UpdateMetric(ctx context.Context, metric metrics.Metrics) error {
	args := m.Called(ctx, metric)
	return args.Error(0)
}

func (m *ExampleMockStorage) UpdateMetrics(ctx context.Context, metricsBatch []metrics.Metrics) error {
	args := m.Called(ctx, metricsBatch)
	return args.Error(0)
}

func (m *ExampleMockStorage) GetMetric(ctx context.Context, metric metrics.Metrics) (metrics.Metrics, error) {
	args := m.Called(ctx, metric)
	return args.Get(0).(metrics.Metrics), args.Error(1)
}

func (m *ExampleMockStorage) GetMetrics(ctx context.Context) (map[string]metrics.GaugeMetric,
	map[string]metrics.CounterMetric, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]metrics.GaugeMetric), args.Get(1).(map[string]metrics.CounterMetric),
		args.Error(2)
}

func (m *ExampleMockStorage) Restore(fName string) error {
	args := m.Called(fName)
	return args.Error(0)
}

func (m *ExampleMockStorage) Save(fName string) error {
	args := m.Called(fName)
	return args.Error(0)
}

func (m *ExampleMockStorage) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *ExampleMockStorage) Close() error {
	args := m.Called()
	return args.Error(0)
}

// ExampleServer_UpdateMetricHandlerWithPathVars shows how to update a metric using request path variables.
func ExampleServer_UpdateMetricHandlerWithPathVars() {
	r := chi.NewRouter()
	storage := ExampleMockStorage{}
	cfg := &config.ServerConfig{}
	log, err := logger.Initialize("info")
	if err != nil {
		return
	}
	service := NewMetricService(&storage, log)
	server := &Server{
		MetricService: service,
		Config:        cfg,
		Logger:        log,
	}
	r.Put("/update/{type}/{name}/{value}", server.UpdateMetricHandlerWithPathVars)

	// valid metric name, type, value
	req := httptest.NewRequest(http.MethodPut, "/update/gauge/first/11.1", nil)
	storage.On("UpdateMetric", mock.Anything, mock.Anything).Return(nil).Once()
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// invalid metric type
	req = httptest.NewRequest(http.MethodPut, "/update/unknown/first/11.1", nil)
	storage.
		On("UpdateMetric", mock.Anything, mock.Anything).
		Return(er.NewInvalidMetricType("invalid metric type", nil)).Once()
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// invalid metric value
	req = httptest.NewRequest(http.MethodPut, "/update/gauge/first/eleven", nil)

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// Output:
	// 200
	// 400
	// 400
}

// ExampleServer_UpdateMetricHandlerWithBody shows how to update a metric using request body.
func ExampleServer_UpdateMetricHandlerWithBody() {
	r := chi.NewRouter()
	storage := ExampleMockStorage{}
	cfg := &config.ServerConfig{}
	log, err := logger.Initialize("info")
	if err != nil {
		return
	}
	service := NewMetricService(&storage, log)
	server := &Server{
		MetricService: service,
		Config:        cfg,
		Logger:        log,
	}
	r.Post("/update", server.UpdateMetricHandlerWithBody)

	// valid metric name, type, value
	body := bytes.NewReader([]byte(`{"id": "first", "type": "gauge", "value": 11.1}`))
	req := httptest.NewRequest(http.MethodPost, "/update", body)
	req.Header.Set("Content-Type", "application/json")
	storage.On("UpdateMetric", mock.Anything, mock.Anything).Return(nil).Once()
	storage.On("GetMetric", mock.Anything, mock.Anything).Return(metrics.Metrics{}, nil).Once()
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// invalid metric type
	body = bytes.NewReader([]byte(`{"id": "first", "type": "unknown", "value": 11.1}`))
	req = httptest.NewRequest(http.MethodPost, "/update", body)
	req.Header.Set("Content-Type", "application/json")
	storage.On("UpdateMetric", mock.Anything, mock.Anything).
		Return(er.NewInvalidMetricType("invalid metric type", nil)).Once()
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// invalid metric value
	body = bytes.NewReader([]byte(`{"id": "first", "type": "gauge", "value": "eleven"}`))
	req = httptest.NewRequest(http.MethodPost, "/update", body)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// Output:
	// 200
	// 400
	// 400
}

// ExampleServer_UpdateMetricsBatchHandlerWithBody shows how to update a batch of metrics using request body.
func ExampleServer_UpdateMetricsBatchHandlerWithBody() {
	r := chi.NewRouter()
	storage := ExampleMockStorage{}
	cfg := &config.ServerConfig{}
	log, err := logger.Initialize("info")
	if err != nil {
		return
	}
	service := NewMetricService(&storage, log)
	server := &Server{
		MetricService: service,
		Config:        cfg,
		Logger:        log,
	}
	r.Post("/updates", server.UpdateMetricsBatchHandlerWithBody)

	// valid metric name, type, value
	body := bytes.NewReader([]byte(`[{"id": "first", "type": "counter", "value": 10}, {"id": "second", "type": "gauge", "value": 11.1}]`))
	req := httptest.NewRequest(http.MethodPost, "/updates", body)
	req.Header.Set("Content-Type", "application/json")
	storage.On("UpdateMetrics", mock.Anything, mock.Anything).Return(nil).Once()
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// invalid metric type
	body = bytes.NewReader([]byte(`[{"id": "first", "type": "unknown", "value": 10}, {"id": "second", "type": "gauge", "value": 11.1}]`))
	req = httptest.NewRequest(http.MethodPost, "/updates", body)
	req.Header.Set("Content-Type", "application/json")
	storage.On("UpdateMetrics", mock.Anything, mock.Anything).
		Return(er.NewInvalidMetricType("invalid metric type", nil)).Once()
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// invalid metric value
	body = bytes.NewReader([]byte(`[{"id": "first", "type": "unknown", "value": "ten"}, {"id": "second", "type": "gauge", "value": "eleven"}]`))
	req = httptest.NewRequest(http.MethodPost, "/updates", body)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// Output:
	// 200
	// 400
	// 400
}

// ExampleServer_GetMetricHandlerWithPathVars shows how to get a metric value by type and name using request path variables.
func ExampleServer_GetMetricHandlerWithPathVars() {
	r := chi.NewRouter()
	storage := ExampleMockStorage{}
	cfg := &config.ServerConfig{}
	log, err := logger.Initialize("info")
	if err != nil {
		return
	}
	service := NewMetricService(&storage, log)
	server := &Server{
		MetricService: service,
		Config:        cfg,
		Logger:        log,
	}
	r.Get("/value/{type}/{name}", server.GetMetricHandlerWithPathVars)

	// valid metric name, type, value
	req := httptest.NewRequest(http.MethodGet, "/value/gauge/first", nil)
	metricValue := 42.0
	storage.On("GetMetric", mock.Anything, mock.Anything).
		Return(metrics.Metrics{ID: "first", MType: "gauge", Value: &metricValue}, nil).Once()
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

	// not found metric name
	req = httptest.NewRequest(http.MethodGet, "/value/gauge/unknown", nil)
	storage.On("GetMetric", mock.Anything, mock.Anything).
		Return(metrics.Metrics{}, er.NewInvalidMetricName("invalid metric name", nil)).Once()
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// not found metric type
	req = httptest.NewRequest(http.MethodGet, "/value/unknown/first", nil)
	storage.On("GetMetric", mock.Anything, mock.Anything).
		Return(metrics.Metrics{}, er.NewInvalidMetricType("invalid metric type", nil)).Once()
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// Output:
	// 200
	// 42
	// 404
	// 404
}

// ExampleServer_GetMetricHandlerWithBody shows how to get a metric value using the request body.
func ExampleServer_GetMetricHandlerWithBody() {
	r := chi.NewRouter()
	storage := ExampleMockStorage{}
	cfg := &config.ServerConfig{}
	log, err := logger.Initialize("info")
	if err != nil {
		return
	}
	service := NewMetricService(&storage, log)
	server := &Server{
		MetricService: service,
		Config:        cfg,
		Logger:        log,
	}
	r.Post("/value", server.GetMetricHandlerWithBody)

	// valid metric name, type, value
	body := bytes.NewReader([]byte(`{"id": "first", "type": "gauge"}`))
	req := httptest.NewRequest(http.MethodPost, "/value", body)
	req.Header.Set("Content-Type", "application/json")
	storage.On("GetMetric", mock.Anything, mock.Anything).Return(metrics.Metrics{}, nil).Once()
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// not found metric name
	body = bytes.NewReader([]byte(`{"id": "unknown", "type": "gauge"}`))
	req = httptest.NewRequest(http.MethodPost, "/value", body)
	req.Header.Set("Content-Type", "application/json")
	storage.On("GetMetric", mock.Anything, mock.Anything).
		Return(metrics.Metrics{}, er.NewInvalidMetricName("invalid metric name", nil)).Once()
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// not found metric type
	body = bytes.NewReader([]byte(`{"id": "first", "type": "unknown"}`))
	req = httptest.NewRequest(http.MethodPost, "/value", body)
	req.Header.Set("Content-Type", "application/json")
	storage.On("GetMetric", mock.Anything, mock.Anything).
		Return(metrics.Metrics{}, er.NewInvalidMetricType("invalid metric type", nil)).Once()
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// Output:
	// 200
	// 404
	// 404
}

// ExampleServer_GetMetricsHandler shows how to get all metrics.
func ExampleServer_GetMetricsHandler() {
	r := chi.NewRouter()
	storage := ExampleMockStorage{}
	cfg := &config.ServerConfig{}
	log, err := logger.Initialize("info")
	if err != nil {
		return
	}
	service := NewMetricService(&storage, log)
	server := &Server{
		MetricService: service,
		Config:        cfg,
		Logger:        log,
	}
	r.Get("/", server.GetMetricsHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	storage.On("GetMetrics", mock.Anything).
		Return(make(map[string]metrics.GaugeMetric), make(map[string]metrics.CounterMetric), nil).Once()
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// Output: 200
}

// ExampleServer_PingDatabaseHandler shows	 how to check the connection to the database.
func ExampleServer_PingDatabaseHandler() {
	r := chi.NewRouter()
	storage := ExampleMockStorage{}
	cfg := &config.ServerConfig{}
	log, err := logger.Initialize("info")
	if err != nil {
		return
	}
	service := NewMetricService(&storage, log)
	server := &Server{
		MetricService: service,
		Config:        cfg,
		Logger:        log,
	}
	r.Get("/ping", server.PingDatabaseHandler)

	// database is available
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	storage.On("Ping", mock.Anything).Return(nil).Once()
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// database in unavailable
	req = httptest.NewRequest(http.MethodGet, "/ping", nil)
	storage.On("Ping", mock.Anything).Return(errors.New("the database is disabled")).Once()
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// Output:
	// 200
	// 500
}
