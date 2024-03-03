package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/errors"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/utils"
)

var db = storage.NewMemStore()

func UpdateMetricHandler(res http.ResponseWriter, req *http.Request) {
	if !isRequestMethodAllowed(req.Method, []string{http.MethodPost}) {
		http.Error(res, "Only POST requests are allowed.", http.StatusMethodNotAllowed)
		return
	}

	// if !isRequqestContainsHeaderValue(req, "Content-Type", "text/plain") {
	// 	http.Error(res, "Only Content-Type: text/plain is allowed.", http.StatusUnsupportedMediaType)
	// 	return
	// }

	metricType, metricName, metricValue, err := parsePathVariables(req.URL.Path, "/update/")
	if err != nil {
		http.Error(res, "The resource you requested has not been found at the specified address", http.StatusNotFound)
		return
	}

	err = db.UpdateMetric(metricType, metricName, metricValue)
	switch err.(type) {
	case errors.InvalidMetricName, errors.InvalidMetricValue:
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func isRequestMethodAllowed(method string, allowedMethods []string) bool {
	return utils.Contains(allowedMethods, method)
}

func isRequqestContainsHeaderValue(req *http.Request, header string, value string) bool {
	contentTypeValues := req.Header.Values(header)
	return utils.Contains(contentTypeValues, value)
}

func parsePathVariables(path string, prefix string) (metricType string, metricName string, metricVakue string, err error) {
	pathVariables := strings.Split(strings.TrimPrefix(path, prefix), "/")
	if len(pathVariables) < 3 {
		return "", "", "", fmt.Errorf("invalid path variables")
	}

	return pathVariables[0], pathVariables[1], strings.Join(pathVariables[2:], "/"), nil
}
