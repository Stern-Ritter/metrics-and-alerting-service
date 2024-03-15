package server

import (
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model"
	"github.com/go-chi/chi"
)

func (s *Server) UpdateMetricHandler(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "type")
	metricName := chi.URLParam(req, "name")
	metricValue := chi.URLParam(req, "value")

	err := s.storage.UpdateMetric(metricType, metricName, metricValue)
	switch err.(type) {
	case errors.InvalidMetricType, errors.InvalidMetricValue:
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (s *Server) GetMetricHandler(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "type")
	metricName := chi.URLParam(req, "name")

	body, err := s.storage.GetMetricValueByTypeAndName(metricType, metricName)
	switch err.(type) {
	case errors.InvalidMetricType, errors.InvalidMetricName:
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
	res.Header().Set("Content-type", "text/plain")
	res.WriteHeader(http.StatusOK)
	_, err = io.WriteString(res, body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) GetMetricsHandler(res http.ResponseWriter, req *http.Request) {
	gauges, counters := s.storage.GetMetrics()
	body := getMetricsString(gauges, counters)

	res.Header().Set("Content-type", "text/html")
	res.WriteHeader(http.StatusOK)
	_, err := io.WriteString(res, body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func getMetricsString(gauges map[string]model.GaugeMetric, counters map[string]model.CounterMetric) string {
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
