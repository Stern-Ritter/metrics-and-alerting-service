package server

import (
	"net/http"

	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	service "github.com/Stern-Ritter/metrics-and-alerting-service/internal/service/server"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func Run(s *service.Server) error {
	config, err := getConfig(config.ServerConfig{
		LoggerLvl: "info",
	})
	if err != nil {
		return err
	}

	err = logger.Initialize(config.LoggerLvl)
	if err != nil {
		return err
	}

	router := chi.NewRouter()
	router.Use(logger.RequestLogger)
	router.Get("/", s.GetMetricsHandler)
	router.Post("/update", s.UpdateMetricHandlerWithBody)
	router.Post("/update/", s.UpdateMetricHandlerWithBody)
	router.Post("/update/{type}/{name}/{value}", s.UpdateMetricHandlerWithPathVars)
	router.Post("/value", s.GetMetricHandlerWithBody)
	router.Post("/value/", s.GetMetricHandlerWithBody)
	router.Get("/value/{type}/{name}", s.GetMetricHandlerWithPathVars)

	err = http.ListenAndServe(config.URL, router)
	if err != nil {
		logger.Log.Fatal(err.Error(), zap.String("event", "start server"))
	}
	return err
}
