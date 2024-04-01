package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func (s *Server) UpdateMetricHandlerWithPathVars(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "type")
	metricName := chi.URLParam(req, "name")
	metricValue := chi.URLParam(req, "value")

	err := s.Storage.UpdateMetric(metricType, metricName, metricValue)
	switch err.(type) {
	case errors.InvalidMetricType, errors.InvalidMetricValue:
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if s.isSyncSaveStorageState() {
		err := s.Storage.Save(s.Config.StorageFilePath)
		if err != nil {
			s.Logger.Error(err.Error(), zap.String("event", "sync save to file storage"))
		} else {
			s.Logger.Info("Success sync save to file storage", zap.String("event", "sync save to file storage"))
		}
	}
	res.WriteHeader(http.StatusOK)
}

func (s *Server) UpdateMetricHandlerWithBody(res http.ResponseWriter, req *http.Request) {
	metric := metrics.Metrics{}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(req.Body)

	if err != nil {
		http.Error(res, "Error decode request JSON body", http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(res, "Error decode request JSON body", http.StatusBadRequest)
		return
	}

	switch metrics.MetricType(metric.MType) {
	case metrics.Gauge:
		updatedMetric, err := s.Storage.UpdateGaugeMetric(metrics.MetricsToGaugeMetric(metric))
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		value := updatedMetric.GetValue()
		metric.Value = &value

	case metrics.Counter:
		updatedMetric, err := s.Storage.UpdateCounterMetric(metrics.MetricsToCounterMetric(metric))
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		delta := updatedMetric.GetValue()
		metric.Delta = &delta

	default:
		http.Error(res, fmt.Sprintf("Invalid metric type: %s", metric.MType), http.StatusBadRequest)
		return
	}

	if s.isSyncSaveStorageState() {
		err := s.Storage.Save(s.Config.StorageFilePath)
		if err != nil {
			s.Logger.Error(err.Error(), zap.String("event", "sync save to file storage"))
		} else {
			s.Logger.Info("Success sync save to file storage", zap.String("event", "sync save to file storage"))
		}
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

func (s *Server) GetMetricHandlerWithPathVars(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "type")
	metricName := chi.URLParam(req, "name")

	body, err := s.Storage.GetMetricValueByTypeAndName(metricType, metricName)
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

func (s *Server) GetMetricHandlerWithBody(res http.ResponseWriter, req *http.Request) {
	metric := metrics.Metrics{}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, "Error decode request JSON body", http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(res, "Error decode request JSON body", http.StatusBadRequest)
		return
	}

	switch metrics.MetricType(metric.MType) {
	case metrics.Gauge:
		savedMetric, err := s.Storage.GetGaugeMetric(metric.ID)
		if err != nil {
			metric.Value = &metrics.ZeroGaugeMetricValue
			break
		}

		value := savedMetric.GetValue()
		metric.Value = &value

	case metrics.Counter:
		savedMetric, err := s.Storage.GetCounterMetric(metric.ID)
		if err != nil {
			metric.Delta = &metrics.ZeroCounterMetricValue
			break
		}

		delta := savedMetric.GetValue()
		metric.Delta = &delta

	default:
		http.Error(res, fmt.Sprintf("Invalid metric type: %s", metric.MType), http.StatusBadRequest)
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
	gauges, counters := s.Storage.GetMetrics()
	body := getMetricsString(gauges, counters)

	res.Header().Set("Content-type", "text/html")
	_, err := io.WriteString(res, body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
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

func (s *Server) isSyncSaveStorageState() bool {
	isFileStorageEnabled := len(strings.TrimSpace(s.Config.StorageFilePath)) != 0
	isSyncSaveStorageState := s.Config.StoreInterval == 0
	return isFileStorageEnabled && isSyncSaveStorageState
}
