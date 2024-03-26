package server

import (
	"net/http"
	"strings"

	compress "github.com/Stern-Ritter/metrics-and-alerting-service/internal/compress/server"
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

	isFileStorageEnabled := len(strings.TrimSpace(config.StorageFilePath)) != 0
	if isFileStorageEnabled {
		err = s.AddFileStorage(config.StorageFilePath)
		if err != nil {
			logger.Log.Fatal(err.Error(), zap.String("event", "add file storage"))
			return err
		}
		logger.Log.Info("Success", zap.String("event", "add file storage"))

		if config.Restore {
			if err := s.FileStorage.Load(); err != nil {
				logger.Log.Fatal(err.Error(), zap.String("event", "restore data from file storage"))
				return err
			}
			logger.Log.Info("Success", zap.String("event", "restore from storage"))
		}

		s.FileStorage.SetSaveInterval(config.StoreInterval)
		defer s.FileStorage.Close()
	}

	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Use(compress.GzipMiddleware)
	r.Get("/", s.GetMetricsHandler)

	r.Route("/update", func(r chi.Router) {
		if isFileStorageEnabled {
			r.Use(s.FileStorage.FileStorageMiddleware)
		}
		r.Post("/", s.UpdateMetricHandlerWithBody)
		r.Post("/{type}/{name}/{value}", s.UpdateMetricHandlerWithPathVars)
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", s.GetMetricHandlerWithBody)
		r.Get("/{type}/{name}", s.GetMetricHandlerWithPathVars)
	})

	err = http.ListenAndServe(config.URL, r)
	if err != nil {
		logger.Log.Fatal(err.Error(), zap.String("event", "start server"))
	}
	return err
}
