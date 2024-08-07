package server

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"

	compress "github.com/Stern-Ritter/metrics-and-alerting-service/internal/compress/server"
	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	service "github.com/Stern-Ritter/metrics-and-alerting-service/internal/service/server"
	storage "github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/migrations"
	pb "github.com/Stern-Ritter/metrics-and-alerting-service/proto/gen/metrics/metricsapi/v1"
)

// Run starts the server, setting up the storage and HTTP handlers.
// It returns an error if there are issues starting the server.
func Run(config *config.ServerConfig, logger *logger.ServerLogger) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	idleConnsClosed := make(chan struct{})

	metricStorage := createStorage(config, logger)
	defer metricStorage.Close()
	metricService := createMetricService(metricStorage, config, logger)
	rsaPrivateKey := getRsaPrivateKey(config.CryptoKeyPath, logger)
	trustedSubnet := getTrustedSubnet(config.TrustedSubnet, logger)

	server := service.NewServer(metricService, config, rsaPrivateKey, trustedSubnet, logger)

	if config.GRPC {
		err := runGrpcServer(server, signals, idleConnsClosed)
		return err
	} else {
		err := runHTTPServer(server, signals, idleConnsClosed)
		return err
	}
}

func createStorage(config *config.ServerConfig, logger *logger.ServerLogger) storage.Storage {
	var store storage.Storage

	isDatabaseEnabled := len(config.DatabaseDSN) > 0
	if isDatabaseEnabled {
		db, err := sql.Open("pgx", config.DatabaseDSN)
		if err != nil {
			logger.Fatal(err.Error(), zap.String("event", "connect database"))
		}

		dbStorage := storage.NewDBStorage(db, logger)
		logger.Info("Success", zap.String("event", "init database schema"))
		store = dbStorage
		logger.Info("Success", zap.String("event", "create database storage"))

		err = migrateDatabase(config.DatabaseDSN)
		if err != nil {
			logger.Fatal(err.Error(), zap.String("event", "migrate database"))
		}
	} else {
		store = storage.NewMemoryStorage(logger)
		logger.Info("Success", zap.String("event", "create in memory storage"))

		err := restoreMemoryStorageState(store, config)
		if err != nil {
			logger.Fatal(err.Error(), zap.String("event", "restore storage state from file"))
		}
		logger.Info("Success", zap.String("event", "restore storage state from file"))
	}

	return store
}

func migrateDatabase(databaseDsn string) error {
	goose.SetBaseFS(migrations.Migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose failed to set postgres dialect: %w", err)
	}

	db, err := goose.OpenDBWithDriver("pgx", databaseDsn)
	if err != nil {
		return fmt.Errorf("goose failed to open database connection: %w", err)
	}

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("goose failed to migrate database: %w", err)
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("goose failed to close database connection: %w", err)
	}

	return nil
}

func restoreMemoryStorageState(storage storage.Storage, config *config.ServerConfig) error {
	needRestoreState := config.Restore
	isFileStorageEnabled := len(config.FileStoragePath) > 0
	if needRestoreState && isFileStorageEnabled {
		return storage.Restore(config.FileStoragePath)
	}
	return nil
}

func createMetricService(storage storage.Storage, config *config.ServerConfig, logger *logger.ServerLogger) *service.MetricService {
	metricService := service.NewMetricService(storage, logger)

	isFileStorageEnabled := len(config.FileStoragePath) > 0
	isDatabaseNotEnabled := len(config.DatabaseDSN) == 0
	if isFileStorageEnabled && isDatabaseNotEnabled {
		metricService.SetSaveStateToFileInterval(config.FileStoragePath, config.StoreInterval)
	}

	return metricService
}

func getRsaPrivateKey(rsaPrivateKeyPath string, logger *logger.ServerLogger) *rsa.PrivateKey {
	rsaPrivateKey, err := service.GetRSAPrivateKey(rsaPrivateKeyPath)
	if err != nil {
		logger.Fatal(err.Error(), zap.String("event", "get rsa private key"))
	}

	return rsaPrivateKey
}

func getTrustedSubnet(cidrTrustedSubnet string, logger *logger.ServerLogger) *net.IPNet {
	trustedSubnet, err := service.GetTrustedSubnet(cidrTrustedSubnet)
	if err != nil {
		logger.Fatal(err.Error(), zap.String("event", "get trusted subnet for agents"))
	}

	return trustedSubnet
}

func runGrpcServer(server *service.Server, signals chan os.Signal, idleConnsClosed chan struct{}) error {
	listen, err := net.Listen("tcp", server.Config.URL)
	if err != nil {
		return err
	}

	opts := make([]grpc.ServerOption, 0)

	opts = append(opts, grpc.ChainUnaryInterceptor(logger.LoggerInterceptor(server.Logger)))
	opts = append(opts, grpc.ChainUnaryInterceptor(server.SubnetInterceptor))
	opts = append(opts, grpc.ChainUnaryInterceptor(server.SignInterceptor))
	isEncryptionEnabled := len(server.Config.TLSCertPath) > 0 && len(server.Config.TLSKeyPath) > 0
	if isEncryptionEnabled {
		creds, err := credentials.NewServerTLSFromFile(server.Config.TLSCertPath, server.Config.TLSKeyPath)
		if err != nil {
			server.Logger.Fatal(err.Error(), zap.String("event", "load credentials"))
		}
		opts = append(opts, grpc.Creds(creds))
	}

	srv := grpc.NewServer(opts...)
	pb.RegisterMetricsV1ServiceServer(srv, server)

	go func() {
		<-signals

		server.Logger.Info("Shutting down server", zap.String("event", "shutdown server"))
		srv.GracefulStop()

		err := server.MetricService.SaveStateToFile(server.Config.FileStoragePath)
		if err != nil {
			server.Logger.Error("Failed to save state to file", zap.String("event", "save state to file"),
				zap.Error(err))
		}

		close(idleConnsClosed)
	}()

	server.Logger.Info("Starting server", zap.String("event", "start server"))
	if err := srv.Serve(listen); err != nil {
		return err
	}

	<-idleConnsClosed
	server.Logger.Info("Server shutdown complete", zap.String("event", "shutdown server"))

	return nil
}

func runHTTPServer(server *service.Server, signals chan os.Signal, idleConnsClosed chan struct{}) error {
	r := addRoutes(server)
	srv := &http.Server{
		Addr:    server.Config.URL,
		Handler: r,
	}

	go func() {
		<-signals
		server.Logger.Info("Shutting down server", zap.String("event", "shutdown server"))

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(server.Config.ShutdownTimeout)*time.Second)
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
	r.Use(s.SubnetMiddleware)
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
