package main

import (
	"database/sql"
	"log"

	app "github.com/Stern-Ritter/metrics-and-alerting-service/internal/app/server"
	config "github.com/Stern-Ritter/metrics-and-alerting-service/internal/config/server"
	logger "github.com/Stern-Ritter/metrics-and-alerting-service/internal/logger/server"
	service "github.com/Stern-Ritter/metrics-and-alerting-service/internal/service/server"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	config, err := app.GetConfig(config.ServerConfig{
		LoggerLvl: "info",
	})
	if err != nil {
		log.Fatalf("%+v", err)
	}

	logger, err := logger.Initialize(config.LoggerLvl)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	db, err := sql.Open("pgx", config.DatabaseDSN)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer db.Close()

	dbStorage := storage.NewMetricStore(db)
	storage := storage.NewServerMemStorage(logger)
	metricService := service.NewMetricService(&dbStorage, &storage, logger)
	server := service.NewServer(metricService, &config, logger)

	err = app.Run(server)
	if err != nil {
		logger.Fatal(err.Error(), zap.String("event", "start server"))
	}
}
