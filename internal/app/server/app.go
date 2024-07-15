package server

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	compress "github.com/Stern-Ritter/metrics-and-alerting-service/internal/compress/server"
	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	crypto "github.com/Stern-Ritter/metrics-and-alerting-service/internal/crypto/server"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	service "github.com/Stern-Ritter/metrics-and-alerting-service/internal/service/server"
	storage "github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage/server"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Run starts the server, setting up the storage and HTTP handlers.
// It returns an error if there are issues starting the server.
func Run(config *config.ServerConfig, logger *logger.ServerLogger) error {
	idleConnsClosed := make(chan struct{})

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	isDatabaseEnabled := len(strings.TrimSpace(config.DatabaseDSN)) != 0

	var store storage.Storage

	if isDatabaseEnabled {
		db, err := sql.Open("pgx", config.DatabaseDSN)
		if err != nil {
			logger.Fatal(err.Error(), zap.String("event", "connect database"))
			return err
		}
		defer db.Close()

		dbStorage := storage.NewDBStorage(db, logger)
		logger.Info("Success", zap.String("event", "init database schema"))
		store = dbStorage
		logger.Info("Success", zap.String("event", "create database storage"))
	} else {
		store = storage.NewMemoryStorage(logger)
		logger.Info("Success", zap.String("event", "create in memory storage"))
	}

	mService := service.NewMetricService(store, logger)
	if isDatabaseEnabled {
		err := mService.MigrateDatabase(config.DatabaseDSN)
		if err != nil {
			logger.Fatal(err.Error(), zap.String("event", "migrate database"))
		}
	}

	var rsaPrivateKey *rsa.PrivateKey

	isEncryptionEnabled := len(strings.TrimSpace(config.CryptoKeyPath)) != 0
	if isEncryptionEnabled {
		key, err := crypto.GetRSAPrivateKey(config.CryptoKeyPath)
		if err != nil {
			logger.Fatal(err.Error(), zap.String("event", "get rsa private key"))
		}
		rsaPrivateKey = key
	}

	server := service.NewServer(mService, config, rsaPrivateKey, logger)

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

	srv := &http.Server{
		Addr:    server.Config.URL,
		Handler: r,
	}

	go func() {
		<-signals
		server.Logger.Info("Shutting down server", zap.String("event", "shutdown server"))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			server.Logger.Error("Failed to shutdown server", zap.String("event", "shutdown server"),
				zap.Error(err))
		}

		err := server.MetricService.SaveStateToFile(server.Config.FileStoragePath)
		if err != nil {
			server.Logger.Error("Failed to save state to file", zap.String("event", "save state to file"),
				zap.Error(err))
		}

		close(idleConnsClosed)
	}()

	server.Logger.Info("Starting server", zap.String("event", "start server"))
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-idleConnsClosed
	server.Logger.Info("Server shutdown complete", zap.String("event", "shutdown server"))

	return nil
}

func addRoutes(s *service.Server) *chi.Mux {
	r := chi.NewRouter()
	r.Use(s.Logger.LoggerMiddleware)
	r.Use(s.SignMiddleware)
	r.Use(s.EncryptMiddleware)
	r.Use(compress.GzipMiddleware)
	r.Get("/", s.GetMetricsHandler)

	r.Mount("/debug/pprof", http.DefaultServeMux)

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
