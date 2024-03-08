package main

import (
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/app"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/config"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
)

func main() {
	storage := storage.NewServerMemStorage()
	server := app.NewMetricsServer(&storage, config.MetricsServerConfig)

	err := server.Run()
	if err != nil {
		panic(err)
	}
}
