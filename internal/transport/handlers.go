package transport

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

func UpdateMetricHandler(storage storage.ServerStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if !isRequestMethodAllowed(req.Method, []string{http.MethodPost}) {
			http.Error(res, "Only POST requests are allowed.", http.StatusMethodNotAllowed)
			return
		}

		metricType, metricName, metricValue, err := parsePathVariables(req.URL.Path, "/update/")
		if err != nil {
			http.Error(res, "The resource you requested has not been found at the specified address", http.StatusNotFound)
			return
		}

		err = storage.UpdateMetric(metricType, metricName, metricValue)
		switch err.(type) {
		case errors.InvalidMetricType, errors.InvalidMetricValue:
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		res.WriteHeader(http.StatusOK)
	}
}

func isRequestMethodAllowed(method string, allowedMethods []string) bool {
	return utils.Contains(allowedMethods, method)
}

func parsePathVariables(path string, prefix string) (string, string, string, error) {
	path = formatPath(path, prefix)
	pathVariables := strings.Split(path, "/")

	if len(pathVariables) != 3 {
		return "", "", "", fmt.Errorf("invalid path variables")
	}

	metricType := strings.TrimSpace(pathVariables[0])
	metricName := strings.TrimSpace(pathVariables[1])
	metricValue := strings.TrimSpace(pathVariables[2])
	if len(metricType) == 0 || len(metricName) == 0 || len(metricValue) == 0 {
		return "", "", "", fmt.Errorf("invalid path variables")
	}

	return metricType, metricName, metricValue, nil
}

func formatPath(path string, prefix string) string {
	path = strings.TrimPrefix(path, prefix)
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	return path
}
