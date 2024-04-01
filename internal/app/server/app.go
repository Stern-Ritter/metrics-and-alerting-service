package server

import (
	"net/http"
	"strings"

	compress "github.com/Stern-Ritter/metrics-and-alerting-service/internal/compress/server"
	service "github.com/Stern-Ritter/metrics-and-alerting-service/internal/service/server"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func Run(s *service.Server) error {
	isFileStorageEnabled := len(strings.TrimSpace(s.Config.StorageFilePath)) != 0
	if isFileStorageEnabled && s.Config.Restore {
		if err := s.Storage.Restore(s.Config.StorageFilePath); err != nil {
			s.Logger.Fatal(err.Error(), zap.String("event", "restore storage state from file"))
			return err
		}
		s.Logger.Info("Success", zap.String("event", "restore storage state from file"))

		s.Storage.SetSaveInterval(s.Config.StorageFilePath, s.Config.StoreInterval)
	}

	r := addRoutes(s)
	err := http.ListenAndServe(s.Config.URL, r)
	if err != nil {
		s.Logger.Fatal(err.Error(), zap.String("event", "start server"))
	}
	return err
}

func addRoutes(s *service.Server) *chi.Mux {
	r := chi.NewRouter()
	r.Use(s.Logger.LoggerMiddleware)
	r.Use(compress.GzipMiddleware)
	r.Get("/", s.GetMetricsHandler)

	r.Route("/update", func(r chi.Router) {
		r.Post("/", s.UpdateMetricHandlerWithBody)
		r.Post("/{type}/{name}/{value}", s.UpdateMetricHandlerWithPathVars)
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", s.GetMetricHandlerWithBody)
		r.Get("/{type}/{name}", s.GetMetricHandlerWithPathVars)
	})

	return r
}
