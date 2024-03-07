package app

import (
	"net/http"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/config"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	handlers "github.com/Stern-Ritter/metrics-and-alerting-service/internal/transport"
)

type MetricsServer struct {
	storage *storage.ServerMemStorage
	config  config.ServerConfig
}

func NewMetricsServer(storage *storage.ServerMemStorage, config config.ServerConfig) MetricsServer {
	return MetricsServer{storage, config}
}

func (s *MetricsServer) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handlers.UpdateMetricHandler(s.storage))
	err := http.ListenAndServe(s.config.URL, mux)
	return err
}
