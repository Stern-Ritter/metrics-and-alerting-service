package server

import (
	"net/http"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	service "github.com/Stern-Ritter/metrics-and-alerting-service/internal/service/server"
	"github.com/go-chi/chi"
)

func Run(s *service.Server) error {
	config, err := getConfig(config.ServerConfig{})
	if err != nil {
		return err
	}

	router := chi.NewRouter()
	router.Get("/", s.GetMetricsHandler)
	router.Post("/update/{type}/{name}/{value}", s.UpdateMetricHandler)
	router.Get("/value/{type}/{name}", s.GetMetricHandler)

	err = http.ListenAndServe(config.URL, router)
	return err
}
