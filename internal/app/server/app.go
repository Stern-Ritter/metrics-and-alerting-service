package server

import (
	"context"
	"log"
	"net/http"
	"strings"

	compress "github.com/Stern-Ritter/metrics-and-alerting-service/internal/compress/server"
	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	service "github.com/Stern-Ritter/metrics-and-alerting-service/internal/service/server"
	storage "github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage/server"
	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5"
)

func Run(config *config.ServerConfig, logger *logger.ServerLogger) error {
	isDatabaseEnabled := len(strings.TrimSpace(config.DatabaseDSN)) != 0

	var store storage.Storage

	if isDatabaseEnabled {
		conn, err := pgx.Connect(context.Background(), config.DatabaseDSN)
		if err != nil {
			logger.Fatal(err.Error(), zap.String("event", "connect database"))
			return err
		}
		defer conn.Close(context.Background())

		dbStorage := storage.NewDBStorage(conn, logger)
		err = dbStorage.Bootstrap(context.Background())
		if err != nil {
			log.Fatal(err.Error(), zap.String("event", "init database schema"))
			return err
		}
		logger.Info("Success", zap.String("event", "init database schema"))
		store = dbStorage
		logger.Info("Success", zap.String("event", "create database storage"))
	} else {
		store = storage.NewMemoryStorage(logger)
		logger.Info("Success", zap.String("event", "create in memory storage"))
	}

	mService := service.NewMetricService(store, logger)
	server := service.NewServer(mService, config, logger)

	isFileStorageEnabled := len(strings.TrimSpace(config.FileStoragePath)) != 0
	if !isDatabaseEnabled && isFileStorageEnabled && config.Restore {
		if err := server.MetricService.RestoreStateFromFile(config.FileStoragePath); err != nil {
			server.Logger.Fatal(err.Error(), zap.String("event", "restore storage state from file"))
			return err
		}
		server.Logger.Info("Success", zap.String("event", "restore storage state from file"))

		server.MetricService.SetSaveStateToFileInterval(server.Config.FileStoragePath, server.Config.StoreInterval)
	}

	r := addRoutes(server)
	err := http.ListenAndServe(server.Config.URL, r)
	if err != nil {
		server.Logger.Fatal(err.Error(), zap.String("event", "start server"))
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

	r.Route("/updates", func(r chi.Router) {
		r.Post("/", s.UpdateMetricsBatchHandlerWithBody)
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", s.GetMetricHandlerWithBody)
		r.Get("/{type}/{name}", s.GetMetricHandlerWithPathVars)
	})

	r.Get("/ping", s.PingDatabaseHandler)

	return r
}
