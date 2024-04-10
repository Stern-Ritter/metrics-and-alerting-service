package server

import (
	"context"
	"database/sql"
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

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Run(config *config.ServerConfig, logger *logger.ServerLogger) error {
	isDatabaseEnabled := len(strings.TrimSpace(config.DatabaseDSN)) != 0

	var store storage.Storage

	if isDatabaseEnabled {
		conn, err := sql.Open("pgx", config.DatabaseDSN)
		if err != nil {
			logger.Fatal(err.Error(), zap.String("event", "connect database"))
			return err
		}
		defer conn.Close()

		dbStorage := storage.NewDBStorage(conn, logger)
		err = dbStorage.Bootstrap(context.TODO())
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

	isFileStorageEnabled := len(strings.TrimSpace(config.StorageFilePath)) != 0
	if !isDatabaseEnabled && isFileStorageEnabled && config.Restore {
		if err := server.MetricService.RestoreMetricsFromStorage(config.StorageFilePath); err != nil {
			server.Logger.Fatal(err.Error(), zap.String("event", "restore storage state from file"))
			return err
		}
		server.Logger.Info("Success", zap.String("event", "restore storage state from file"))

		server.MetricService.SetMetricsSaveInterval(server.Config.StorageFilePath, server.Config.StoreInterval)
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

	r.Route("/value", func(r chi.Router) {
		r.Post("/", s.GetMetricHandlerWithBody)
		r.Get("/{type}/{name}", s.GetMetricHandlerWithPathVars)
	})

	r.Get("/ping", s.PingDatabaseHandler)

	return r
}
