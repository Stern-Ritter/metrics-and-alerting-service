package main

import (
	"flag"

	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/app"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/config"
	"github.com/Stern-Ritter/metrics-and-alerting-service/internal/storage"
)

var metricsServerConfig = config.ServerConfig{
	URL: config.URL{
		Host: "localhost",
		Port: 8080,
	},
}

func parseFlags() {
	flag.Var(&metricsServerConfig.URL, "a", "address and port to run server in format <host>:<port>")
	flag.Parse()
}

func main() {
	parseFlags()

	storage := storage.NewServerMemStorage()
	server := app.NewMetricsServer(&storage, metricsServerConfig)

	server.Run()
}
