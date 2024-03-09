package app

import (
	"log"
	"net/http"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/config"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	handlers "github.com/Stern-Ritter/metrics-and-alerting-service/internal/transport"
	"github.com/go-chi/chi"
)

type MetricsServer struct {
	storage *storage.ServerMemStorage
	config  config.ServerConfig
}

func NewMetricsServer(storage *storage.ServerMemStorage, config config.ServerConfig) MetricsServer {
	return MetricsServer{storage, config}
}

func (s *MetricsServer) Run() {
	router := chi.NewRouter()
	router.Get("/", handlers.GetMetricsHandler(s.storage))
	router.Post("/update/{type}/{name}/{value}", handlers.UpdateMetricHandler(s.storage))
	router.Get("/value/{type}/{name}", handlers.GetMetricHandler(s.storage))

	log.Fatal(http.ListenAndServe(s.config.URL.String(), router))
}
