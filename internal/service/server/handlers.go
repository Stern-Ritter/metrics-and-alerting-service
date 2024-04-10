package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/go-chi/chi"
)

func (s *Server) UpdateMetricHandlerWithPathVars(res http.ResponseWriter, req *http.Request) {
	mName := chi.URLParam(req, "name")
	mType := chi.URLParam(req, "type")
	mValue := chi.URLParam(req, "value")

	err := s.MetricService.UpdateMetricWithPathVars(req.Context(), mName, mType, mValue,
		s.isSyncSaveStorageState(), s.Config.StorageFilePath)

	switch err.(type) {
	case errors.InvalidMetricType, errors.InvalidMetricValue:
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (s *Server) UpdateMetricHandlerWithBody(res http.ResponseWriter, req *http.Request) {
	metric, err := decodeMetrics(req.Body)
	if err != nil {
		http.Error(res, "Error decode request JSON body", http.StatusBadRequest)
		return
	}

	updatedMetric, err := s.MetricService.UpdateMetricWithBody(req.Context(), metric, s.isSyncSaveStorageState(),
		s.Config.StorageFilePath)
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

func (s *Server) GetMetricHandlerWithPathVars(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "type")
	metricName := chi.URLParam(req, "name")

	value, err := s.MetricService.GetMetricValueByTypeAndName(req.Context(), metricType, metricName)
	switch err.(type) {
	case errors.InvalidMetricType, errors.InvalidMetricName:
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	res.Header().Set("Content-type", "text/plain")
	res.WriteHeader(http.StatusOK)
	_, err = io.WriteString(res, value)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) GetMetricHandlerWithBody(res http.ResponseWriter, req *http.Request) {
	metric, err := decodeMetrics(req.Body)
	if err != nil {
		http.Error(res, "Error decode request JSON body", http.StatusBadRequest)
		return
	}

	metric, err = s.MetricService.GetMetricHandlerWithBody(req.Context(), metric)

	switch err.(type) {
	case errors.InvalidMetricType, errors.InvalidMetricName:
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(metric)
	if err != nil {
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
	}
	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(resp)
	if err != nil {
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
	}
}

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

func (s *Server) isSyncSaveStorageState() bool {
	isFileStorageEnabled := len(strings.TrimSpace(s.Config.StorageFilePath)) != 0
	isSyncSaveStorageState := s.Config.StoreInterval == 0
	return isFileStorageEnabled && isSyncSaveStorageState
}
