package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi"

	er "github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
)

// UpdateMetricHandlerWithPathVars updates a metric using request path variables.
func (s *Server) UpdateMetricHandlerWithPathVars(res http.ResponseWriter, req *http.Request) {
	mName := chi.URLParam(req, "name")
	mType := chi.URLParam(req, "type")
	mValue := chi.URLParam(req, "value")

	err := s.MetricService.UpdateMetricWithPathVars(req.Context(), mName, mType, mValue,
		s.isSyncSaveStorageState(), s.Config.FileStoragePath)

	if err != nil {
		var invalidMetricType er.InvalidMetricType
		var invalidMetricValue er.InvalidMetricValue
		if errors.As(err, &invalidMetricType) || errors.As(err, &invalidMetricValue) {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

// UpdateMetricHandlerWithBody updates a metric using request body.
func (s *Server) UpdateMetricHandlerWithBody(res http.ResponseWriter, req *http.Request) {
	metric, err := decodeMetrics(req.Body)
	if err != nil {
		http.Error(res, "Error decode request JSON body", http.StatusBadRequest)
		return
	}

	updatedMetric, err := s.MetricService.UpdateMetricWithBody(req.Context(), metric, s.isSyncSaveStorageState(),
		s.Config.FileStoragePath)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(updatedMetric)
	if err != nil {
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
	}
	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(resp)
	if err != nil {
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
	}
}

// UpdateMetricsBatchHandlerWithBody updates a batch of metrics using the request body.
func (s *Server) UpdateMetricsBatchHandlerWithBody(res http.ResponseWriter, req *http.Request) {
	metrics, err := decodeMetricsBatch(req.Body)
	if err != nil {
		http.Error(res, "Error decode request JSON body", http.StatusBadRequest)
		return
	}

	err = s.MetricService.UpdateMetricsBatchWithBody(req.Context(), metrics,
		s.isSyncSaveStorageState(), s.Config.FileStoragePath)

	if err != nil {
		var invalidMetricType er.InvalidMetricType
		var invalidMetricValue er.InvalidMetricValue
		if errors.As(err, &invalidMetricType) || errors.As(err, &invalidMetricValue) {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

// GetMetricHandlerWithPathVars  returns the value of a metric by type and name using request path variables.
func (s *Server) GetMetricHandlerWithPathVars(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "type")
	metricName := chi.URLParam(req, "name")

	value, err := s.MetricService.GetMetricValueByTypeAndName(req.Context(), metricType, metricName)

	if err != nil {
		var invalidMetricType er.InvalidMetricType
		var invalidMetricName er.InvalidMetricName
		if errors.As(err, &invalidMetricType) || errors.As(err, &invalidMetricName) {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-type", "text/plain")
	res.WriteHeader(http.StatusOK)
	_, err = io.WriteString(res, value)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

// GetMetricHandlerWithBody returns a metric by type and name defined in the request body.
func (s *Server) GetMetricHandlerWithBody(res http.ResponseWriter, req *http.Request) {
	metric, err := decodeMetrics(req.Body)
	if err != nil {
		http.Error(res, "Error decode request JSON body", http.StatusBadRequest)
		return
	}

	savedMetric, err := s.MetricService.GetMetricValueWithBody(req.Context(), metric)

	if err != nil {
		var invalidMetricType er.InvalidMetricType
		var invalidMetricName er.InvalidMetricName
		if errors.As(err, &invalidMetricType) || errors.As(err, &invalidMetricName) {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(savedMetric)
	if err != nil {
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
	}
	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(resp)
	if err != nil {
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
	}
}

// GetMetricsHandler returns all metrics.
func (s *Server) GetMetricsHandler(res http.ResponseWriter, req *http.Request) {
	gauges, counters, err := s.MetricService.GetMetrics(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	body := getMetricsString(gauges, counters)

	res.Header().Set("Content-type", "text/html")
	_, err = io.WriteString(res, body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
}

// PingDatabaseHandler checks the connection to the database and return connection status.
func (s *Server) PingDatabaseHandler(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second)
	defer cancel()

	err := s.MetricService.PingDatabase(ctx)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func getMetricsString(gauges map[string]metrics.GaugeMetric, counters map[string]metrics.CounterMetric) string {
	metricsNames := make([]string, 0)
	for _, metric := range gauges {
		metricsNames = append(metricsNames, metric.Name)
	}
	for _, metric := range counters {
		metricsNames = append(metricsNames, metric.Name)
	}
	sort.Strings(metricsNames)

	return strings.Join(metricsNames, ",\n")
}

func decodeMetrics(source io.ReadCloser) (metrics.Metrics, error) {
	metric := metrics.Metrics{}
	var buf bytes.Buffer
	_, err := buf.ReadFrom(source)
	if err != nil {
		return metric, err
	}

	err = json.Unmarshal(buf.Bytes(), &metric)
	return metric, err
}

func decodeMetricsBatch(source io.ReadCloser) ([]metrics.Metrics, error) {
	metricsBatch := make([]metrics.Metrics, 0)
	var buf bytes.Buffer
	_, err := buf.ReadFrom(source)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf.Bytes(), &metricsBatch)
	if err != nil {
		return nil, err
	}

	return metricsBatch, nil
}

func (s *Server) isSyncSaveStorageState() bool {
	isFileStorageEnabled := len(strings.TrimSpace(s.Config.FileStoragePath)) != 0
	isSyncSaveStorageState := s.Config.StoreInterval == 0
	return isFileStorageEnabled && isSyncSaveStorageState
}
